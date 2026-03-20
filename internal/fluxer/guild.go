package fluxer

import (
	"encoding/json"
	"fmt"

	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/disgoorg/snowflake/v2"
)

type Guild struct {
	ID                          snowflake.ID         `json:"id"`
	Name                        string               `json:"name"`
	Icon                        *string              `json:"icon"`
	Splash                      *string              `json:"splash"`
	Owner                       *bool                `json:"owner,omitempty"`
	OwnerID                     snowflake.ID         `json:"owner_id"`
	Permissions                 *discord.Permissions `json:"permissions,omitempty"`
	AFKChannelID                *snowflake.ID        `json:"afk_channel_id"`
	AFKTimeout                  int                  `json:"afk_timeout"`
	VerificationLevel           int                  `json:"verification_level"`
	DefaultMessageNotifications int                  `json:"default_message_notifications"`
	ExplicitContentFilter       int                  `json:"explicit_content_filter"`
	Features                    []string             `json:"features"`
	MFALevel                    int                  `json:"mfa_level"`
	ApplicationID               *snowflake.ID        `json:"application_id"`
	SystemChannelID             *snowflake.ID        `json:"system_channel_id"`
	SystemChannelFlags          uint                 `json:"system_channel_flags"`
	RulesChannelID              *snowflake.ID        `json:"rules_channel_id"`
	VanityURLCode               *string              `json:"vanity_url_code"`
	Description                 *string              `json:"description"`
	Banner                      *string              `json:"banner"`
	PreferredLocale             string               `json:"preferred_locale"`
	ApproximateMemberCount      *int                 `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount    *int                 `json:"approximate_presence_count,omitempty"`
	NSFWLevel                   int                  `json:"nsfw_level"`
}

type GuildMember struct {
	User                       *UserPartial   `json:"user,omitempty"`
	Nick                       *string        `json:"nick"`
	Avatar                     *string        `json:"avatar"`
	Banner                     *string        `json:"banner"`
	Roles                      []snowflake.ID `json:"roles"`
	JoinedAt                   *string        `json:"joined_at"`
	Deaf                       bool           `json:"deaf"`
	Mute                       bool           `json:"mute"`
	CommunicationDisabledUntil *string        `json:"communication_disabled_until"`
}

type GuildMemberUpdate struct {
	Nick                       *string
	Roles                      []snowflake.ID
	Mute                       *bool
	Deaf                       *bool
	ChannelID                  *snowflake.ID
	CommunicationDisabledUntil *string

	// NOTE: Fluxer's semantics are different
	// an empty nick/communication_disabled_until string is not semantically equivilent to null
	ClearNick    bool
	ClearChannel bool
	ClearTimeout bool
}

func (u GuildMemberUpdate) MarshalJSON() ([]byte, error) {
	var raw struct {
		Nick                       json.RawMessage `json:"nick,omitempty"`
		Roles                      json.RawMessage `json:"roles,omitempty"`
		Mute                       json.RawMessage `json:"mute,omitempty"`
		Deaf                       json.RawMessage `json:"deaf,omitempty"`
		ChannelID                  json.RawMessage `json:"channel_id,omitempty"`
		CommunicationDisabledUntil json.RawMessage `json:"communication_disabled_until,omitempty"`
	}

	if u.ClearNick {
		raw.Nick = []byte("null")
	} else if u.Nick != nil {
		data, err := json.Marshal(u.Nick)
		if err != nil {
			return nil, fmt.Errorf("marshalling GuildMemberUpdate.Nick: %w", err)
		}

		raw.Nick = data
	}

	if u.Roles != nil {
		data, err := json.Marshal(u.Roles)
		if err != nil {
			return nil, fmt.Errorf("marshalling GuildMemberUpdate.Nick: %w", err)
		}

		raw.Roles = data
	}

	if u.Mute != nil {
		data, err := json.Marshal(u.Mute)
		if err != nil {
			return nil, fmt.Errorf("marshalling GuildMemberUpdate.Mute: %w", err)
		}

		raw.Mute = data
	}

	if u.Deaf != nil {
		data, err := json.Marshal(u.Deaf)
		if err != nil {
			return nil, fmt.Errorf("marshalling GuildMemberUpdate.Deaf: %w", err)
		}

		raw.Deaf = data
	}

	if u.ClearChannel {
		raw.ChannelID = []byte("null")
	} else if u.ChannelID != nil {
		data, err := json.Marshal(u.ChannelID)
		if err != nil {
			return nil, fmt.Errorf("marshalling GuildMemberUpdate.ChannelID: %w", err)
		}

		raw.ChannelID = data
	}

	if u.ClearTimeout {
		raw.CommunicationDisabledUntil = []byte("null")
	} else if u.CommunicationDisabledUntil != nil {
		data, err := json.Marshal(u.CommunicationDisabledUntil)
		if err != nil {
			return nil, fmt.Errorf("marshalling GuildMemberUpdate.CommunicationDisabledUntil: %w", err)
		}

		raw.CommunicationDisabledUntil = data
	}

	return json.Marshal(raw)
}

type GuildBanCreate struct {
	DeleteMessageDays int    `json:"delete_message_days"`
	Reason            string `json:"reason"`
}
