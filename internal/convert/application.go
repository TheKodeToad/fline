package convert

import (
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/TheKodeToad/fine/internal/fluxer"
	"github.com/TheKodeToad/fine/internal/misc"
)

func BotToDiscord(bot fluxer.Bot) discord.User {
	return discord.User{
		ID:            bot.ID,
		Username:      bot.Username,
		Discriminator: bot.Discriminator,
		Avatar:        bot.Avatar,
		Bot:           misc.New(true),
		Banner:        bot.Banner,
		Flags:         misc.New(UserFlagsToDiscord(bot.Flags)),
	}
}

func ApplicationFlagsToDiscord() discord.ApplicationFlags {
	// NOTE: Fluxer has no flags of its own and these are the flags that should always appear to be present
	return discord.ApplicationFlagGatewayPresence |
		discord.ApplicationFlagGatewayGuildMembers |
		discord.ApplicationFlagGatewayMessageContent
}

func ApplicationToDiscord(app fluxer.Application) discord.Application {
	result := discord.Application{
		ID:                  app.ID,
		Name:                app.Name,
		Icon:                app.Icon,
		BotPublic:           app.BotPublic,
		BotRequireCodeGrant: app.BotRequireCodeGrant,
		Flags:               ApplicationFlagsToDiscord(),
		Owner: &discord.User{
			ID:            1130650140672000000,
			Username:      "Hampus",
			Discriminator: "0001",
			GlobalName:    misc.New("Hampus"),
			Avatar:        misc.New("c4c8b5d0"),
			Flags:         misc.New(discord.UserFlags(1)),
		},
	}

	if app.Bot != nil {
		result.Bot = misc.New(BotToDiscord(*app.Bot))

		if app.Bot.Bio != nil {
			result.Description = *app.Bot.Bio
		}
	}

	return result
}
