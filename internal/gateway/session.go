package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	fline "github.com/TheKodeToad/fline/internal"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
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
	readErr := make(chan error, 1)

	go func() {
		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				readErr <- err
				return
			}

			// FIXME: maybe in theory this could block forever if just as a message is read the run loop is ending?
			// that's pretty bad
			read <- wsMessage{msgType, msg}
		}
	}()

	write := make(chan wsMessage, 256)
	writeErr := make(chan error, 1)

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

type sessionInfo struct {
	// host contains the host header in the initial request.
	host string
	// apiVersion contains the version query parameter in the original request.
	apiVersion int
}

type session struct {
	info       sessionInfo
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

	versionStr := r.URL.Query().Get("v")
	apiVersion, err := strconv.Atoi(versionStr)
	if err != nil {
		// FIXME: this should wait for a response!
		clientCloseMsgErr := clientConn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, ""),
		)
		if clientCloseMsgErr != nil {
			clientCloseMsgErr = fmt.Errorf("failed to write close message to client: %w", clientCloseMsgErr)
		}

		clientCloseErr := clientConn.Close()
		if clientCloseErr != nil {
			clientCloseErr = fmt.Errorf("failed to close client: %w", clientCloseErr)
		}

		return session{}, errors.Join(clientCloseMsgErr, clientCloseErr)
	}

	fluxerURL := misc.New(*conf.FluxerGatewayURL)
	fluxerQuery := fluxerURL.Query()
	fluxerQuery.Add("v", fline.FluxerAPIVersion)
	fluxerURL.RawQuery = fluxerQuery.Encode()

	fluxerConn, _, err := websocket.DefaultDialer.Dial(fluxerURL.String(), nil)
	if err != nil {
		fluxerConnErr := fmt.Errorf("failed to connect to fluxer: %w", err)

		// FIXME: this should wait for a response!
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

	return session{
		info: sessionInfo{
			host:       r.Host,
			apiVersion: apiVersion,
		},
		clientConn: clientConn,
		fluxerConn: fluxerConn,
		client:     setupWSChannels(clientConn),
		fluxer:     setupWSChannels(fluxerConn),
	}, nil
}

func logPacket(msg string, packet discord.Packet) error {
	if !slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		return nil
	}

	var params []any

	if packet.SequenceNum != nil {
		params = append(
			params,
			slog.Int("seq", *packet.SequenceNum),
		)
	}

	data, err := json.MarshalIndent(packet.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal packet data: %w", err)
	}

	params = append(
		params,
		slog.String("opcode", packet.Opcode.String()),
		slog.String("event", packet.Event),
		slog.String("data", string(data)),
	)

	slog.Debug(msg, params...)
	return nil
}

// errNonConvertiblePacket signals that the packet is not convertible and nothing should be sent to the destination.
var errNonConvertiblePacket = errors.New("non-convertible packet")

// packetToFluxer converts from a Discord packet to a Fluxer packet.
func packetToFluxer(packet discord.Packet) (discord.Packet, error) {
	switch packet.Opcode {
	case discord.GatewayOpHeartbeat,
		discord.GatewayOpRequestGuildMembers,
		discord.GatewayOpResume:
		return packet, nil
	case discord.GatewayOpIdentify:
		var inPayload discord.IdentifyPayload
		err := json.Unmarshal(packet.Data, &inPayload)
		if err != nil {
			return discord.Packet{}, err
		}

		outPayload := convert.IdentifyPayloadToFluxer(inPayload)
		newData, err := json.Marshal(outPayload)
		if err != nil {
			return discord.Packet{}, err
		}

		packet.Data = newData
		return packet, nil
	case discord.GatewayOpPresenceUpdate:
		var inPayload discord.UpdatePresencePayload
		err := json.Unmarshal(packet.Data, &inPayload)
		if err != nil {
			return discord.Packet{}, err
		}

		outPayload := convert.PresenceUpdatePayloadToFluxer(inPayload)
		newData, err := json.Marshal(outPayload)
		if err != nil {
			return discord.Packet{}, err
		}

		packet.Data = newData
		return packet, nil
	default:
		return discord.Packet{}, errNonConvertiblePacket
	}
}

