package convert

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
)

func UserFlagsToDiscord(flags fluxer.UserFlags) discord.UserFlags {
	var result discord.UserFlags

	if flags&fluxer.UserFlagStaff != 0 {
		result |= discord.UserFlagStaff
	}

	if flags&fluxer.UserFlagBugHunter != 0 {
		result |= discord.UserFlagBugHunterLevel1
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

func UserPartialToDiscord(user fluxer.UserPartial) discord.User {
	return discord.User{
		ID:            user.ID,
		Username:      user.Username,
		Discriminator: user.Discriminator,
		GlobalName:    user.GlobalName,
		Avatar:        user.Avatar,
		Bot:           user.Bot,
		System:        user.System,
		Flags:         misc.New(UserFlagsToDiscord(user.Flags)),
		PublicFlags:   misc.New(UserFlagsToDiscord(user.Flags)),
	}
}

func UserPrivateToDiscord(user fluxer.UserPrivate) discord.User {
	result := UserPartialToDiscord(user.UserPartial)

	result.AccentColor = user.AccentColor
	result.MFAEnabled = &user.MFAEnabled
	result.Banner = user.Banner

	return result
}
