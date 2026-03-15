package convert

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
)

func commonChannelTypeSubset(kind discord.ChannelType) bool {
	switch kind {
	case discord.ChannelTypeGuildText:
	case discord.ChannelTypeDM:
	case discord.ChannelTypeGuildVoice:
	case discord.ChannelTypeGroupDM:
	case discord.ChannelTypeGuildCategory:
		return true
	}
	return false
}

func ChannelToDiscord(channel fluxer.Channel) (discord.Channel, bool) {
	if !commonChannelTypeSubset(channel.Type) {
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
