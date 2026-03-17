package discord

import "github.com/disgoorg/snowflake/v2"

type ChannelMention struct {
	ID      snowflake.ID `json:"id"`
	GuildID snowflake.ID `json:"guild_id"`
	Type    ChannelType  `json:"type"`
	Name    string       `json:"name"`
}

type Attachment struct {
	ID          snowflake.ID `json:"id"`
	Filename    string       `json:"filename"`
	Title       *string      `json:"title,omitempty"`
	Description *string      `json:"description,omitempty"`
	ContentType *string      `json:"content_type,omitempty"`
	Size        int          `json:"size"`
	URL         string       `json:"url"`
	ProxyURL    string       `json:"proxy_url"`
	Width       *int         `json:"width,omitempty"`
	Height      *int         `json:"height,omitempty"`
	Ephemeral   *bool        `json:"ephemeral,omitempty"`
}

type EmbedFooter struct {
	Text         *string  `json:"text,omitempty"`
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
	Inline *bool  `json:"inline"`
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
	Content          *string           `json:"content"`
	Nonce            *Nonce            `json:"nonce"`
	TTS              *bool             `json:"tts"`
	Embeds           []Embed           `json:"embeds"`
	AllowedMentions  *AllowedMentions  `json:"allowed_mentions"`
	MessageReference *MessageReference `json:"message_reference"`
	Flags            int               `json:"flags"`
	EnforceNonce     *bool             `json:"enforce_nonce"`
}
