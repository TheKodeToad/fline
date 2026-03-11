package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TheKodeToad/fine/internal/config"
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/TheKodeToad/fine/internal/misc"
	"github.com/gorilla/websocket"
)

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

type session struct {
	clientConn *websocket.Conn
	fluxerConn *websocket.Conn
	client     wsChannels
	fluxer     wsChannels
}

func startSession(conf *config.Config, w http.ResponseWriter, r *http.Request) (session, error) {
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return session{}, fmt.Errorf("failed to upgrade to websocket: %w", err)
	}

	slog.Debug("starting websocket session")

	fluxerURL := misc.New(*conf.FluxerGatewayURL)
	fluxerQuery := fluxerURL.Query()
	fluxerQuery.Add("v", conf.FluxerAPIVersion)
	fluxerURL.RawQuery = fluxerQuery.Encode()

	fluxerConn, _, err := websocket.DefaultDialer.Dial(fluxerURL.String(), nil)
	if err != nil {
		fluxerConnErr := fmt.Errorf("failed to connect to fluxer: %w", err)

		clientCloseMsgErr := clientConn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(
				discord.GatewayClosedUnknownError,
				"Connection to Fluxer gateway failed.",
			),
		)
		if clientCloseMsgErr != nil {
			clientCloseMsgErr = fmt.Errorf("failed to write close message to client: %w", clientCloseMsgErr)
		}

		clientCloseErr := clientConn.Close()
		if clientCloseErr != nil {
			clientCloseErr = fmt.Errorf("failed to close client connection: %w", clientCloseErr)
		}

		return session{}, errors.Join(fluxerConnErr, clientCloseMsgErr, clientCloseErr)
	}

	// allow both sides to naturally perform the close handshake
	clientConn.SetCloseHandler(func(int, string) error { return nil })
	fluxerConn.SetCloseHandler(func(int, string) error { return nil })

	var s session

	s.clientConn = clientConn
	s.fluxerConn = fluxerConn

	s.client = setupWSChannels(clientConn)
	s.fluxer = setupWSChannels(fluxerConn)

	return s, nil
}

// packetToFluxer converts from a Discord packet to a Fluxer packet.
func packetToFluxer(packet discord.Packet) (discord.Packet, bool) {
	fmt.Printf("converting packet to fluxer %+v %s\n", packet, string(packet.Data))

	return discord.Packet{}, false
}

func (s *session) handleClientMsg(msg wsMessage) error {
	if msg.messageType != websocket.TextMessage {
		return nil
	}

	var packet discord.Packet
	err := json.Unmarshal(msg.data, &packet)
	if err != nil {
		slog.Debug("failed to unmarshal client packet", slog.Any("err", err))
		s.client.write <- wsMessage{
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
		}
		return nil
	}

	packet, keep := packetToFluxer(packet)
	if !keep {
		return nil
	}

	newData, err := json.Marshal(packet)
	if err != nil {
		slog.Warn("failed to marshal modified client packet", slog.Any("err", err))
		return err
	}

	s.client.write <- wsMessage{
		websocket.TextMessage,
		newData,
	}
	return nil
}

// packetToDiscord converts from a Fluxer packet to a Discord packet.
func packetToDiscord(packet discord.Packet) (discord.Packet, bool) {
	fmt.Printf("converting packet to discord %+v %s\n", packet, string(packet.Data))

	switch packet.Opcode {
	case discord.GatewayOpHello:
		return packet, true
	}

	return discord.Packet{}, false
}

func (s *session) handleFluxerMsg(msg wsMessage) error {
	if msg.messageType != websocket.TextMessage {
		return nil
	}

	var packet discord.Packet
	err := json.Unmarshal(msg.data, &packet)
	if err != nil {
		s.client.write <- wsMessage{
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
		}
		return fmt.Errorf("failed to unmarshal fluxer packet: %w", err)
	}

	packet, keep := packetToDiscord(packet)
	if !keep {
		return nil
	}

	newData, err := json.Marshal(packet)
	if err != nil {
		slog.Debug("failed to marshal modified fluxer packet", slog.Any("err", err))
		return nil
	}

	s.client.write <- wsMessage{
		websocket.TextMessage,
		newData,
	}
	return nil
}

func (s *session) run(shutdown <-chan struct{}) {
	for {
		select {
		case msg := <-s.client.read:
			err := s.handleClientMsg(msg)
			if err != nil {
				slog.Warn("error handling client message", slog.Any("err", err))
			}
		case msg := <-s.fluxer.read:
			err := s.handleFluxerMsg(msg)
			if err != nil {
				slog.Warn("error handling fluxer message", slog.Any("err", err))
			}
		case err := <-s.client.readErr:
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				s.fluxer.write <- wsMessage{
					websocket.CloseMessage,
					websocket.FormatCloseMessage(closeErr.Code, closeErr.Text),
				}
			} else {
				slog.Warn("error reading from client; ending session", slog.Any("err", err))
				return
			}
		case err := <-s.fluxer.readErr:
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				s.fluxer.write <- wsMessage{
					websocket.CloseMessage,
					websocket.FormatCloseMessage(closeErr.Code, closeErr.Text),
				}
			} else {
				slog.Warn("error reading from fluxer; ending session", slog.Any("err", err))
				return
			}
		case err := <-s.client.writeErr:
			slog.Warn("error writing to client; ending session", slog.Any("err", err))
			s.client.writeErr = nil // don't receive again
			return
		case err := <-s.fluxer.writeErr:
			slog.Warn("error writing to fluxer; ending session", slog.Any("err", err))
			s.fluxer.writeErr = nil // ditto
			return
		case <-shutdown:
			return
		}
	}

}

func (s *session) close() error {
	// close the channels so that reading does not block when empty (but instead yields errWriteChanClosed)
	close(s.client.write)
	close(s.fluxer.write)

	// make sure all buffered messages are sent
	var clientWriteErr error
	if s.client.writeErr != nil {
		slog.Debug("waiting for messages to client to be sent")

		err := <-s.client.writeErr
		if err != errWriteChanClosed {
			clientWriteErr = fmt.Errorf("error writing to client after session ended: %w", err)
		}
	}

	var fluxerWriteErr error
	if s.fluxer.writeErr != nil {
		slog.Debug("waiting for messages to fluxer to be sent")

		err := <-s.fluxer.writeErr
		if err != errWriteChanClosed {
			fluxerWriteErr = fmt.Errorf("error writing to fluxer after session ended: %w", err)
		}
	}

	clientCloseErr := s.clientConn.Close()
	if clientCloseErr != nil {
		clientCloseErr = fmt.Errorf("error closing client connection: %w", clientCloseErr)
	}

	fluxerCloseErr := s.fluxerConn.Close()
	if fluxerCloseErr != nil {
		fluxerCloseErr = fmt.Errorf("error closing fluxer connection: %w", fluxerCloseErr)
	}

	return errors.Join(clientWriteErr, fluxerWriteErr, clientCloseErr, fluxerCloseErr)

}
