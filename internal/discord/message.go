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

type ChannelMention struct {
	ID      snowflake.ID `json:"id"`
	GuildID snowflake.ID `json:"guild_id"`
	Type    ChannelType  `json:"type"`
	Name    string       `json:"name"`
}

type Attachment struct {
	ID          snowflake.ID `json:"id"`
	Filename    *string      `json:"filename,omitempty"`
	Title       *string      `json:"title,omitempty"`
	Description *string      `json:"description,omitempty"`
	ContentType *string      `json:"content_type,omitempty"`
	Size        *int         `json:"size,omitempty"`
	URL         *string      `json:"url,omitempty"`
	ProxyURL    *string      `json:"proxy_url,omitempty"`
	Width       *int         `json:"width,omitempty"`
	Height      *int         `json:"height,omitempty"`
}

type EmbedFooter struct {
	Text         *string `json:"text,omitempty"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

type EmbedMedia struct {
	URL      *string `json:"url,omitempty"`
	ProxyURL *string `json:"proxy_url,omitempty"`
	Width    *int    `json:"width,omitempty"`
	Height   *int    `json:"height,omitempty"`
}

type EmbedProvider struct {
	Name *string `json:"name,omitempty"`
	URL  *string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         *string `json:"name,omitempty"`
	URL          *string `json:"url,omitempty"`
	IconURL      *string `json:"icon_url,omitempty"`
	ProxyIconURL *string `json:"proxy_icon_url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline *bool  `json:"inline,omitempty"`
}

type Embed struct {
	Title       *string        `json:"title,omitempty"`
	Type        *string        `json:"type,omitempty"`
	Description *string        `json:"description,omitempty"`
	URL         *string        `json:"url,omitempty"`
	Timestamp   *string        `json:"timestamp"`
	Color       *uint          `json:"color,omitempty"`
	Footer      *EmbedFooter   `json:"footer,omitempty"`
	Image       *EmbedMedia    `json:"image,omitempty"`
	Thumbnail   *EmbedMedia    `json:"thumbnail,omitempty"`
	Provider    *EmbedProvider `json:"provider,omitempty"`
	Author      *EmbedAuthor   `json:"author,omitempty"`
	Fields      []EmbedField   `json:"fields,omitzero"`
}

type ReactionCountDetails struct {
	Burst  int `json:"burst"`
	Normal int `json:"normal"`
}

type ReactionType uint

const (
	ReactionTypeNormal ReactionType = 0
)

type Reaction struct {
	Count        int                  `json:"count"`
	CountDetails ReactionCountDetails `json:"count_details"`
	Me           bool                 `json:"me"`
	MeBurst      bool                 `json:"me_burst"`
	Emoji        Emoji                `json:"emoji"`
	BurstColors  []uint               `json:"burst_colors"`
}

type MessageType uint

type MessageFlags uint

var (
	MessageFlagSupressEmbeds        MessageFlags = 1 << 2
	MessageFlagSupressNotifications MessageFlags = 1 << 12
	MessageFlagVoiceMessage         MessageFlags = 1 << 13
)

type MessageReference struct {
	MessageID       *snowflake.ID `json:"message_id,omitempty"`
	ChannelID       *snowflake.ID `json:"channel_id,omitempty"`
	GuildID         *snowflake.ID `json:"guild_id,omitempty"`
	FailIfNotExists *bool         `json:"fail_if_not_exists,omitempty"`
}

type Message struct {
	ID               snowflake.ID      `json:"id"`
	ChannelID        snowflake.ID      `json:"channel_id"`
	Author           User              `json:"author"`
	Content          string            `json:"content"`
	Timestamp        string            `json:"timestamp"`
	EditedTimestamp  *string           `json:"edited_timestamp"`
	TTS              bool              `json:"tts"`
	MentionEveryone  bool              `json:"mention_everyone"`
	Mentions         []User            `json:"mentions"`
	MentionRoles     []snowflake.ID    `json:"mention_roles"`
	MentionChannels  []ChannelMention  `json:"mention_channels,omitzero"`
	Attachments      []Attachment      `json:"attachments"`
	Embeds           []Embed           `json:"embeds"`
	Reactions        []Reaction        `json:"reactions"`
	Nonce            *Nonce            `json:"nonce,omitempty"`
	Pinned           bool              `json:"pinned"`
	WebhookID        *snowflake.ID     `json:"webhook_id"`
	Type             MessageType       `json:"type"`
	Flags            MessageFlags      `json:"flags"`
	MessageReference *MessageReference `json:"message_reference,omitempty"`
	// TODO: referenced message, resolved data?
}

type AllowedMentions struct {
	Parse       []string       `json:"parse,omitzero"`
	Roles       []snowflake.ID `json:"roles,omitzero"`
	Users       []snowflake.ID `json:"users,omitzero"`
	RepliedUser *bool          `json:"replied_user,omitempty"`
}

type MessageCreate struct {
	Content          *string                   `json:"content,omitempty"`
	Nonce            *Nonce                    `json:"nonce,omitempty"`
	TTS              *bool                     `json:"tts,omitempty"`
	Embeds           []Embed                   `json:"embeds,omitzero"`
	AllowedMentions  *AllowedMentions          `json:"allowed_mentions,omitempty"`
	MessageReference *MessageReference         `json:"message_reference,omitempty"`
	StickerIDs       []snowflake.ID            `json:"sticker_ids,omitzero"`
	Files            []multipartx.InMemoryFile `json:"-"`
	Attachments      []Attachment              `json:"attachments,omitzero"`
	Flags            MessageFlags              `json:"flags,omitempty"`
	EnforceNonce     *bool                     `json:"enforce_nonce,omitempty"`
}

