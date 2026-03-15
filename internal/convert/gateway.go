package convert

import (
	"fmt"

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

func GuildCreateEventToDiscord(event fluxer.GuildCreateEvent) discord.GuildCreateEvent {
	members := make([]discord.GuildMember, 0, len(event.Members))
	for _, member := range event.Members {
		members = append(members, GuildMemberToDiscord(member))
	}

	channels := make([]discord.Channel, 0, len(event.Channels))
	for _, channel := range event.Channels {
		conv, ok := ChannelToDiscord(channel)
		if !ok {
			continue
		}

		channels = append(channels, conv)
	}

	guild := GuildToDiscord(event.Properties)

	guild.Roles = make([]discord.Role, 0, len(event.Roles))
	for _, role := range event.Roles {
		guild.Roles = append(guild.Roles, RoleToDiscord(role))
	}

	guild.Emojis = event.Emojis
	guild.Stickers = event.Stickers

	return discord.GuildCreateEvent{
		Guild:       guild,
		JoinedAt:    event.JoinedAt,
		Large:       event.Large,
		Unavailable: event.Unavailable,
		MemberCount: event.MemberCount,
		Members:     members,
		Channels:    channels,
	}
}

func MessageCreateEventToDiscord(event fluxer.MessageCreateEvent) discord.MessageCreateEvent {
	return discord.MessageCreateEvent{
		Message: MessageToDiscord(event.Message),
		GuildID: event.GuildID,
	}
}
