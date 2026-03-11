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

// IntrinsicFluxerApplicationFlags contains the Discord flags apps on Fluxer intrinsically have the equivilent of.
// NOTE: Fluxer is yet to add any application flags of its own as far as I can tell, so converting Fluxer flags to Discord is as simple as just using these.
var IntrinsicFluxerApplicationFlags = discord.ApplicationGatewayPresence |
	discord.ApplicationGatewayGuildMembers |
	discord.ApplicationGatewayMessageContent

func ApplicationToDiscord(app fluxer.Application) discord.Application {
	result := discord.Application{
		ID:                  app.ID,
		Name:                app.Name,
		Icon:                app.Icon,
		BotPublic:           app.BotPublic,
		BotRequireCodeGrant: app.BotRequireCodeGrant,
		Flags:               IntrinsicFluxerApplicationFlags,
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
