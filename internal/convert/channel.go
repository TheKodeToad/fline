package convert

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
)

func IsChannelConvertible(kind discord.ChannelType) bool {
	switch kind {
	case discord.ChannelTypeGuildText,
		discord.ChannelTypeDM,
		discord.ChannelTypeGuildVoice,
		discord.ChannelTypeGroupDM,
		discord.ChannelTypeGuildCategory:
		return true
	default:
		return false
	}
}

func ChannelToDiscord(channel fluxer.Channel) (discord.Channel, bool) {
	if !IsChannelConvertible(channel.Type) {
		return discord.Channel{}, false
	}

	var recipients []discord.User
	if channel.Recipients != nil {
		recipients = make([]discord.User, 0, len(channel.Recipients))
		for _, recipient := range channel.Recipients {
			recipients = append(recipients, UserPartialToDiscord(recipient))
		}
	}

	return discord.Channel{
		ID:                   channel.ID,
		Type:                 channel.Type,
		GuildID:              channel.GuildID,
		Position:             channel.Position,
		PermissionOverwrites: channel.PermissionOverwrites,
		Name:                 channel.Name,
		Topic:                channel.Topic,
		NSFW:                 channel.NSFW,
		LastMessageID:        channel.LastMessageID,
		Bitrate:              channel.Bitrate,
		UserLimit:            channel.UserLimit,
		RateLimitPerUser:     channel.RateLimitPerUser,
		Recipients:           recipients,
	}, true
}

func ChannelCreateToFluxer(create discord.ChannelCreate) fluxer.ChannelCreate {
	// NOTE: yep, Fluxer just makes this required
	var kind discord.ChannelType
	if create.Type != nil {
		kind = *create.Type
	}

	return fluxer.ChannelCreate{
		Name:                 create.Name,
		Type:                 kind,
		Topic:                create.Topic,
		Bitrate:              create.Bitrate,
		UserLimit:            create.UserLimit,
		Position:             create.Position,
		PermissionOverwrites: create.PermissionOverwrites,
		ParentID:             create.ParentID,
		NSFW:                 create.NSFW,
		RTCRegion:            create.RTCRegion,
	}
}
