package convert

import (
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/TheKodeToad/fine/internal/fluxer"
	"github.com/TheKodeToad/fine/internal/misc"
)

func UserFlagsToDiscord(flags fluxer.UserFlags) discord.UserFlags {
	var result discord.UserFlags

	if flags&fluxer.UserStaff != 0 {
		result |= discord.UserStaff
	}

	if flags&fluxer.UserBugHunter != 0 {
		result |= discord.UserBugHunterLevel1
	}

	return result
}

func PremiumTypeToDiscord(kind fluxer.PremiumType) discord.PremiumType {
	if kind == fluxer.PremiumTypeNone {
		return discord.PremiumTypeNone
	} else {
		return discord.PremiumTypeNitro
	}
}

func PartialUserToDiscord(user fluxer.PartialUser) discord.User {
	return discord.User{
		ID:            user.ID,
		Username:      user.Username,
		Discriminator: user.Discriminator,
		GlobalName:    user.GlobalName,
		Avatar:        user.Avatar,
		Bot:           user.Bot,
		System:        user.System,
		Flags:         misc.New(UserFlagsToDiscord(user.Flags)),
	}
}

func PrivateUserToDiscord(user fluxer.PrivateUser) discord.User {
	result := PartialUserToDiscord(user.PartialUser)

	result.AccentColor = user.AccentColor
	result.MFAEnabled = &user.MFAEnabled
	result.Banner = user.Banner

	return result
}
