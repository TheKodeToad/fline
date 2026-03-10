package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"sync"

	"github.com/TheKodeToad/fine/internal/config"
	"github.com/TheKodeToad/fine/internal/convert"
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/TheKodeToad/fine/internal/misc"
	"github.com/gorilla/websocket"
)

// Gateway keeps track of gateway sessions and can be shut down with Shutdown.
type Gateway struct {
	sessionsLock sync.Mutex
	sessions     map[sessionHandle]bool
}

var upgrader websocket.Upgrader

func (g *Gateway) ServeHTTP(conf *config.Config, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(fmt.Errorf("failed to upgrade to websocket: %w", err))
	}

	slog.Debug("received websocket connection")
	handleSession(conf, g, conn)
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
	g.sessionsLock.Lock()
	defer g.sessionsLock.Unlock()

	if g.sessions == nil {
		return
	}

	var wg sync.WaitGroup

	for sesh := range g.sessions {
		wg.Add(1)
		go func() {
			sesh.shutdown <- struct{}{}
			<-sesh.shutdownFinished
			wg.Done()
		}()
	}

	g.sessions = nil

	wg.Wait()
}

type sessionHandle struct {
	shutdown         chan<- struct{}
	shutdownFinished <-chan struct{}
}

type wsMessage struct {
	messageType int
	data        []byte
}

type wsIO struct {
	read     <-chan wsMessage
	readErr  <-chan error
	write    chan<- wsMessage
	writeErr <-chan error
}

var errWriteChanClosed = errors.New("websocket write channel closed")

// wsHandleIO sets up some channels and spawns goroutines handling them.
func wsHandleIO(conn *websocket.Conn) wsIO {
	read := make(chan wsMessage)
	readErr := make(chan error)

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

	return wsIO{read, readErr, write, writeErr}
}

// handleSession starts a session and blocks until it is done.
// The connection is closed at the end.
func handleSession(conf *config.Config, gateway *Gateway, clientConn *websocket.Conn) {
	fluxerURL := misc.New(*conf.FluxerGatewayURL)
	fluxerQuery := fluxerURL.Query()
	fluxerQuery.Add("v", conf.FluxerAPIVersion)
	fluxerURL.RawQuery = fluxerQuery.Encode()

	fluxerConn, _, err := websocket.DefaultDialer.Dial(fluxerURL.String(), nil)
	if err != nil {
		clientConn.Close()

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

	shutdown := make(chan struct{})
	shutdownFinished := make(chan struct{}, 1)

	gateway.registerSession(sessionHandle{shutdown, shutdownFinished})
	defer gateway.unregisterSession(sessionHandle{shutdown, shutdownFinished})

	// allow both sides to naturally perform the close handshake
	clientConn.SetCloseHandler(func(int, string) error { return nil })
	fluxerConn.SetCloseHandler(func(int, string) error { return nil })

	client := wsHandleIO(clientConn)
	fluxer := wsHandleIO(fluxerConn)

loop:
	for {
		select {
		case msg := <-client.read:
			if msg.messageType != websocket.TextMessage {
				break // out of select
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

			fmt.Printf("received from client %+v\n", packet)
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
				break // out of select
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

			fmt.Printf("received from fluxer gateway %+v\n", packet)
		case err := <-fluxer.readErr:
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				fluxer.write <- wsMessage{
					websocket.CloseMessage,
					websocket.FormatCloseMessage(
						convert.GatewayCloseToDiscord(closeErr.Code, closeErr.Text),
					),
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

	clientConn.Close()
	fluxerConn.Close()

	shutdownFinished <- struct{}{}
}
