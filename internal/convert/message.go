package convert

import (
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/TheKodeToad/fine/internal/fluxer"
	"github.com/TheKodeToad/fine/internal/misc"
	"github.com/disgoorg/snowflake/v2"
)

func ReactionToDiscord(reaction fluxer.Reaction) discord.Reaction {
	return discord.Reaction{
		Count: reaction.Count,
		CountDetails: discord.ReactionCountDetails{
			Burst:  0,
			Normal: reaction.Count,
		},
		Me:          reaction.Me,
		MeBurst:     false,
		Emoji:       reaction.Emoji,
		BurstColors: []uint{},
	}
}

func MessageFlagsToDiscord(flags discord.MessageFlags) discord.MessageFlags {
	commonSubset := discord.MessageFlagSupressEmbeds |
		discord.MessageFlagSupressNotifications |
		discord.MessageFlagVoiceMessage

	return flags & commonSubset
}

func MessageToDiscord(message fluxer.Message) discord.Message {
	mentions := make([]discord.User, 0, len(message.Mentions))
	for _, mention := range message.Mentions {
		mentions = append(mentions, UserPartialToDiscord(mention))
	}

	mentionRoles := message.MentionRoles
	if mentionRoles == nil {
		mentionRoles = []snowflake.ID{}
	}

	reactions := make([]discord.Reaction, 0, len(message.Reactions))
	for _, reaction := range message.Reactions {
		reactions = append(reactions, ReactionToDiscord(reaction))
	}

	var nonce *discord.Nonce
	if message.Nonce != nil {
		nonce = misc.New(discord.NonceFromString(*message.Nonce))
	}

	return discord.Message{
		ID:               message.ID,
		ChannelID:        message.ChannelID,
		Author:           UserPartialToDiscord(message.Author),
		Content:          message.Content,
		Timestamp:        message.Timestamp,
		EditedTimestamp:  message.EditedTimestamp,
		TTS:              message.TTS,
		MentionEveryone:  message.MentionEveryone,
		Mentions:         mentions,
		MentionRoles:     mentionRoles,
		MentionChannels:  message.MentionChannels,
		Attachments:      message.Attachments,
		Embeds:           message.Embeds,
		Reactions:        reactions,
		Nonce:            nonce,
		Pinned:           message.Pinned,
		Type:             message.Type,
		Flags:            MessageFlagsToDiscord(message.Flags),
		MessageReference: message.MessageReference,
	}
}
