package convert

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
)

func commonChannelTypeSubset(kind discord.ChannelType) bool {
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

func WebhookToDiscord(webhook fluxer.Webhook) discord.Webhook {
	return discord.Webhook{
		ID:        webhook.ID,
		Type:      discord.WebhookTypeIncoming,
		GuildID:   &webhook.ChannelID,
		ChannelID: &webhook.ChannelID,
		User:      misc.New(UserPartialToDiscord(webhook.User)),
		Name:      &webhook.Name,
		Avatar:    webhook.Avatar,
		Token:     &webhook.Token,
	}
}