func (mc *MessageCreate) UnmarshalForm(form multipartx.InMemoryForm) error {
	if payloadJSON := form.Value["payload_json"]; len(payloadJSON) != 0 {
		err := json.Unmarshal([]byte(payloadJSON[0]), mc)
		if err != nil {
			return fmt.Errorf("failed to unmarshal payload_json: %w", err)
		}
	} else {
		if content := form.Value["content"]; len(content) != 0 {
			mc.Content = misc.New(content[0])
		}

		if nonce := form.Value["nonce"]; len(nonce) != 0 {
			mc.Nonce = misc.New(NonceFromString(nonce[0]))
		}

		if tts := form.Value["tts"]; len(tts) != 0 {
			switch tts[0] {
			case "false":
				mc.TTS = misc.New(false)
			case "true":
				mc.TTS = misc.New(true)
			default:
				return errors.New("expected true/false for tts form value")
			}
		}

		if stickerIDs := form.Value["sticker_ids"]; len(stickerIDs) != 0 {
			mc.StickerIDs = make([]snowflake.ID, 0, len(stickerIDs))
			for _, stickerID := range stickerIDs {
				id, err := strconv.ParseUint(stickerID, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse sticker_ids form value: %w", err)
				}

				mc.StickerIDs = append(mc.StickerIDs, snowflake.ID(id))
			}
		}

		if flags := form.Value["flags"]; len(flags) != 0 {
			flagsInt, err := strconv.Atoi(flags[0])
			if err != nil {
				return fmt.Errorf("failed to parse flags form value: %w", err)
			}

			mc.Flags = MessageFlags(flagsInt)
		}

		if enforceNonce := form.Value["enforce_nonce"]; len(enforceNonce) != 0 {
			switch enforceNonce[0] {
			case "false":
				mc.EnforceNonce = misc.New(false)
			case "true":
				mc.EnforceNonce = misc.New(true)
			default:
				return errors.New("expected true/false for enforce_nonce form value")
			}
		}
	}

	mc.Files = form.Files
	return nil
}

type MessageEdit struct {
	Content         *string                   `json:"content,omitempty"`
	Embeds          []Embed                   `json:"embeds,omitzero"`
	Flags           *MessageFlags             `json:"flags,omitempty"`
	AllowedMentions *AllowedMentions          `json:"allowed_mentions,omitempty"`
	Files           []multipartx.InMemoryFile `json:"-"`
	Attachments     []Attachment              `json:"attachments,omitzero"`
}

func (me *MessageEdit) UnmarshalJSON(data []byte) error {
	var raw struct {
		Content         json.RawMessage `json:"content"`
		Embeds          json.RawMessage `json:"embeds"`
		Flags           json.RawMessage `json:"flags"`
		AllowedMentions json.RawMessage `json:"allowed_mentions"`
		Attachments     json.RawMessage `json:"attachments"`
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	// NOTE: we unfortunately need this manual parsing to distinguish between null and undefined

	if string(raw.Content) == "null" {
		me.Content = misc.New("")
	} else if raw.Content != nil {
		err := json.Unmarshal(raw.Content, &me.Content)
		if err != nil {
			return fmt.Errorf("unmarshalling into MessageEdit.Content: %w", err)
		}
	}

	if string(raw.Embeds) == "null" {
		me.Embeds = []Embed{}
	} else if raw.Embeds != nil {
		err := json.Unmarshal(raw.Embeds, &me.Embeds)
		if err != nil {
			return fmt.Errorf("unmarshalling into MessageEdit.Embeds: %w", err)
		}
	}

	if string(raw.Flags) == "null" {
		me.Flags = misc.New(MessageFlags(0))
	} else if raw.Flags != nil {
		err := json.Unmarshal(raw.Embeds, &me.Embeds)
		if err != nil {
			return fmt.Errorf("unmarshalling into MessageEdit.Flags: %w", err)
		}
	}

	if raw.AllowedMentions != nil {
		err = json.Unmarshal(raw.AllowedMentions, &me.AllowedMentions)
		if err != nil {
			return fmt.Errorf("unmarshalling into MessageEdit.AllowedMentions: %w", err)
		}
	}

	if string(raw.Attachments) == "null" {
		me.Attachments = []Attachment{}
	} else if raw.Attachments != nil {
		err = json.Unmarshal(raw.Attachments, &me.Attachments)
		if err != nil {
			return fmt.Errorf("unmarshalling into MessageEdit.Attachments: %w", err)
		}
	}

	return nil
}

func (me *MessageEdit) UnmarshalForm(form multipartx.InMemoryForm) error {
	if payloadJSON := form.Value["payload_json"]; len(payloadJSON) != 0 {
		err := json.Unmarshal([]byte(payloadJSON[0]), me)
		if err != nil {
			return fmt.Errorf("failed to unmarshal payload_json: %w", err)
		}
	} else {
		if content := form.Value["content"]; len(content) != 0 {
			me.Content = misc.New(content[0])
		}

		if flags := form.Value["flags"]; len(flags) != 0 {
			flagsInt, err := strconv.Atoi(flags[0])
			if err != nil {
				return fmt.Errorf("failed to parse flags form value: %w", err)
			}

			me.Flags = misc.New(MessageFlags(flagsInt))
		}
	}

	me.Files = form.Files
	return nil
}

func (me *MessageEdit) EncodeForm(form *multipart.Writer) error {
	payloadJSON, err := json.Marshal(me)
	if err != nil {
		return fmt.Errorf("failed to marshal payload_json: %w", err)
	}

	form.WriteField("payload_json", string(payloadJSON))

	for _, file := range me.Files {
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
	Messages []snowflake.ID `json:"messages"`
}
