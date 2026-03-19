package fluxer

import "github.com/disgoorg/snowflake/v2"

type Sticker struct {
	ID          snowflake.ID `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Tags        []string     `json:"tags"`
	Animated    bool         `json:"animated"`
}
