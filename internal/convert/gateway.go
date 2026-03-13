package convert

import (
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/TheKodeToad/fine/internal/fluxer"
)

// ReadyEventToDiscord converts a Fluxer gateway ready event to a Discord one.
// The ResumeGatewayURL had better be non-nil, you have been warned.
func ReadyEventToDiscord(data fluxer.ReadyEvent) discord.ReadyEvent {
	return discord.ReadyEvent{
		V:                data.Version,
		User:             UserPrivateToDiscord(data.User),
		Guilds:           data.Guilds,
		SessionID:        data.SessionID,
		ResumeGatewayURL: *data.ResumeGatewayURL,
		Application: discord.ReadyEventApplication{
			ID: data.User.ID,
			// NOTE: if Fluxer adds any application flags this may no longer be accurate
			Flags: ApplicationFlagsToDiscord(),
		},
	}
}

func MessageCreateEventToDiscord(event fluxer.MessageCreateEvent) discord.MessageCreateEvent {
	return discord.MessageCreateEvent{
		Message: MessageToDiscord(event.Message),
		GuildID: event.GuildID,
	}
}
