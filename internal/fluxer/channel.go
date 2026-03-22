package fluxer

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/disgoorg/snowflake/v2"
)

type Channel struct {
	ID                   snowflake.ID                  `json:"id"`
	Type                 discord.ChannelType           `json:"type"`
	GuildID              *snowflake.ID                 `json:"guild_id,omitempty"`
	Position             *int                          `json:"position,omitempty"`
	PermissionOverwrites []discord.PermissionOverwrite `json:"permission_overwrites,omitzero"`
	Name                 *string                       `json:"name,omitempty"`
	Topic                *string                       `json:"topic"`
	NSFW                 *bool                         `json:"nsfw"`
	LastMessageID        *snowflake.ID                 `json:"last_message_id"`
	Bitrate              *int                          `json:"bitrate,omitempty"`
	UserLimit            *int                          `json:"user_limit,omitempty"`
	RateLimitPerUser     *int                          `json:"rate_limit_per_user,omitempty"`
	Recipients           []UserPartial                 `json:"recipients,omitempty"`
	Icon                 *string                       `json:"icon"`
	OwnerID              *snowflake.ID                 `json:"owner_id,omitempty"`
	ApplicationID        *snowflake.ID                 `json:"application_id,omitempty"`
	ParentID             *snowflake.ID                 `json:"parent_id"`
	LastPinTimestamp     *string                       `json:"last_pin_timestamp,omitempty"`
	RTCRegion            *string                       `json:"rtc_region"`
}

type Webhook struct {
	ID        snowflake.ID `json:"id"`
	GuildID   snowflake.ID `json:"guild_id"`
	ChannelID snowflake.ID `json:"channel_id"`
	User      UserPartial  `json:"user"`
	Name      string       `json:"name"`
	Avatar    *string      `json:"avatar"`
	Token     string       `json:"token,omitempty"`
}
