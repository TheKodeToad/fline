package discord_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/misc"
	"github.com/disgoorg/snowflake/v2"
)

func TestGuildMemberUpdateUnmarshal(t *testing.T) {
	expect := func(str string, expected discord.GuildMemberUpdate) {
		var got discord.GuildMemberUpdate
		err := json.Unmarshal([]byte(str), &got)
		if err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if !reflect.DeepEqual(expected, got) {
			t.Errorf("unmarshalling %s: expected %+v but got %+v", str, expected, got)
		}
	}

	expect("{}", discord.GuildMemberUpdate{})
	expect(
		`{"channel_id": null}`,
		discord.GuildMemberUpdate{ClearChannel: true},
	)
	expect(
		`{ "nick": null, "roles": null, "mute": null, "deaf": null, "communication_disabled_until": null }`,
		discord.GuildMemberUpdate{
			Nick:                       misc.New(""),
			Roles:                      []snowflake.ID{},
			Mute:                       misc.New(false),
			Deaf:                       misc.New(false),
			CommunicationDisabledUntil: misc.New(""),
		},
	)
	expect(
		`{ "nick": "joe", "roles": ["1234"], "mute": true, "deaf": true, "channel_id": "5678", "communication_disabled_until": "the great refactor" }`,
		discord.GuildMemberUpdate{
			Nick: misc.New("joe"),
			Roles: []snowflake.ID{1234},
			Mute: misc.New(true),
			Deaf: misc.New(true),
			ChannelID: misc.New(snowflake.ID(5678)),
			CommunicationDisabledUntil: misc.New("the great refactor"),
		},
	)
}
