package discord

import (
	"encoding/json"

	"github.com/disgoorg/snowflake/v2"
)

const (
	GatewayClosedUnknownError         = 4000
	GatewayClosedUnknownOpcode        = 4001
	GatewayClosedDecodeError          = 4002
	GatewayClosedNotAuthenticated     = 4003
	GatewayClosedAuthenicationFailed  = 4004
	GatewayClosedAlreadyAuthenticated = 4005
	GatewayClosedInvalidSeq           = 4007
	GatewayClosedRateLimited          = 4008
	GatewayClosedSessionTimedOut      = 4009
	GatewayClosedInvalidShard         = 4010
	GatewayClosedShardingRequired     = 4011
	GatewayClosedInvalidAPIVersion    = 4012
	GatewayClosedInvalidIntents       = 4013
	GatewayClosedDisallowedIntents    = 4014
)

// NOTE: Fluxer's structure is identical.
type GatewayInfo struct {
	URL               string `json:"url"`
	Shards            int    `json:"shards"`
	SessionStartLimit struct {
		Total          int `json:"total"`
		Remaining      int `json:"remaining"`
		ResetAfter     int `json:"reset_after"`
		MaxConcurrency int `json:"max_concurrency"`
	} `json:"session_start_limit"`
}

type GatewayOpcode uint

var (
	GatewayOpDispatch            GatewayOpcode = 0
	GatewayOpHeartbeat           GatewayOpcode = 1
	GatewayOpIdentify            GatewayOpcode = 2
	GatewayOpPresenceUpdate      GatewayOpcode = 3
	GatewayOpVoiceStateUpdate    GatewayOpcode = 4
	GatewayOpResume              GatewayOpcode = 6
	GatewayOpReconnect           GatewayOpcode = 7
	GatewayOpRequestGuildMembers GatewayOpcode = 8
	GatewayOpInvalidSession      GatewayOpcode = 9
	GatewayOpHello               GatewayOpcode = 10
	GatewayOpHeartbeatAck        GatewayOpcode = 11
)

func (op GatewayOpcode) String() string {
	switch op {
	case GatewayOpDispatch:
		return "Dispatch"
	case GatewayOpHeartbeat:
		return "Heartbeat"
	case GatewayOpIdentify:
		return "Identify"
	case GatewayOpPresenceUpdate:
		return "PresenceUpdate"
	case GatewayOpVoiceStateUpdate:
		return "VoiceStateUpdate"
	case GatewayOpResume:
		return "Resume"
	case GatewayOpReconnect:
		return "Reconnect"
	case GatewayOpRequestGuildMembers:
		return "RequestGuildMembers"
	case GatewayOpInvalidSession:
		return "InvalidSession"
	case GatewayOpHello:
		return "Hello"
	case GatewayOpHeartbeatAck:
		return "HeartbeatAck"
	default:
		return ""
	}
}

// NOTE: Fluxer's structure is identical.
type Packet struct {
	Opcode      GatewayOpcode   `json:"op"`
	Data        json.RawMessage `json:"d"`
	SequenceNum *int            `json:"s,omitempty"`
	Event       string          `json:"t,omitempty"`
}

type ReadyEventApplication struct {
	ID    snowflake.ID     `json:"id"`
	Flags ApplicationFlags `json:"flags"`
}

type UnavailableGuild struct {
	ID          snowflake.ID `json:"id"`
	Unavailable *bool        `json:"unavailable,omitempty"`
}

type ReadyEvent struct {
	V                int                   `json:"v"`
	User             User                  `json:"user"`
	Guilds           []UnavailableGuild    `json:"guilds"`
	SessionID        string                `json:"session_id"`
	ResumeGatewayURL string                `json:"resume_gateway_url"`
	Shard            *[2]int               `json:"shard,omitempty"`
	Application      ReadyEventApplication `json:"application"`
	// NOTE: this is just an (undocumented) empty array Discord keeps for backwards compatibility.
	// Without this Eris will not like the packet.
	PrivateChannels [0]struct{} `json:"private_channels"`
}

type MessageCreateEvent struct {
	Message
	GuildID *snowflake.ID `json:"guild_id,omitempty"`
}
