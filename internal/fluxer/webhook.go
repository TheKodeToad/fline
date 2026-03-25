package fluxer

import "github.com/disgoorg/snowflake/v2"

type Webhook struct {
	ID        snowflake.ID `json:"id"`
	GuildID   snowflake.ID `json:"guild_id"`
	ChannelID snowflake.ID `json:"channel_id"`
	User      UserPartial  `json:"user"`
	Name      string       `json:"name"`
	Avatar    *string      `json:"avatar"`
	Token     string       `json:"token,omitempty"`
}
