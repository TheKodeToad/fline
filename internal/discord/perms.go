package discord

import (
	"github.com/TheKodeToad/fline/internal/flags"
	"github.com/disgoorg/snowflake/v2"
)

// NOTE: future proofing this for if/when Discord has more than 64 permissions
type Permissions = flags.Set

var (
	PermCreateInstantInvite = flags.Bit(0)
	PermKickMembers         = flags.Bit(1)
	PermBanMembers          = flags.Bit(2)
	PermAdministrator       = flags.Bit(3)
	PermManageChannels      = flags.Bit(4)
	PermManageGuild         = flags.Bit(5)
	PermAddReactions        = flags.Bit(6)
	PermViewAuditLog        = flags.Bit(7)
	PermPrioritySpeaker     = flags.Bit(8)
	PermStream              = flags.Bit(9)
	PermViewChannel         = flags.Bit(10)
	PermSendMessages        = flags.Bit(11)
	PermSendTTSMessages     = flags.Bit(12)
	PermManageMessages      = flags.Bit(13)
	PermEmbedLinks          = flags.Bit(14)
	PermAttachFiles         = flags.Bit(15)
	PermReadMessageHistory  = flags.Bit(16)
	PermMentionEveryone     = flags.Bit(17)
	PermUseExternalEmojis   = flags.Bit(18)
	PermConnect             = flags.Bit(20)
	PermSpeak               = flags.Bit(21)
	PermMuteMembers         = flags.Bit(22)
	PermDeafenMembers       = flags.Bit(23)
	PermMoveMembers         = flags.Bit(24)
	PermUseVAD              = flags.Bit(25)
	PermChangeNickname      = flags.Bit(26)
	PermManageNicknames     = flags.Bit(27)
	PermManageRoles         = flags.Bit(28)
	PermManageWebhooks      = flags.Bit(29)
	PermManageExpressions   = flags.Bit(30)
	PermUseExternalStickers = flags.Bit(37)
	PermModerateMembers     = flags.Bit(40)
	PermCreateExpressions   = flags.Bit(43)
	PermPinMembers          = flags.Bit(51)
	PermBypassSlowmode      = flags.Bit(52)
)

type RoleColors struct {
	PrimaryColor   int  `json:"primary_color"`
	SecondaryColor *int `json:"secondary_color"`
	TertiaryColor  *int `json:"tertiary_color"`
}

type Role struct {
	ID           snowflake.ID `json:"id"`
	Name         string       `json:"name"`
	Color        int          `json:"color"`
	Colors       RoleColors   `json:"colors"`
	Hoist        bool         `json:"hoist"`
	UnicodeEmoji *string      `json:"unicode_emoji,omitempty"`
	Position     int          `json:"position"`
	Permissions  Permissions  `json:"permissions"`
	Managed      bool         `json:"managed"`
	Mentionable  bool         `json:"mentionable"`
	Flags        uint         `json:"flags"`
}
