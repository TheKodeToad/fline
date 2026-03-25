package convert

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
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
	var bot *discord.User
	var description string
	if app.Bot != nil {
		bot = misc.New(BotToDiscord(*app.Bot))

		if app.Bot.Bio != nil {
			description = *app.Bot.Bio
		}
	}

	return discord.Application{
		ID:                  app.ID,
		Name:                app.Name,
		Icon:                app.Icon,
		Description:         description,
		BotPublic:           app.BotPublic,
		BotRequireCodeGrant: app.BotRequireCodeGrant,
		Bot:                 bot,
		Flags:               ApplicationFlagsToDiscord(),
		Owner: &discord.User{
			ID:            fluxer.DeletedUserID,
			Username:      fluxer.DeletedUserUsername,
			Discriminator: fluxer.DeletedUserDiscrim,
			System:        misc.New(true),
		},
	}
}
