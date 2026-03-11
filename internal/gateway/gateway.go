package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/TheKodeToad/fine/internal/config"
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/TheKodeToad/fine/internal/misc"
	"github.com/gorilla/websocket"
)

type sessionHandle struct {
	shutdown         chan<- struct{}
	shutdownFinished <-chan struct{}
}

// Gateway keeps track of gateway sessions and can be shut down with Shutdown.
type Gateway struct {
	sessionsLock sync.Mutex
	sessions     map[sessionHandle]bool
}

var upgrader websocket.Upgrader

func (g *Gateway) ServeHTTP(conf *config.Config, w http.ResponseWriter, r *http.Request) {
	shutdown := make(chan struct{})
	shutdownFinished := make(chan struct{}, 1)
	defer func() { shutdownFinished <- struct{}{} }()

	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(fmt.Errorf("failed to upgrade to websocket: %w", err))
	}
	defer clientConn.Close()

	slog.Debug("starting websocket session")

	fluxerURL := misc.New(*conf.FluxerGatewayURL)
	fluxerQuery := fluxerURL.Query()
	fluxerQuery.Add("v", conf.FluxerAPIVersion)
	fluxerURL.RawQuery = fluxerQuery.Encode()

	fluxerConn, _, err := websocket.DefaultDialer.Dial(fluxerURL.String(), nil)
	if err != nil {
		slog.Warn("connection to fluxer gateway failed", slog.Any("err", err))

		clientCloseMsgErr := clientConn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(
				discord.GatewayClosedUnknownError,
				"Connection to Fluxer gateway failed.",
			),
		)
		clientCloseErr := clientConn.Close()
		if clientCloseMsgErr != nil || clientCloseErr != nil {
			slog.Warn(
				"closing client connection failed",
				slog.Any("err", errors.Join(clientCloseMsgErr, clientCloseErr)),
			)
		}

		return
	}
	defer fluxerConn.Close()

	g.registerSession(sessionHandle{shutdown, shutdownFinished})
	defer g.unregisterSession(sessionHandle{shutdown, shutdownFinished})

	handleSession(clientConn, fluxerConn, shutdown)

	slog.Debug("websocket session finished")
}

func (g *Gateway) registerSession(sesh sessionHandle) {
	g.sessionsLock.Lock()
	defer g.sessionsLock.Unlock()

	if g.sessions == nil {
		g.sessions = map[sessionHandle]bool{}
	}
	g.sessions[sesh] = true
}

func (g *Gateway) unregisterSession(sesh sessionHandle) {
	g.sessionsLock.Lock()
	defer g.sessionsLock.Unlock()

	if g.sessions != nil {
		delete(g.sessions, sesh)
	}
}

func (g *Gateway) Shutdown() {
	// NOTE: unfortunately we can't defer the unlock as if we only unlock it after wg.Wait() it causes a deadlock -
	// we are waiting for the sessions to shut down which are waiting
	g.sessionsLock.Lock()

	if g.sessions == nil {
		g.sessionsLock.Unlock()
		return
	}

	var wg sync.WaitGroup

	wg.Add(len(g.sessions))
	for sesh := range g.sessions {
		go func() {
			sesh.shutdown <- struct{}{}
			<-sesh.shutdownFinished
			wg.Done()
		}()
	}

	g.sessions = nil
	g.sessionsLock.Unlock()

	wg.Wait()
}

type wsMessage struct {
	messageType int
	data        []byte
}

type wsChannels struct {
	read     <-chan wsMessage
	readErr  <-chan error
	write    chan<- wsMessage
	writeErr <-chan error
}

var errWriteChanClosed = errors.New("websocket write channel closed")

