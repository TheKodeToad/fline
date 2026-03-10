package convert

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/flags"
	"github.com/TheKodeToad/fline/internal/fluxer"
)

var permsCommonSubset = flags.New(
	discord.PermCreateInstantInvite,
	discord.PermKickMembers,
	discord.PermBanMembers,
	discord.PermAdministrator,
	discord.PermManageChannels,
	discord.PermManageGuild,
	discord.PermAddReactions,
	discord.PermViewAuditLog,
	discord.PermPrioritySpeaker,
	discord.PermStream,
	discord.PermViewChannel,
	discord.PermSendMessages,
	discord.PermSendTTSMessages,
	discord.PermManageMessages,
	discord.PermEmbedLinks,
	discord.PermAttachFiles,
	discord.PermReadMessageHistory,
	discord.PermMentionEveryone,
	discord.PermUseExternalEmojis,
	discord.PermConnect,
	discord.PermSpeak,
	discord.PermMuteMembers,
	discord.PermDeafenMembers,
	discord.PermMoveMembers,
	discord.PermUseVAD,
	discord.PermChangeNickname,
	discord.PermManageNicknames,
	discord.PermManageRoles,
	discord.PermManageWebhooks,
	discord.PermManageExpressions,
	discord.PermUseExternalStickers,
	discord.PermModerateMembers,
	discord.PermCreateExpressions,
	discord.PermPinMembers,
	discord.PermBypassSlowmode,
)

func PermissionsToDiscord(perms discord.Permissions) discord.Permissions {
	return perms.Intersection(permsCommonSubset)
}

func RoleToDiscord(role fluxer.Role) discord.Role {
	return discord.Role{
		ID: role.ID,
		Name: role.Name,
		Color: role.Color,
		Colors: discord.RoleColors{
			PrimaryColor: role.Color,
		},
		Hoist: role.Hoist,
		UnicodeEmoji: role.UnicodeEmoji,
		Position: role.Position,
		Permissions: role.Permissions,
		Mentionable: role.Mentionable,
	}
}
