package convert

import (
	"fmt"

	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
	"github.com/TheKodeToad/fline/internal/multipartx"
	"github.com/disgoorg/snowflake/v2"
)

func EmbedFieldToFluxer(field discord.EmbedField) discord.EmbedField {
	if field.Value == "" {
		// NOTE: Fluxer doesn't allow any empty string here unlike Discord
		field.Value = "-"
	}

	return field
}

func EmbedToFluxer(embed discord.Embed) discord.Embed {
	for i, field := range embed.Fields {
		embed.Fields[i] = EmbedFieldToFluxer(field)
	}

	// NOTE: even though Discord documents these as not being optional they are
	if embed.Footer != nil && embed.Footer.Text == nil {
		embed.Footer = nil
	}
	if embed.Image != nil && embed.Image.URL == nil {
		embed.Image = nil
	}
	if embed.Thumbnail != nil && embed.Thumbnail.URL == nil {
		embed.Thumbnail = nil
	}
	if embed.Author != nil && embed.Author.Name == nil {
		embed.Author = nil
	}

	return embed
}

func ReactionToDiscord(reaction fluxer.Reaction) discord.Reaction {
	return discord.Reaction{
		Count: reaction.Count,
		CountDetails: discord.ReactionCountDetails{
			Burst:  0,
			Normal: reaction.Count,
		},
		Me:          reaction.Me,
		MeBurst:     false,
		Emoji:       EmojiToDiscord(reaction.Emoji),
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

func UploadsToFluxer(inAttachments []discord.Attachment, inFiles []multipartx.InMemoryFile) ([]discord.Attachment, []multipartx.InMemoryFile, bool) {
	pendingFiles := make([]multipartx.InMemoryFile, 0, len(inFiles))
	pendingFiles = append(pendingFiles, inFiles...)

	outAttachments := make([]discord.Attachment, 0, len(inAttachments))
	outFiles := make([]multipartx.InMemoryFile, 0, len(inFiles))

	var newID int
	for _, attachment := range inAttachments {
		oldID := attachment.ID
		attachment.ID = snowflake.ID(newID)

		oldFieldName := fmt.Sprintf("files[%d]", oldID)
		for i, file := range pendingFiles {
			if file.FieldName == oldFieldName {
				pendingFiles[i] = pendingFiles[len(pendingFiles)-1]
				pendingFiles = pendingFiles[:len(pendingFiles)-1]

				if attachment.Filename == nil {
					attachment.Filename = &file.FileName
				}

				outFiles = append(outFiles, multipartx.InMemoryFile{
					FieldName: fmt.Sprintf("files[%d]", newID),
					FileName:  file.FileName,
					Data:      file.Data,
				})
				break
			}
		}

		outAttachments = append(outAttachments, attachment)
		newID++
	}

	for _, file := range pendingFiles {
		outFiles = append(outFiles, multipartx.InMemoryFile{
			FieldName: fmt.Sprintf("files[%d]", newID),
			FileName:  file.FileName,
			Data:      file.Data,
		})

		outAttachments = append(outAttachments, discord.Attachment{
			ID:       snowflake.ID(newID),
			Filename: &file.FileName,
		})

		newID++
	}

	return outAttachments, outFiles, true
}

func AllowedMentionsToFluxer(mentions discord.AllowedMentions) discord.AllowedMentions {
	if mentions.Parse == nil {
		mentions.Parse = []string{}
	}
	if mentions.Roles == nil {
		mentions.Roles = []snowflake.ID{}
	}
	if mentions.Users == nil {
		mentions.Users = []snowflake.ID{}
	}
	if mentions.RepliedUser == nil {
		mentions.RepliedUser = misc.New(true)
	}

	return mentions
}

func MessageCreateToFluxer(create discord.MessageCreate) (fluxer.MessageCreate, bool) {
	var nonce *string
	if create.Nonce != nil {
		nonce = misc.New(NonceToFluxer(*create.Nonce))
	}

	var allowedMentions *discord.AllowedMentions
	if create.AllowedMentions != nil {
		allowedMentions = misc.New(AllowedMentionsToFluxer(*create.AllowedMentions))
	}

	embeds := make([]discord.Embed, 0, len(create.Embeds))
	for _, embed := range create.Embeds {
		embeds = append(embeds, EmbedToFluxer(embed))
	}

	attachments, files, ok := UploadsToFluxer(create.Attachments, create.Files)
	if !ok {
		return fluxer.MessageCreate{}, false
	}

	return fluxer.MessageCreate{
		Content:          create.Content,
		Nonce:            nonce,
		TTS:              create.TTS,
		Embeds:           embeds,
		AllowedMentions:  allowedMentions,
		MessageReference: create.MessageReference,
		StickerIDs:       create.StickerIDs,
		Files:            files,
		Attachments:      attachments,
		Flags:            MessageFlagsToFluxer(create.Flags),
		EnforceNonce:     create.EnforceNonce,
	}, true
}

func MessageEditToFluxer(edit discord.MessageEdit) (discord.MessageEdit, bool) {
	var embeds []discord.Embed
	if edit.Embeds != nil {
		embeds := make([]discord.Embed, 0, len(edit.Embeds))
		for _, embed := range edit.Embeds {
			embeds = append(embeds, EmbedToFluxer(embed))
		}
	}

	var flags *discord.MessageFlags
	if edit.Flags != nil {
		flags = misc.New(MessageFlagsToDiscord(*edit.Flags))
	}

	var allowedMentions *discord.AllowedMentions
	if edit.AllowedMentions != nil {
		allowedMentions = misc.New(AllowedMentionsToFluxer(*edit.AllowedMentions))
	}

	var attachments []discord.Attachment
	var files []multipartx.InMemoryFile

	if len(edit.Attachments) != 0 || len(edit.Files) != 0 {
		var ok bool
		attachments, files, ok = UploadsToFluxer(edit.Attachments, edit.Files)
		if !ok {
			return discord.MessageEdit{}, false
		}
	} else if edit.Attachments != nil {
		attachments = []discord.Attachment{}
	}

	return discord.MessageEdit{
		Content:         edit.Content,
		Embeds:          embeds,
		Flags:           flags,
		AllowedMentions: allowedMentions,
		Files:           files,
		Attachments:     attachments,
	}, true
}

func MessageBulkDeleteToFluxer(delete discord.MessageBulkDelete) fluxer.MessageBulkDelete {
	return fluxer.MessageBulkDelete{MessageIDs: delete.Messages}
}
