package fluxer

import (
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/multipartx"
	"github.com/disgoorg/snowflake/v2"
)

type Reaction struct {
	Emoji Emoji `json:"emoji"`
	Count int   `json:"count"`
	Me    bool  `json:"me"`
}

type Message struct {
	ID               snowflake.ID              `json:"id"`
	ChannelID        snowflake.ID              `json:"channel_id"`
	Author           UserPartial               `json:"author"`
	Content          string                    `json:"content"`
	Timestamp        string                    `json:"timestamp"`
	EditedTimestamp  *string                   `json:"edited_timestamp"`
	TTS              bool                      `json:"tts"`
	MentionEveryone  bool                      `json:"mention_everyone"`
	Mentions         []UserPartial             `json:"mentions"`
	MentionRoles     []snowflake.ID            `json:"mention_roles"`
	MentionChannels  []discord.ChannelMention  `json:"mention_channels"`
	Attachments      []discord.Attachment      `json:"attachments"`
	Embeds           []discord.Embed           `json:"embeds"`
	Reactions        []Reaction                `json:"reactions"`
	Nonce            *string                   `json:"nonce,omitempty"`
	Pinned           bool                      `json:"pinned"`
	WebhookID        *snowflake.ID             `json:"webhook_id"`
	Type             discord.MessageType       `json:"type"`
	Flags            discord.MessageFlags      `json:"flags"`
	MessageReference *discord.MessageReference `json:"message_reference"`
	// TODO: referenced message, resolved data?
}

type MessageCreate struct {
	Content          *string                   `json:"content,omitempty"`
	Nonce            *string                   `json:"nonce,omitempty"`
	TTS              *bool                     `json:"tts,omitempty"`
	Embeds           []discord.Embed           `json:"embeds,omitempty"`
	AllowedMentions  *discord.AllowedMentions  `json:"allowed_mentions,omitempty"`
	MessageReference *discord.MessageReference `json:"message_reference,omitempty"`
	StickerIDs       []snowflake.ID            `json:"sticker_ids,omitzero"`
	Files            []multipartx.InMemoryFile `json:"-"`
	Attachments      []discord.Attachment      `json:"attachments,omitzero"`
	Flags            discord.MessageFlags      `json:"flags"`
	EnforceNonce     *bool                     `json:"enforce_nonce,omitempty"`
}

func (mc MessageCreate) EncodeForm(form *multipart.Writer) error {
	payloadJSON, err := json.Marshal(mc)
	if err != nil {
		return fmt.Errorf("failed to marshal payload_json: %w", err)
	}

	form.WriteField("payload_json", string(payloadJSON))

	for _, file := range mc.Files {
		writer, err := form.CreateFormFile(file.FieldName, file.FileName)
		if err != nil {
			return fmt.Errorf("failed to add form file: %w", err)
		}

		_, err = writer.Write(file.Data)
		if err != nil {
			return fmt.Errorf("failed to write message file to form: %w", err)
		}
	}

	return nil
}

type MessageBulkDelete struct {
	MessageIDs []snowflake.ID `json:"message_ids"`
}
