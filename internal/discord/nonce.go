package discord

import (
	"encoding/json"
	"strconv"
)

// Nonce represents a Discord nonce value which is either a string or integer.
type Nonce struct {
	// NOTE: these are not exported to avoid invalid states :)
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

func (n Nonce) IsString() bool {
	return n.isStr
}

func (n Nonce) IsInt() bool {
	return !n.isStr
}

// StringValue returns the string value, or an empty string if it is not a nonce string.
// You should check IsString before calling this.
func (n Nonce) StringValue() string {
	return n.strValue
}

// StringValue returns the int value, or 0 if it is not a nonce int.
// You should check IsString before calling this.
func (n Nonce) IntValue() int64 {
	return n.intValue
}

func (n Nonce) MarshalJSON() ([]byte, error) {
	if n.isStr {
		return json.Marshal(n.strValue)
	} else {
		return json.Marshal(n.intValue)
	}
}

func (n *Nonce) UnmarshalJSON(data []byte) error {
	// FIXME: probably handles invalid input badly
	if data[0] == '"' {
		// NOTE: resetting the value to avoid both fields being set
		*n = Nonce{isStr: true}
		return json.Unmarshal(data, &n.strValue)
	} else {
		*n = Nonce{isStr: false}
		return json.Unmarshal(data, &n.intValue)
	}
}
