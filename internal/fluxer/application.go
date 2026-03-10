package fluxer

import "github.com/disgoorg/snowflake/v2"

type Bot struct {
	ID            snowflake.ID `json:"id"`
	Username      string       `json:"username"`
	Discriminator string       `json:"discriminator"`
	Avatar        *string      `json:"avatar"`
	Banner        *string      `json:"banner"`
	Bio           *string      `json:"bio"`
	Token         *string      `json:"token"`
	MFAEnabled    bool         `json:"mfa_enabled"`
	Flags         UserFlags    `json:"flags"`
}

type ApplicationFlags uint

type Application struct {
	ID                  snowflake.ID     `json:"id"`
	Name                string           `json:"name"`
	Icon                *string          `json:"icon"`
	Description         *string          `json:"description"`
	BotPublic           bool             `json:"bot_public"`
	BotRequireCodeGrant bool             `json:"bot_require_code_grant"`
	Flags               ApplicationFlags `json:"flags"`
	Bot                 *Bot             `json:"bot"`
}
