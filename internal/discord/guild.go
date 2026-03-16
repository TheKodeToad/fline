package discord

import (
	"github.com/disgoorg/snowflake/v2"
)

type Guild struct {
	ID                          snowflake.ID  `json:"id"`
	Name                        string        `json:"name"`
	Icon                        *string       `json:"icon"`
	Splash                      *string       `json:"splash"`
	Owner                       *bool         `json:"owner,omitempty"`
	OwnerID                     snowflake.ID  `json:"owner_id"`
	Permissions                 *Permissions  `json:"permissions,omitempty"`
	AFKChannelID                *snowflake.ID `json:"afk_channel_id"`
	AFKTimeout                  int           `json:"afk_timeout"`
	VerificationLevel           int           `json:"verification_level"`
	DefaultMessageNotifications int           `json:"default_message_notifications"`
	ExplicitContentFilter       int           `json:"explicit_content_filter"`
	Features                    []string      `json:"features"`
	MFALevel                    int           `json:"mfa_level"`
	ApplicationID               *snowflake.ID `json:"application_id"`
	SystemChannelID             *snowflake.ID `json:"system_channel_id"`
	SystemChannelFlags          uint          `json:"system_channel_flags"`
	RulesChannelID              *snowflake.ID `json:"rules_channel_id"`
	VanityURLCode               *string       `json:"vanity_url_code"`
	Description                 *string       `json:"description"`
	Banner                      *string       `json:"banner"`
	PreferredLocale             string        `json:"preferred_locale"`
	ApproximateMemberCount      *int          `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    *int          `json:"approximate_presence_count,omitempty"`
	NSFWLevel                   int           `json:"nsfw_level"`
	Roles                       []Role        `json:"roles"`
	Emojis                      []Emoji       `json:"emojis"`
	Stickers                    []Sticker     `json:"stickers,omitzero"`
}

type GuildMember struct {
	User                       *User          `json:"user,omitempty"`
	Nick                       *string        `json:"nick"`
	Avatar                     *string        `json:"avatar"`
	Banner                     *string        `json:"banner"`
	Roles                      []snowflake.ID `json:"roles"`
	JoinedAt                   *string        `json:"joined_at"`
	Deaf                       bool           `json:"deaf"`
	Mute                       bool           `json:"mute"`
	Flags                      uint           `json:"flags"`
	CommunicationDisabledUntil *string        `json:"communication_disabled_until"`
}

type GuildBanCreate struct {
	DeleteMessageDays    *int `json:"delete_message_days,omitempty"`
	DeleteMessageSeconds *int `json:"delete_message_seconds,omitempty"`
}
