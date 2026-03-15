package flags

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// Set represents a set of flags which may exceed 64-bits.
// It exists mostly to mitigate the poor ergonomics of big.Int, and native ints should be prefered when possible.
// Unlike *big.Int, the zero value is valid and behaves like 0.
type Set struct {
	val *big.Int
}

// New creates a Set as a union of several Sets.
func New(sets ...Set) Set {
	var result Set

	for _, set := range sets {
		result = result.Union(set)
	}

	return result
}

// Bit creates a set only containing 1 << bit.
func Bit(bit uint) Set {
	result := big.NewInt(1)
	result.Lsh(result, bit)

	return Set{result}
}

func (s Set) Intersection(s2 Set) Set {
	if s.val == nil || s2.val == nil {
		return Set{}
	}

	result := big.NewInt(0)
	result.And(s.val, s2.val)

	return Set{result}
}

func (s Set) Union(s2 Set) Set {
	if s.val == nil {
		return s2
	}
	if s2.val == nil {
		return s
	}

	result := big.NewInt(0)
	result.Or(s.val, s2.val)

	return Set{result}
}

// NOTE: unused, this is just to avoid the annoyance if it's needed later
func (s Set) Equal(s2 Set) bool {
	if s.val == nil {
		return s2.val == nil || s2.val.BitLen() == 0
	}
	if s2.val == nil {
		return s.val == nil || s.val.BitLen() == 0
	}

	return s.val.Cmp(s2.val) == 0
}

func (s Set) MarshalJSON() ([]byte, error) {
	text, err := s.val.MarshalText()
	if err != nil {
		return nil, fmt.Errorf("val.MarshalText failed: %w", err)
	}

	return json.Marshal(string(text))
}

func (s *Set) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	if s.val == nil {
		s.val = big.NewInt(0)
	}
	err = s.val.UnmarshalText([]byte(str))
	if err != nil {
		return fmt.Errorf("val.UnmarshalText failed: %w", err)
	}

	return nil
}
