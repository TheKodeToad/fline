package fluxer

import "github.com/disgoorg/snowflake/v2"

type VoiceState struct {
	GuildID                 *snowflake.ID `json:"guild_id,omitempty"`
	ChannelID               *snowflake.ID `json:"channel_id"`
	UserID                  snowflake.ID  `json:"user_id"`
	Member                  GuildMember   `json:"member"`
	SessionID               string        `json:"session_id"`
	Deaf                    bool          `json:"deaf"`
	Mute                    bool          `json:"mute"`
	SelfDeaf                bool          `json:"self_deaf"`
	SelfMute                bool          `json:"self_mute"`
	SelfVideo               bool          `json:"self_video"`
	Supress                 bool          `json:"supress"`
	RequestToSpeakTimestamp *string       `json:"request_to_speak_timestamp"`
}
