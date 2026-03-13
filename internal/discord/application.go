package discord

import "github.com/disgoorg/snowflake/v2"

type EventWebhookStatus uint

const (
	EventWebhookStatusDisabled EventWebhookStatus = 0
)

type ApplicationFlags uint

var (
	ApplicationFlagGatewayPresence       ApplicationFlags = 1 << 12
	ApplicationFlagGatewayGuildMembers   ApplicationFlags = 1 << 14
	ApplicationFlagGatewayMessageContent ApplicationFlags = 1 << 18
)

type Application struct {
	ID                  snowflake.ID       `json:"id"`
	Name                string             `json:"name"`
	Icon                *string            `json:"icon"`
	Description         string             `json:"description"`
	BotPublic           bool               `json:"bot_public"`
	BotRequireCodeGrant bool               `json:"bot_require_code_grant"`
	Bot                 *User              `json:"bot"`
	Owner               *User              `json:"owner,omitempty"`
	Flags               ApplicationFlags   `json:"flags"`
	VerifyKey           string             `json:"verify_key"`
	RedirectURIs        *[]string          `json:"redirect_uris,omitempty"`
	EventWebhookStatus  EventWebhookStatus `json:"event_webhook_status"`
}
