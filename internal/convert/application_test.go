package convert_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/TheKodeToad/fine/internal/convert"
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/TheKodeToad/fine/internal/fluxer"
	"github.com/TheKodeToad/fine/internal/misc"
)

func TestApplicationToDiscord(t *testing.T) {
	input := fluxer.Application{
		ID:                  1234,
		Name:                "Test",
		Icon:                misc.New("appicon"),
		Description:         misc.New("App Description"),
		BotPublic:           true,
		BotRequireCodeGrant: false,
		Bot: &fluxer.Bot{
			ID:            5678,
			Username:      "TestBot",
			Discriminator: "0000",
			Avatar:        misc.New("botavatar"),
			Banner:        misc.New("botbanner"),
			Bio:           misc.New("This is a cool bio"),
			MFAEnabled:    false,
		},
	}
	expected := discord.Application{
		ID:                  1234,
		Name:                "Test",
		Icon:                misc.New("appicon"),
		Description:         "This is a cool bio",
		BotPublic:           true,
		BotRequireCodeGrant: false,
		Bot: &discord.User{
			ID:            5678,
			Username:      "TestBot",
			Discriminator: "0000",
			Avatar:        misc.New("botavatar"),
			Bot:           misc.New(true),
			Banner:        misc.New("botbanner"),
			Flags:         misc.New(discord.UserFlags(0)),
		},
		Flags: discord.ApplicationGatewayPresence |
			discord.ApplicationGatewayGuildMembers |
			discord.ApplicationGatewayMessageContent,
	}
	got := convert.ApplicationToDiscord(input)

	if !reflect.DeepEqual(expected, got) {
		gotJSON, err := json.MarshalIndent(got, "", "  ")
		if err != nil {
			t.Fatalf("json.Marshal(got) failed: %v", err)
		}

		expectedJSON, err := json.MarshalIndent(expected, "", "  ")
		if err != nil {
			t.Fatalf("json.Marshal(expected) failed: %v", err)

		}

		t.Errorf("\nexpected: %s\ngot: %s", string(expectedJSON), string(gotJSON))
	}

}