func (s *session) handleClientMsg(msg wsMessage) error {
	if msg.messageType != websocket.TextMessage {
		return nil
	}

	var inPacket discord.Packet
	err := json.Unmarshal(msg.data, &inPacket)
	if err != nil {
		s.client.write <- wsMessage{
			websocket.CloseMessage,
			websocket.FormatCloseMessage(discord.GatewayClosedDecodeError, ""),
		}
		return fmt.Errorf("failed to unmarshal client packet: %w", err)
	}

	err = logPacket("received discord packet", inPacket)
	if err != nil {
		slog.Warn("failed to log discord packet", slog.Any("err", err))
	}

	outPacket, err := packetToFluxer(inPacket)
	if errors.Is(err, errNonConvertiblePacket) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to convert packet to fluxer: %w", err)
	}

	newData, err := json.Marshal(outPacket)
	if err != nil {
		return fmt.Errorf("failed to marshal modified client packet: %w", err)
	}

	s.fluxer.write <- wsMessage{
		websocket.TextMessage,
		newData,
	}
	return nil
}

func eventToDiscord(name string, payload json.RawMessage, info sessionInfo) (json.RawMessage, error) {
	switch name {
	case "READY":
		var inEvent fluxer.ReadyEvent

		err := json.Unmarshal(payload, &inEvent)
		if err != nil {
			return json.RawMessage{}, err
		}

		// do libraries actually care about this?
		// for now just reassuringly echo back the version that was connected with...
		inEvent.Version = info.apiVersion

		if inEvent.ResumeGatewayURL == nil {
			inEvent.ResumeGatewayURL = misc.New(fline.GatewayURL(info.host))
		}

		outEvent := convert.ReadyEventToDiscord(inEvent)
		// NOTE: Fluxer doesn't currently support sharding, but some libs may break without this :)
		outEvent.Shard = misc.New([2]int{0, 1})
		return json.Marshal(outEvent)
	case "RESUMED":
		return payload, nil
	case "CHANNEL_CREATE", "CHANNEL_UPDATE":
		var inChannel fluxer.Channel

		err := json.Unmarshal(payload, &inChannel)
		if err != nil {
			return json.RawMessage{}, err
		}

		outChannel, ok := convert.ChannelToDiscord(inChannel)
		if !ok {
			return nil, errNonConvertiblePacket
		}

		return json.Marshal(outChannel)
	case "GUILD_STICKERS_UPDATE":
		var inEvent fluxer.GuildStickersUpdateEvent

		err := json.Unmarshal(payload, &inEvent)
		if err != nil {
			return json.RawMessage{}, nil
		}
		
		outEvent := convert.GuildStickersUpdateEventToDiscord(inEvent)
		return json.Marshal(outEvent)
	case "GUILD_EMOJI_UPDATE", "GUILD_DELETE", "GUILD_ROLE_DELETE":
		return payload, nil
	case "GUILD_CREATE":
		var inEvent fluxer.GuildCreateEvent

		err := json.Unmarshal(payload, &inEvent)
		if err != nil {
			return json.RawMessage{}, err
		}

		outEvent := convert.GuildCreateEventToDiscord(inEvent)
		return json.Marshal(outEvent)
	case "GUILD_MEMBER_ADD":
		var inEvent fluxer.GuildMemberAddEvent

		err := json.Unmarshal(payload, &inEvent)
		if err != nil {
			return json.RawMessage{}, err
		}

		outEvent := convert.GuildMemberAddEventToDiscord(inEvent)
		return json.Marshal(outEvent)
	case "GUILD_MEMBER_UPDATE":
		var inEvent fluxer.GuildMemberUpdateEvent

		err := json.Unmarshal(payload, &inEvent)
		if err != nil {
			return json.RawMessage{}, err
		}

		outEvent := convert.GuildMemberUpdateEventToDiscord(inEvent)
		return json.Marshal(outEvent)
	case "GUILD_MEMBERS_CHUNK":
		var inEvent fluxer.GuildMembersChunkEvent

		err := json.Unmarshal(payload, &inEvent)
		if err != nil {
			return json.RawMessage{}, err
		}

		outEvent := convert.GuildMembersChunkEventToDiscord(inEvent)
		return json.Marshal(outEvent)
	case "GUILD_ROLE_CREATE", "GUILD_ROLE_UPDATE":
		var inEvent fluxer.GuildRoleEvent

		err := json.Unmarshal(payload, &inEvent)
		if err != nil {
			return json.RawMessage{}, err
		}

		outEvent := convert.GuildRoleEventToDiscord(inEvent)
		return json.Marshal(outEvent)
	case "MESSAGE_CREATE":
		var inEvent fluxer.MessageCreateEvent

		err := json.Unmarshal(payload, &inEvent)
		if err != nil {
			return json.RawMessage{}, err
		}

		outEvent := convert.MessageCreateEventToDiscord(inEvent)
		return json.Marshal(outEvent)
	case "MESSAGE_REACTION_ADD":
		var inEvent fluxer.MessageReactionAddEvent

		err := json.Unmarshal(payload, &inEvent)
		if err != nil {
			return json.RawMessage{}, err
		}

		outEvent := convert.MessageReactionAddEventToDiscord(inEvent)
		return json.Marshal(outEvent)
	case "MESSAGE_REACTION_REMOVE":
		var inEvent fluxer.MessageReactionRemoveEvent

		err := json.Unmarshal(payload, &inEvent)
		if err != nil {
			return json.RawMessage{}, err
		}

		outEvent := convert.MessageReactionRemoveEventToDiscord(inEvent)
		return json.Marshal(outEvent)
	case "MESSAGE_REACTION_REMOVE_ALL":
		return payload, nil
	default:
		slog.Warn("received unknown event from fluxer: " + name)
		return json.RawMessage{}, errNonConvertiblePacket
	}
}

