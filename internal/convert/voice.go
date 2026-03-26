package convert

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
)

func VoiceStateToDiscord(state fluxer.VoiceState) discord.VoiceState {
	return discord.VoiceState{
		GuildID:                 state.GuildID,
		ChannelID:               state.ChannelID,
		UserID:                  state.UserID,
		Member:                  GuildMemberToDiscord(state.Member),
		SessionID:               state.SessionID,
		Deaf:                    state.Deaf,
		Mute:                    state.Mute,
		SelfDeaf:                state.SelfDeaf,
		SelfMute:                state.SelfMute,
		SelfVideo:               state.SelfVideo,
		RequestToSpeakTimestamp: state.RequestToSpeakTimestamp,
	}
}
