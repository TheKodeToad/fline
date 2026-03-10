package fluxer

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/disgoorg/snowflake/v2"
)

type RoleColors struct {
	PrimaryColor   int  `json:"primary_color"`
	SecondaryColor *int `json:"secondary_color"`
	TertiaryColor  *int `json:"tertiary_color"`
}

type Role struct {
	ID           snowflake.ID        `json:"id"`
	Name         string              `json:"name"`
	Color        int                 `json:"color"`
	Hoist        bool                `json:"hoist"`
	UnicodeEmoji *string             `json:"unicode_emoji,omitempty"`
	Position     int                 `json:"position"`
	Permissions  discord.Permissions `json:"permissions"`
	Mentionable  bool                `json:"mentionable"`
}
