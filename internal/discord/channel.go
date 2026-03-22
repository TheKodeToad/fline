package discord

import "github.com/disgoorg/snowflake/v2"

type ChannelType uint

var (
	ChannelTypeGuildText     ChannelType = 0
	ChannelTypeDM            ChannelType = 1
	ChannelTypeGuildVoice    ChannelType = 2
	ChannelTypeGroupDM       ChannelType = 3
	ChannelTypeGuildCategory ChannelType = 4
)

type PermissionOverwrite struct {
	ID    snowflake.ID `json:"id"`
	Type  uint         `json:"type"`
	Allow Permissions  `json:"allow"`
	Deny  Permissions  `json:"deny"`
}

type Channel struct {
	ID                   snowflake.ID          `json:"id"`
	Type                 ChannelType           `json:"type"`
	GuildID              *snowflake.ID         `json:"guild_id,omitempty"`
	Position             *int                  `json:"position,omitempty"`
	PermissionOverwrites []PermissionOverwrite `json:"permission_overwrites,omitzero"`
	Name                 *string               `json:"name"`
	Topic                *string               `json:"topic"`
	NSFW                 *bool                 `json:"nsfw"`
	LastMessageID        *snowflake.ID         `json:"last_message_id"`
	Bitrate              *int                  `json:"bitrate,omitempty"`
	UserLimit            *int                  `json:"user_limit,omitempty"`
	RateLimitPerUser     *int                  `json:"rate_limit_per_user,omitempty"`
	Recipients           []User                `json:"recipients,omitempty"`
	Icon                 *string               `json:"icon"`
	OwnerID              *snowflake.ID         `json:"owner_id,omitempty"`
	ApplicationID        *snowflake.ID         `json:"application_id,omitempty"`
	ParentID             *snowflake.ID         `json:"parent_id"`
	LastPinTimestamp     *string               `json:"last_pin_timestamp,omitempty"`
	RTCRegion            *string               `json:"rtc_region"`
}

type WebhookType uint

const (
	WebhookTypeIncoming WebhookType = 1
)

type WebhookCreate struct {
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
}

type Webhook struct {
	ID        snowflake.ID  `json:"id"`
	Type      WebhookType   `json:"type"`
	GuildID   *snowflake.ID `json:"guild_id"`
	ChannelID *snowflake.ID `json:"channel_id"`
	User      *User         `json:"user,omitempty"`
	Name      *string       `json:"name"`
	Avatar    *string       `json:"avatar"`
	Token     *string       `json:"token,omitempty"`
}
