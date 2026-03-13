package fluxer

import (
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/disgoorg/snowflake/v2"
)

type ReadyEvent struct {
	Version          int                        `json:"version"`
	User             UserPrivate                `json:"user"`
	Guilds           []discord.UnavailableGuild `json:"guilds"`
	SessionID        string                     `json:"session_id"`
	ResumeGatewayURL *string                    `json:"resume_gateway_url"`
}

type MessageCreateEvent struct {
	Message
	GuildID *snowflake.ID `json:"guild_id,omitempty"`
}
