package fluxer

import "github.com/disgoorg/snowflake/v2"

type ReadyEvent struct {
	Version int         `json:"version"`
	User    PrivateUser `json:"user"`
	Guilds  []struct {
		ID          snowflake.ID `json:"id"`
		Unavailable bool         `json:"unavailable"`
	} `json:"guilds"`
	SessionID        string  `json:"session_id"`
	ResumeGatewayURL *string `json:"resume_gateway_url"`
}
