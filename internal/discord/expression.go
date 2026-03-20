package discord

import "github.com/disgoorg/snowflake/v2"

type Emoji struct {
	ID            *snowflake.ID  `json:"id"`
	Name          *string        `json:"name"`
	Roles         []snowflake.ID `json:"roles,omitzero"`
	RequireColons *bool          `json:"require_colons,omitempty"`
	Managed       *bool          `json:"managed,omitempty"`
	Animated      *bool          `json:"animated,omitempty"`
	Available     *bool          `json:"available,omitempty"`
}

type StickerType uint

const (
	StickerTypeGuild StickerType = 2
)

type StickerFormat uint

const (
	StickerFormatPNG StickerFormat = 1
	StickerFormatGIF StickerFormat = 4
)

type Sticker struct {
	ID          snowflake.ID  `json:"id"`
	Name        string        `json:"name"`
	Description *string       `json:"description"`
	Tags        string        `json:"tags"`
	Type        StickerType   `json:"type"`
	FormatType  StickerFormat `json:"format_type"`
	Available   *bool         `json:"available,omitempty"`
	GuildID     *snowflake.ID `json:"guild_id,omitempty"`
	User        *User         `json:"user,omitempty"`
}
