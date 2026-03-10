package gateway

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
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
	sessions     map[*session]bool
}

var upgrader websocket.Upgrader

func (g *Gateway) ServeHTTP(conf *config.Config, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(fmt.Errorf("failed to upgrade to websocket: %w", err))
	}

	slog.Debug("received websocket connection")
	sesh, err := startSession(conf, g, conn)
	if err != nil {
		slog.Warn("failed to start websocket session", slog.Any("err", err))
	}

	sesh.run()
}

func (g *Gateway) registerSession(sesh *session) {
	g.sessionsLock.Lock()
	defer g.sessionsLock.Unlock()

	if g.sessions == nil {
		g.sessions = map[*session]bool{}
	}
	g.sessions[sesh] = true
}

func (g *Gateway) unregisterSession(session *session) {
	g.sessionsLock.Lock()
	defer g.sessionsLock.Unlock()

	if g.sessions != nil {
		delete(g.sessions, session)
	}
}

func (g *Gateway) Shutdown() {
	g.sessionsLock.Lock()
	defer g.sessionsLock.Unlock()

	if g.sessions == nil {
		return
	}

	for sesh := range g.sessions {
		err := sesh.close()
		if err != nil {
			slog.Warn(fmt.Sprintf("error closing session %p on shutdown", sesh), slog.Any("err", err))
		}
	}

	g.sessions = nil
}

type session struct {
	logger          slog.Logger
	clientConn      *websocket.Conn
	fluxerConn      *websocket.Conn
	fluxerWriteLock sync.Mutex
	clientWriteLock sync.Mutex
}

func startSession(conf *config.Config, gateway *Gateway, clientConn *websocket.Conn) (*session, error) {
	fluxerURL := misc.New(*conf.FluxerGatewayURL)
	fluxerQuery := fluxerURL.Query()
	fluxerQuery.Add("v", conf.FluxerAPIVersion)
	fluxerURL.RawQuery = fluxerQuery.Encode()

	fluxerConn, _, err := websocket.DefaultDialer.Dial(fluxerURL.String(), nil)
	if err != nil {
		clientCloseMsgErr := clientConn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(
				discord.GatewayClosedUnknownError,
				"Connection to Fluxer gateway failed.",
			),
		)
		clientCloseErr := clientConn.Close()
		if clientCloseMsgErr != nil || clientCloseErr != nil {
			return nil, fmt.Errorf(
				"connection to fluxer gateway failed: %w; closing client connection failed: %w",
				err,
				errors.Join(clientCloseMsgErr, clientCloseErr),
			)
		}

		return nil, fmt.Errorf("connection to fluxer gateway failed: %w", err)
	}

	// allow both sides to naturally perform the close handshake
	clientConn.SetCloseHandler(func(int, string) error { return nil })
	fluxerConn.SetCloseHandler(func(int, string) error { return nil })

	var s session

	s.logger = *slog.Default().With("session", fmt.Sprintf("%p", &s))

	s.clientConn = clientConn
	s.fluxerConn = fluxerConn

	gateway.registerSession(&s)

	go s.readFromClient(gateway)
	go s.readFromFluxer(gateway)

	return &s, nil
}

// close closes the session ignoring any errors if any of the connections are already closed.
func (s *session) close() error {
	s.logger.Debug("closing session")

	clientCloseErr := s.clientConn.Close()
	if errors.Is(clientCloseErr, net.ErrClosed) {
		clientCloseErr = nil
	}

	fluxerCloseErr := s.fluxerConn.Close()
	if errors.Is(fluxerCloseErr, net.ErrClosed) {
		fluxerCloseErr = nil
	}

	if clientCloseErr != nil || fluxerCloseErr != nil {
		return errors.Join(clientCloseErr, fluxerCloseErr)
	}

	return nil
}

func (s *session) writeClientMsg(messageType int, data []byte) error {
	s.clientWriteLock.Lock()
	defer s.clientWriteLock.Unlock()

	return s.clientConn.WriteMessage(messageType, data)
}

func (s *session) writeFluxerMsg(messageType int, data []byte) error {
	s.fluxerWriteLock.Lock()
	defer s.fluxerWriteLock.Unlock()

	return s.fluxerConn.WriteMessage(messageType, data)
}

func (s *session) readFromClient(gateway *Gateway) {
	for {
		var packet discord.Packet

		err := s.clientConn.ReadJSON(&packet)
		if err != nil {
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				fmt.Printf("Forwarding close to fluxer: %v\n", err)
				err := s.writeFluxerMsg(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(closeErr.Code, closeErr.Text),
				)
				if err != nil {
					s.logger.Warn("failed to forward client close message to fluxer gateway", slog.Any("err", err))
				}
			} else {
				if !errors.Is(err, net.ErrClosed) {
					s.logger.Warn("error reading message from client", slog.Any("err", err))
				}

				err := s.close()
				if err != nil {
					s.logger.Warn("failed to close session", slog.Any("err", err))
				}
				gateway.unregisterSession(s)
			}

			return
		}

		fmt.Printf("from client: %+v\n", packet)
	}
}

func (s *session) readFromFluxer(gateway *Gateway) {
	for {
		var packet discord.Packet

		err := s.fluxerConn.ReadJSON(&packet)
		if err != nil {
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				var code int
				var msg string
				if closeErr != nil {
					code, msg = convert.GatewayCloseToDiscord(closeErr.Code, closeErr.Text)
				} else {
					code, msg = discord.GatewayClosedUnknownError, "Fluxer gateway connection disappeared."
				}

				err := s.writeClientMsg(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(code, msg),
				)
				if err != nil {
					s.logger.Warn("failed to forward fluxer gateway close message to client", slog.Any("err", err))
				}
			} else {
				if !errors.Is(err, net.ErrClosed) {
					s.logger.Warn("error reading message from fluxer gateway", slog.Any("err", err))
				}

				err := s.close()
				if err != nil {
					s.logger.Warn("failed to close session", slog.Any("err", err))
				}
				gateway.unregisterSession(s)
			}

			return
		}

		fmt.Printf("from fluxer: %+v\n", packet)
	}
}
