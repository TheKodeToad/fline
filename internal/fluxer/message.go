package fluxer

import (
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/disgoorg/snowflake/v2"
)

type Reaction struct {
	Emoji discord.Emoji `json:"emoji"`
	Count int           `json:"count"`
	Me    bool          `json:"me"`
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
