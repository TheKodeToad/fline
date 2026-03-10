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

func ApplicationToDiscord(app fluxer.Application) discord.Application {
	flags := discord.ApplicationGatewayPresence |
		discord.ApplicationGatewayGuildMembers |
		discord.ApplicationGatewayMessageContent

	result := discord.Application{
		ID:                   app.ID,
		Name:                 app.Name,
		Icon:                 app.Icon,
		BotPublic:            app.BotPublic,
		BotRequireCodeGrant: app.BotRequireCodeGrant,
		Flags:                flags,
	}

	if app.Bot != nil {
		result.Bot = misc.New(BotToDiscord(*app.Bot))

		if app.Bot.Bio != nil {
			result.Description = *app.Bot.Bio
		}
	}

	return result
}
