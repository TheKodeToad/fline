package convert

import (
	"math"

	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
)

func GuildToDiscord(guild fluxer.Guild) discord.Guild {
	return discord.Guild{
		ID:                          guild.ID,
		Name:                        guild.Name,
		Icon:                        guild.Splash,
		Owner:                       guild.Owner,
		OwnerID:                     guild.OwnerID,
		Permissions:                 guild.Permissions,
		AFKChannelID:                guild.AFKChannelID,
		AFKTimeout:                  guild.AFKTimeout,
		VerificationLevel:           guild.VerificationLevel,
		DefaultMessageNotifications: guild.DefaultMessageNotifications,
		ExplicitContentFilter:       guild.ExplicitContentFilter,
		Features:                    guild.Features,
		MFALevel:                    guild.MFALevel,
		SystemChannelID:             guild.SystemChannelID,
		SystemChannelFlags:          guild.SystemChannelFlags,
		RulesChannelID:              guild.RulesChannelID,
		VanityURLCode:               guild.VanityURLCode,
		Banner:                      guild.Banner,
		ApproximateMemberCount:      guild.ApproximateMemberCount,
		ApproximatePresenceCount:    guild.ApproximatePresenceCount,
		NSFWLevel:                   guild.NSFWLevel,
		Roles:                       []discord.Role{},
		Emojis:                      []discord.Emoji{},
		Stickers:                    []discord.Sticker{},
	}
}

func GuildMemberToDiscord(member fluxer.GuildMember) discord.GuildMember {
	var user *discord.User
	if member.User != nil {
		user = misc.New(UserPartialToDiscord(*member.User))
	}

	return discord.GuildMember{
		User:                       user,
		Nick:                       member.Nick,
		Avatar:                     member.Avatar,
		Banner:                     member.Banner,
		Roles:                      member.Roles,
		JoinedAt:                   member.JoinedAt,
		Deaf:                       member.Deaf,
		Mute:                       member.Mute,
		CommunicationDisabledUntil: member.CommunicationDisabledUntil,
	}
}

func GuildBanCreateToFluxer(create discord.GuildBanCreate) fluxer.GuildBanCreate {
	var deleteMessageDays int
	// TODO: check which takes precidence on Discord and update it accordingly if necessary
	if create.DeleteMessageSeconds != nil {
		const secsInDay = 60 * 60 * 24

		fractionalDays := float64(*create.DeleteMessageSeconds) / secsInDay
		deleteMessageDays = int(math.Round(float64(fractionalDays)))
	} else if create.DeleteMessageDays != nil {
		deleteMessageDays = *create.DeleteMessageDays
	}

	return fluxer.GuildBanCreate{DeleteMessageDays: deleteMessageDays}
}
