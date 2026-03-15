package discord_test

import (
	"encoding/json"
	"testing"

	"github.com/TheKodeToad/fline/internal/discord"
)

func TestNonceStringJSON(t *testing.T) {
	const nonceStr = "moonwelotcosidme"
	const jsonStr = `"moonwelotcosidme"`

	gotJSONStr, err := json.Marshal(discord.NonceFromString(nonceStr))
	if err != nil {
		t.Fatalf("failed to marshal nonce: %v", err)
	}

	if string(gotJSONStr) != jsonStr {
		t.Errorf("expected json.Marshal to return %s but got %s", jsonStr, gotJSONStr)
	}

	var gotNonce discord.Nonce
	err = json.Unmarshal([]byte(jsonStr), &gotNonce)
	if err != nil {
		t.Fatalf("failed to unmarshal nonce: %v", err)
	}

	if !gotNonce.IsString() {
		t.Errorf("expected json.Unmarshal to yield a string nonce")
	} else if gotNonce.StringValue() != nonceStr {
		t.Errorf("expected json.Unmarshal to yield a nonce of %s but got %s", nonceStr, gotNonce.StringValue())
	}
}

func TestNonceIntJSON(t *testing.T) {
	const nonceInt = 12345678
	const jsonStr = "12345678"

	gotJSONStr, err := json.Marshal(discord.NonceFromInt(nonceInt))
	if err != nil {
		t.Fatalf("failed to marshal nonce: %v", err)
	}

	if string(gotJSONStr) != jsonStr {
		t.Errorf("expected json.Marshal to return %s but got %s", jsonStr, gotJSONStr)
	}

	var gotNonce discord.Nonce
	err = json.Unmarshal([]byte(jsonStr), &gotNonce)
	if err != nil {
		t.Fatalf("failed to unmarshal nonce: %v", err)
	}

	if !gotNonce.IsInt() {
		t.Errorf("expected json.Unmarshal to yield a int nonce")
	} else if gotNonce.IntValue() != nonceInt {
		t.Errorf("expected json.Unmarshal to yield a nonce of %d but got %v", nonceInt, gotNonce)
	}
}