// setupWSChannels sets up some channels and spawns goroutines handling them.
func setupWSChannels(conn *websocket.Conn) wsChannels {
	read := make(chan wsMessage)
	readErr := make(chan error)

	// FIXME: goroutine leak
	go func() {
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				readErr <- err
				return
			}

			read <- wsMessage{msgType, msg}
		}
	}()

	write := make(chan wsMessage, 256)
	writeErr := make(chan error)

	go func() {
		for {
			msg, ok := <-write
			if !ok {
				// NOTE: we want to be able to close the channel and wait for all messages to be delivered
				// this signals the end of messages
				writeErr <- errWriteChanClosed
				return
			}

			err := conn.WriteMessage(msg.messageType, msg.data)
			if err != nil {
				writeErr <- err
				return
			}
		}
	}()

	return wsChannels{read, readErr, write, writeErr}
}

func handleSession(clientConn *websocket.Conn, fluxerConn *websocket.Conn, shutdown <-chan struct{}) {
	// allow both sides to naturally perform the close handshake
	clientConn.SetCloseHandler(func(int, string) error { return nil })
	fluxerConn.SetCloseHandler(func(int, string) error { return nil })

	client := setupWSChannels(clientConn)
	fluxer := setupWSChannels(fluxerConn)

loop:
	for {
		select {
		case msg := <-client.read:
			if msg.messageType != websocket.TextMessage {
				continue loop
			}

			var packet discord.Packet
			err := json.Unmarshal(msg.data, &packet)
			if err != nil {
				slog.Debug("failed to unmarshal client packet", slog.Any("err", err))
				client.write <- wsMessage{
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
				}
				break loop
			}

			packet, keep := packetToFluxer(packet)
			if !keep {
				continue loop
			}

			newData, err := json.Marshal(packet)
			if err != nil {
				slog.Debug("failed to marshal modified client packet", slog.Any("err", err))
				continue loop
			}

			client.write <- wsMessage{
				websocket.TextMessage,
				newData,
			}
		case err := <-client.readErr:
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				fluxer.write <- wsMessage{
					websocket.CloseMessage,
					websocket.FormatCloseMessage(closeErr.Code, closeErr.Text),
				}
			} else {
				slog.Warn("error reading from client; ending session", slog.Any("err", err))
				break loop
			}
		case err := <-client.writeErr:
			slog.Warn("error writing to client; ending session", slog.Any("err", err))
			client.writeErr = nil // don't receive again
			break loop
		case msg := <-fluxer.read:
			if msg.messageType != websocket.TextMessage {
				continue loop
			}

			var packet discord.Packet
			err := json.Unmarshal(msg.data, &packet)
			if err != nil {
				slog.Debug("failed to unmarshal fluxer gateway packet", slog.Any("err", err))
				client.write <- wsMessage{
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
				}
				break loop
			}

			packet, keep := packetToDiscord(packet)
			if !keep {
				continue loop
			}

			newData, err := json.Marshal(packet)
			if err != nil {
				slog.Debug("failed to marshal modified fluxer gateway packet", slog.Any("err", err))
				continue loop
			}

			client.write <- wsMessage{
				websocket.TextMessage,
				newData,
			}
		case err := <-fluxer.readErr:
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				fluxer.write <- wsMessage{
					websocket.CloseMessage,
					websocket.FormatCloseMessage(closeErr.Code, closeErr.Text),
				}
			} else {
				slog.Warn("error reading from fluxer gateway; ending session", slog.Any("err", err))
				break loop
			}
		case err := <-fluxer.writeErr:
			slog.Warn("error writing to fluxer gateway; ending session", slog.Any("err", err))
			fluxer.writeErr = nil // don't receive again
			break loop
		case <-shutdown:
			break loop
		}
	}

	// NOTE: close the channels so that reading does not block when empty
	close(client.write)
	close(fluxer.write)

	// make sure all buffered messages are sent
	if client.writeErr != nil {
		slog.Debug("waiting for messages to client to be sent")

		err := <-client.writeErr
		if err != errWriteChanClosed {
			slog.Warn("error writing to client after session ended", slog.Any("err", err))
		}
	}

	if fluxer.writeErr != nil {
		slog.Debug("waiting for messages to fluxer gateway to be sent")

		err := <-fluxer.writeErr
		if err != errWriteChanClosed {
			slog.Warn("error writing to fluxer gateway after session ended", slog.Any("err", err))
		}
	}
}