func packetToDiscord(packet discord.Packet, info sessionInfo) (discord.Packet, error) {
	switch packet.Opcode {
	case discord.GatewayOpHello,
		discord.GatewayOpHeartbeat,
		discord.GatewayOpHeartbeatAck,
		discord.GatewayOpReconnect,
		discord.GatewayOpInvalidSession:
		// passthrough
		return packet, nil
	case discord.GatewayOpDispatch:
		newData, err := eventToDiscord(packet.Event, packet.Data, info)
		if err != nil {
			return discord.Packet{}, fmt.Errorf("failed to convert event to discord: %w", err)
		}

		packet.Data = newData
		return packet, nil
	default:
		return discord.Packet{}, errNonConvertiblePacket
	}
}

func (s *session) handleFluxerMsg(msg wsMessage) error {
	if msg.messageType != websocket.TextMessage {
		return nil
	}

	var inPacket discord.Packet
	err := json.Unmarshal(msg.data, &inPacket)
	if err != nil {
		// FIXME: this should wait for a response!
		s.client.write <- wsMessage{
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
		}
		return fmt.Errorf("failed to unmarshal fluxer packet: %w", err)
	}

	err = logPacket("received fluxer packet", inPacket)
	if err != nil {
		slog.Warn("failed to log fluxer packet", slog.Any("err", err))
	}

	outPacket, err := packetToDiscord(inPacket, s.info)
	if errors.Is(err, errNonConvertiblePacket) {
		if inPacket.SequenceNum != nil {
			// NOTE: make sure there is no skipping of sequence numbers
			// this may generate other warnings, but at least the sequence number the client sends won't be outdated
			outPacket = discord.Packet{
				Opcode: discord.GatewayOpDispatch,
				SequenceNum: inPacket.SequenceNum,
				Event: "FLINE_NON_CONVERTIBLE",
			}
		} else {
			return nil
		}
	} else if err != nil {
		return fmt.Errorf("failed to convert packet to discord: %w", err)
	}

	newData, err := json.Marshal(outPacket)
	if err != nil {
		return fmt.Errorf("failed to marshal modified fluxer packet: %w", err)
	}

	s.client.write <- wsMessage{
		websocket.TextMessage,
		newData,
	}
	return nil
}

func (s *session) run(shutdown <-chan struct{}) {
	var sentClientClose bool
	var sentFluxerClose bool

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
				if closeErr.Code == websocket.CloseAbnormalClosure {
					slog.Debug("client connection closed without sending close message")
					return
				}

				slog.Debug(
					"forwarding close message to fluxer",
					slog.Int("code", closeErr.Code),
					slog.String("text", closeErr.Text),
				)

				sentFluxerClose = true
				s.fluxer.write <- wsMessage{
					websocket.CloseMessage,
					websocket.FormatCloseMessage(closeErr.Code, closeErr.Text),
				}

				if sentClientClose {
					return
				}
			} else {
				slog.Warn("error reading from client; ending session", slog.Any("err", err))
				return
			}
		case err := <-s.fluxer.readErr:
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				if closeErr.Code == websocket.CloseAbnormalClosure {
					slog.Debug("fluxer connection closed without sending close message")
					return
				}

				slog.Debug(
					"forwarding close message to client",
					slog.Int("code", closeErr.Code),
					slog.String("text", closeErr.Text),
				)

				sentClientClose = true
				s.client.write <- wsMessage{
					websocket.CloseMessage,
					websocket.FormatCloseMessage(closeErr.Code, closeErr.Text),
				}

				if sentFluxerClose {
					return
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
