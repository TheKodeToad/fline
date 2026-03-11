package convert

import (
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/TheKodeToad/fine/internal/fluxer"
)

// ReadyEventToDiscord converts a Fluxer gateway ready event to a Discord one.
// The ResumeGatewayURL had better be not-nil, you have been warned.
func ReadyEventToDiscord(data fluxer.ReadyEvent) discord.ReadyEvent {
	return discord.ReadyEvent{
		V:                data.Version,
		User:             PrivateUserToDiscord(data.User),
		Guilds:           data.Guilds,
		SessionID:        data.SessionID,
		ResumeGatewayURL: *data.ResumeGatewayURL,
		Application: discord.ReadyEventApplication{
			ID:    data.User.ID,
			Flags: IntrinsicFluxerApplicationFlags,
		},
	}
}
