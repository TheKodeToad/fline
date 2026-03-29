package discord

import (
	"encoding/json"
	"fmt"

	"github.com/TheKodeToad/fline/internal/misc"
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

type GuildMemberUpdate struct {
	Nick                       *string
	Roles                      []snowflake.ID
	Mute                       *bool
	Deaf                       *bool
	ChannelID                  *snowflake.ID
	CommunicationDisabledUntil *string

	ClearChannel bool
}

func (u *GuildMemberUpdate) UnmarshalJSON(data []byte) error {
	var raw struct {
		Nick                       json.RawMessage `json:"nick"`
		Roles                      json.RawMessage `json:"roles"`
		Mute                       json.RawMessage `json:"mute"`
		Deaf                       json.RawMessage `json:"deaf"`
		ChannelID                  json.RawMessage `json:"channel_id"`
		CommunicationDisabledUntil json.RawMessage `json:"communication_disabled_until"`
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	// NOTE: we unfortunately need this manual parsing to distinguish between null and undefined
	// fortunately for most of the fields, null is semantically equivilent to specifying the default value

	if string(raw.Nick) == "null" {
		u.Nick = misc.New("")
	} else if raw.Nick != nil {
		err := json.Unmarshal(raw.Nick, &u.Nick)
		if err != nil {
			return fmt.Errorf("unmarshalling into GuildMemberUpdate.Nick: %w", err)
		}
	}

	if string(raw.Roles) == "null" {
		u.Roles = []snowflake.ID{}
	} else if raw.Roles != nil {
		err := json.Unmarshal(raw.Roles, &u.Roles)
		if err != nil {
			return fmt.Errorf("unmarshalling into GuildMemberUpdate.Roles: %w", err)
		}
	}

	if string(raw.Mute) == "null" {
		u.Mute = misc.New(false)
	} else if raw.Mute != nil {
		err := json.Unmarshal(raw.Mute, &u.Mute)
		if err != nil {
			return fmt.Errorf("unmarshalling into GuildMemberUpdate.Mute: %w", err)
		}
	}

	if string(raw.Deaf) == "null" {
		u.Deaf = misc.New(false)
	} else if raw.Deaf != nil {
		err := json.Unmarshal(raw.Deaf, &u.Deaf)
		if err != nil {
			return fmt.Errorf("unmarshalling into GuildMemberUpdate.Deaf: %w", err)
		}
	}

	if string(raw.ChannelID) == "null" {
		u.ClearChannel = true
	} else if raw.ChannelID != nil {
		u.ClearChannel = false

		err := json.Unmarshal(raw.ChannelID, &u.ChannelID)
		if err != nil {
			return fmt.Errorf("unmarshalling into GuildMemberUpdate.ChannelID: %w", err)
		}
	}

	if string(raw.CommunicationDisabledUntil) == "null" {
		u.CommunicationDisabledUntil = misc.New("")
	} else if raw.CommunicationDisabledUntil != nil {
		err := json.Unmarshal(raw.CommunicationDisabledUntil, &u.CommunicationDisabledUntil)
		if err != nil {
			return fmt.Errorf("unmarshalling into GuildMemberUpdate.CommunicationDisabledUntil: %w", err)
		}
	}

	return nil
}

type GuildBanCreate struct {
	DeleteMessageDays    *int `json:"delete_message_days,omitempty"`
	DeleteMessageSeconds *int `json:"delete_message_seconds,omitempty"`
}
