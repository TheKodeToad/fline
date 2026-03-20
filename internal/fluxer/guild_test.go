package fluxer_test

import (
	"encoding/json"
	"testing"

	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
	"github.com/disgoorg/snowflake/v2"
)

func TestGuildMemberUpdateMarshal(t *testing.T) {
	expect := func(update fluxer.GuildMemberUpdate, expected string) {
		data, err := json.Marshal(update)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		if string(data) != expected {
			t.Errorf("marshalling %+v: expected '%s' but got '%s'", update, expected, string(data))
		}
	}

	expect(fluxer.GuildMemberUpdate{}, "{}")
	expect(
		fluxer.GuildMemberUpdate{
			ClearNick: true,
		},
		`{"nick":null}`,
	)
	expect(
		fluxer.GuildMemberUpdate{
			ClearChannel: true,
		},
		`{"channel_id":null}`,
	)
	expect(
		fluxer.GuildMemberUpdate{
			ClearTimeout: true,
		},
		`{"communication_disabled_until":null}`,
	)
	expect(
		fluxer.GuildMemberUpdate{
			Nick:                       misc.New("hello"),
			Roles:                      []snowflake.ID{},
			Mute:                       misc.New(true),
			Deaf:                       misc.New(true),
			ChannelID:                  misc.New(snowflake.ID(1234)),
			CommunicationDisabledUntil: misc.New("heat death"),
		},
		`{"nick":"hello","roles":[],"mute":true,"deaf":true,"channel_id":"1234","communication_disabled_until":"heat death"}`,
	)
}
