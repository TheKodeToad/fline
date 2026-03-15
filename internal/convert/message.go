package convert

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
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

var messageFlagsCommonSubset = discord.MessageFlagSupressEmbeds |
	discord.MessageFlagSupressNotifications |
	discord.MessageFlagVoiceMessage

func MessageFlagsToDiscord(flags discord.MessageFlags) discord.MessageFlags {
	return flags & messageFlagsCommonSubset
}

func MessageFlagsToFluxer(flags discord.MessageFlags) discord.MessageFlags {
	return flags & messageFlagsCommonSubset
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
		nonce = misc.New(NonceToDiscord(*message.Nonce))
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

func MessageCreateToFluxer(create discord.MessageCreate) fluxer.MessageCreate {
	var nonce *string
	if create.Nonce != nil {
		nonce = misc.New(NonceToFluxer(*create.Nonce))
	}

	allowedMentions := create.AllowedMentions
	if allowedMentions != nil {
		// NOTE: Discord defaults differ from Fluxer; apply Discord defaults here
		if allowedMentions.Parse == nil {
			allowedMentions.Parse = []string{}
		}
		if allowedMentions.Roles == nil {
			allowedMentions.Roles = []snowflake.ID{}
		}
		if allowedMentions.Users == nil {
			allowedMentions.Users = []snowflake.ID{}
		}
		if allowedMentions.RepliedUser == nil {
			allowedMentions.RepliedUser = misc.New(true)
		}
	}

	return fluxer.MessageCreate{
		Content:          create.Content,
		Nonce:            nonce,
		TTS:              create.TTS,
		Embeds:           create.Embeds,
		AllowedMentions:  allowedMentions,
		MessageReference: create.MessageReference,
		Flags:            create.Flags,
		EnforceNonce:     create.EnforceNonce,
	}
}
