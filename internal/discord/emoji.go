package discord

import "github.com/disgoorg/snowflake/v2"

type Emoji struct {
	ID *snowflake.ID `json:"id"`
	Name *string `json:"name"`
	Animated *bool `json:"animated,omitempty"`
}
