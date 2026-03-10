package fluxer

import "github.com/disgoorg/snowflake/v2"

type UserFlags uint

const (
	UserStaff     UserFlags = 1 << 0
	UserCTPMember UserFlags = 1 << 1
	UserPartner   UserFlags = 1 << 2
	UserBugHunter UserFlags = 1 << 3
)

type PremiumType uint

const (
	PremiumTypeNone  PremiumType = 0
)

type PartialUser struct {
	ID            snowflake.ID `json:"id"`
	Username      string       `json:"username"`
	Discriminator string       `json:"discriminator"`
	GlobalName    *string      `json:"global_name"`
	Avatar        *string      `json:"avatar"`
	Bot           *bool        `json:"bot"`
	System        *bool        `json:"system"`
	Flags         UserFlags    `json:"flags"`
}

type PrivateUser struct {
	PartialUser
	AccentColor *int    `json:"accent_color"`
	MFAEnabled  bool    `json:"mfa_enabled"`
	Banner      *string `json:"banner"`
}
