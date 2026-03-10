package discord

import "github.com/disgoorg/snowflake/v2"

type UserFlags uint

const (
	UserFlagStaff           UserFlags = 1 << 0
	UserFlagBugHunterLevel1 UserFlags = 1 << 3
	UserFlagBugHunterLevel2 UserFlags = 1 << 14
)

type PremiumType uint

const (
	PremiumTypeNone  PremiumType = 0
	PremiumTypeNitro PremiumType = 2
)

type User struct {
	ID            snowflake.ID `json:"id"`
	Username      string       `json:"username"`
	Discriminator string       `json:"discriminator"`
	GlobalName    *string      `json:"global_name"`
	Avatar        *string      `json:"avatar"`
	Bot           *bool        `json:"bot,omitempty"`
	System        *bool        `json:"system,omitempty"`
	MFAEnabled    *bool        `json:"mfa_enabled,omitempty"`
	Banner        *string      `json:"banner,omitempty"`
	AccentColor   *uint        `json:"accent_color,omitempty"`
	Locale        *string      `json:"locale,omitempty"`
	Verified      *bool        `json:"verified,omitempty"`
	Email         *string      `json:"email,omitempty"`
	Flags         *UserFlags   `json:"flags,omitempty"`
	PremiumType   PremiumType  `json:"premium_type,omitempty"`
	PublicFlags   *UserFlags   `json:"public_flags,omitempty"`
}
