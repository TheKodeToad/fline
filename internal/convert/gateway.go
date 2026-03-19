package convert

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
)

func PresenceUpdatePayloadToFluxer(payload discord.UpdatePresencePayload) fluxer.UpdatePresencePayload {
	var customStatus *fluxer.CustomStatus
	for _, activity := range payload.Activities {
		if activity.Type == discord.ActivityCustom {
			customStatus = &fluxer.CustomStatus{Text: activity.State}
			break
		}
	}

	return fluxer.UpdatePresencePayload{
		Status:       payload.Status,
		AFK:          payload.AFK,
		CustomStatus: customStatus,
	}
}

func IdentifyPayloadToFluxer(payload discord.IdentifyPayload) fluxer.IdentifyPayload {
	var presence *fluxer.UpdatePresencePayload
	if payload.Presence != nil {
		presence = misc.New(PresenceUpdatePayloadToFluxer(*payload.Presence))
	}

	return fluxer.IdentifyPayload{
		Token:          payload.Token,
		Properties:     payload.Properties,
		Compress:       payload.Compress,
		LargeThreshold: payload.LargeThreshold,
		Presence:       presence,
	}
}

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

	guild.Stickers = make([]discord.Sticker, 0, len(event.Stickers))
	for _, sticker := range event.Stickers {
		guild.Stickers = append(guild.Stickers, StickerToDiscord(sticker))
	}

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

func GuildMembersChunkEventToDiscord(event fluxer.GuildMembersChunkEvent) discord.GuildMembersChunkEvent {
	members := make([]discord.GuildMember, 0, len(event.Members))
	for _, member := range event.Members {
		members = append(members, GuildMemberToDiscord(member))
	}

	return discord.GuildMembersChunkEvent{
		GuildID:    event.GuildID,
		Members:    members,
		ChunkIndex: event.ChunkIndex,
		ChunkCount: event.ChunkCount,
		NotFound:   event.NotFound,
		Presences:  event.Presences,
		Nonce:      event.Nonce,
	}
}

func GuildStickersUpdateEventToDiscord(event fluxer.GuildStickersUpdateEvent) discord.GuildStickersUpdateEvent {
	stickers := make([]discord.Sticker, 0, len(event.Stickers))
	for _, sticker := range event.Stickers {
		stickers = append(stickers, StickerToDiscord(sticker))
	}

	return discord.GuildStickersUpdateEvent{
		GuildID:  event.GuildID,
		Stickers: stickers,
	}
}

func MessageCreateEventToDiscord(event fluxer.MessageCreateEvent) discord.MessageCreateEvent {
	var member *discord.GuildMember
	if event.Member != nil {
		member = misc.New(GuildMemberToDiscord(*event.Member))
	}

	return discord.MessageCreateEvent{
		Message: MessageToDiscord(event.Message),
		GuildID: event.GuildID,
		Member:  member,
	}
}
