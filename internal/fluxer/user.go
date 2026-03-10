package fluxer

import "github.com/disgoorg/snowflake/v2"

type UserFlags uint

const (
	UserFlagStaff     UserFlags = 1 << 0
	UserFlagCTPMember UserFlags = 1 << 1
	UserFlagPartner   UserFlags = 1 << 2
	UserFlagBugHunter UserFlags = 1 << 3
)

type PremiumType uint

const (
	PremiumTypeNone PremiumType = 0
)

type UserPartial struct {
	ID            snowflake.ID `json:"id"`
	Username      string       `json:"username"`
	Discriminator string       `json:"discriminator"`
	GlobalName    *string      `json:"global_name"`
	Avatar        *string      `json:"avatar"`
	Bot           *bool        `json:"bot"`
	System        *bool        `json:"system"`
	Flags         UserFlags    `json:"flags"`
}

type UserPrivate struct {
	UserPartial
	AccentColor *uint   `json:"accent_color"`
	MFAEnabled  bool    `json:"mfa_enabled"`
	Banner      *string `json:"banner"`
}
