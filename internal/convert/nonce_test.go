package convert_test

import (
	"testing"

	"github.com/TheKodeToad/fine/internal/convert"
	"github.com/TheKodeToad/fine/internal/discord"
)

func TestNonceInt(t *testing.T) {
	const val = 0xDEADBEEF

	fluxer := convert.NonceToFluxer(discord.NonceFromInt(val))
	discord := convert.NonceToDiscord(fluxer)

	if discord.IntValue() != val {
		t.Errorf("expected nonce value of %d to preserved but got %v", val, discord)
	}
}

func TestNonceString(t *testing.T) {
	const val = "alivebeef"

	fluxer := convert.NonceToFluxer(discord.NonceFromString(val))
	discord := convert.NonceToDiscord(fluxer)

	if discord.StringValue() != val {
		t.Errorf("expected nonce value of %s to preserved but got %v", val, discord)
	}
}
