package discord

import "encoding/json"

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

// NOTE: Fluxer's structure is identical.
type Packet struct {
	Opcode      int             `json:"op"`
	Data        json.RawMessage `json:"d"`
	SequenceNum int             `json:"s"`
	Event       string          `json:"t"`
}
