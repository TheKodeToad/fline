package fluxer

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/disgoorg/snowflake/v2"
)

type UpdatePresencePayload struct {
	Status       string        `json:"status"`
	Mobile       bool          `json:"mobile"`
	AFK          bool          `json:"afk"`
	CustomStatus *CustomStatus `json:"custom_status,omitempty"`
}

type IdentifyPayload struct {
	Token          string                     `json:"token"`
	Properties     discord.IdentifyProperties `json:"properties"`
	Compress       *bool                      `json:"compress,omitempty"`
	LargeThreshold *int                       `json:"large_thresold,omitempty"`
	Presence       *UpdatePresencePayload     `json:"presence,omitempty"`
}

type ReadyEvent struct {
	Version          int                        `json:"version"`
	User             UserPrivate                `json:"user"`
	Guilds           []discord.UnavailableGuild `json:"guilds"`
	SessionID        string                     `json:"session_id"`
	ResumeGatewayURL *string                    `json:"resume_gateway_url"`
}

type GuildCreateEvent struct {
	Properties  Guild         `json:"properties"`
	JoinedAt    string        `json:"joined_at"`
	Large       bool          `json:"large"`
	Unavailable *bool         `json:"unavailable,omitempty"`
	MemberCount int           `json:"member_count"`
	Members     []GuildMember `json:"members"`
	Channels    []Channel     `json:"channels"`
	// NOTE: these fields are included inside a Discord guild, but Fluxer only sends them in the gateway event
	Roles    []Role            `json:"roles"`
	Emojis   []discord.Emoji   `json:"emojis"`
	Stickers []discord.Sticker `json:"stickers,omitzero"`
}

type GuildMembersChunkEvent struct {
	GuildID    snowflake.ID   `json:"guild_id"`
	Members    []GuildMember  `json:"members"`
	ChunkIndex int            `json:"chunk_index"`
	ChunkCount int            `json:"chunk_count"`
	NotFound   []snowflake.ID `json:"not_found,omitzero"`
	Presences  []todo         `json:"presences,omitzero"`
	Nonce      *string        `json:"nonce,omitempty"`
}

type MessageCreateEvent struct {
	Message
	GuildID *snowflake.ID `json:"guild_id,omitempty"`
	Member  *GuildMember  `json:"member,omitempty"`
}

type CustomStatus struct {
	Text          *string       `json:"text"`
	EmojiID       *snowflake.ID `json:"emoji_id"`
	EmojiName     *string       `json:"emoji_name"`
	EmojiAnimated bool          `json:"emoji_animated"`
	ExpiresAt     *string       `json:"expires_at"`
}

type PresenceUpdateEvent struct {
	GuildID *snowflake.ID `json:"guild_id"`
	User    UserPartial   `json:"user"`
	// TODO: draw the rest of the owl...
}
