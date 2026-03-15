package discord

import "github.com/disgoorg/snowflake/v2"

type Emoji struct {
	ID       *snowflake.ID `json:"id"`
	Name     *string       `json:"name"`
	Animated *bool         `json:"animated,omitempty"`
}

type Sticker struct {
	ID          snowflake.ID  `json:"id"`
	Name        string        `json:"name"`
	Description *string       `json:"description"`
	Tags        string        `json:"tags"`
	Type        int           `json:"type"`
	FormatType  int           `json:"format_type"`
	Available   *bool         `json:"available,omitempty"`
	GuildID     *snowflake.ID `json:"guild_id,omitempty"`
	User        *User         `json:"user,omitempty"`
}
