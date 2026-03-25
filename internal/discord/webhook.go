package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"

	"github.com/TheKodeToad/fline/internal/misc"
	"github.com/TheKodeToad/fline/internal/multipartx"
	"github.com/disgoorg/snowflake/v2"
)

type WebhookType uint

const (
	WebhookTypeIncoming WebhookType = 1
)

type WebhookCreate struct {
	Name   string  `json:"name"`
	Avatar *string `json:"avatar"`
}

type Webhook struct {
	ID        snowflake.ID  `json:"id"`
	Type      WebhookType   `json:"type"`
	GuildID   *snowflake.ID `json:"guild_id"`
	ChannelID *snowflake.ID `json:"channel_id"`
	User      *User         `json:"user,omitempty"`
	Name      *string       `json:"name"`
	Avatar    *string       `json:"avatar"`
	Token     *string       `json:"token,omitempty"`
}

type WebhookExecute struct {
	Wait             bool                      `json:"-"`
	Content          *string                   `json:"content,omitempty"`
	Username         *string                   `json:"username,omitempty"`
	AvatarURL        *string                   `json:"avatar_url,omitempty"`
	TTS              *bool                     `json:"tts,omitempty"`
	Embeds           []Embed                   `json:"embeds,omitzero"`
	AllowedMentions  *AllowedMentions          `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference         `json:"message_reference,omitempty"`
	StickerIDs       []snowflake.ID            `json:"sticker_ids,omitzero"`
	Files            []multipartx.InMemoryFile `json:"-"`
	Attachments      []Attachment              `json:"attachments,omitzero"`
	Flags            int                       `json:"flags,omitempty"`
}

func (we WebhookExecute) EncodeForm(form *multipart.Writer) error {
	payloadJSON, err := json.Marshal(we)
	if err != nil {
		return fmt.Errorf("failed to marshal payload_json: %w", err)
	}

	form.WriteField("payload_json", string(payloadJSON))

	for _, file := range we.Files {
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

func (we *WebhookExecute) UnmarshalForm(form multipartx.InMemoryForm) error {
	if payloadJSON := form.Value["payload_json"]; len(payloadJSON) != 0 {
		err := json.Unmarshal([]byte(payloadJSON[0]), we)
		if err != nil {
			return fmt.Errorf("failed to unmarshal payload_json: %w", err)
		}
	} else {
		if content := form.Value["content"]; len(content) != 0 {
			we.Content = misc.New(content[0])
		}

		if username := form.Value["username"]; len(username) != 0 {
			we.Username = misc.New(username[0])
		}

		if avatarURL := form.Value["avatar_url"]; len(avatarURL) != 0 {
			we.AvatarURL = misc.New(avatarURL[0])
		}

		if tts := form.Value["tts"]; len(tts) != 0 {
			switch tts[0] {
			case "false":
				we.TTS = misc.New(false)
			case "true":
				we.TTS = misc.New(true)
			default:
				return errors.New("expected true/false for tts form value")
			}
		}

		if stickerIDs := form.Value["sticker_ids"]; len(stickerIDs) != 0 {
			we.StickerIDs = make([]snowflake.ID, 0, len(stickerIDs))
			for _, stickerID := range stickerIDs {
				id, err := strconv.ParseUint(stickerID, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse sticker_ids form value: %w", err)
				}

				we.StickerIDs = append(we.StickerIDs, snowflake.ID(id))
			}
		}

		if flags := form.Value["flags"]; len(flags) != 0 {
			flagsInt, err := strconv.Atoi(flags[0])
			if err != nil {
				return fmt.Errorf("failed to parse flags form value: %w", err)
			}

			we.Flags = flagsInt
		}
	}

	we.Files = form.Files
	return nil
}
