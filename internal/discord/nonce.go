// Package Discord contains models which represent data from Discord's API.
package discord

import (
	"encoding/json"
	"strconv"
)

// Nonce represents a Discord nonce value which is either a string or integer.
type Nonce struct {
	isStr    bool
	strValue string
	intValue int64
}

func NonceFromString(s string) Nonce {
	return Nonce{isStr: true, strValue: s}
}

func NonceFromInt(i int64) Nonce {
	return Nonce{isStr: false, intValue: i}
}

// String coerces the nonce value to a string.
// This is for the sake of implementing Stringer - you should use StringValue instead.
func (n Nonce) String() string {
	if n.isStr {
		return n.strValue
	} else {
		return strconv.FormatInt(n.intValue, 10)
	}
}

func (n Nonce) StringValue() (string, bool) {
	return n.strValue, n.isStr
}

func (n Nonce) IntValue() (int64, bool) {
	return n.intValue, !n.isStr
}

func (n Nonce) MarshalJSON() ([]byte, error) {
	if n.isStr {
		return json.Marshal(n.strValue)
	} else {
		return json.Marshal(n.intValue)
	}
}

func (n Nonce) UnmarshalJSON(data []byte) error {
	// FIXME: probably handles invalid input badly
	if data[0] == '"' {
		return json.Unmarshal(data, &n.strValue)
	} else {
		return json.Unmarshal(data, &n.intValue)
	}
}
