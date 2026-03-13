package convert

import (
	"strconv"
	"strings"

	"github.com/TheKodeToad/fine/internal/discord"
)

const numericNoncePrefix = "numeric!!"

func NonceToDiscord(nonce string) discord.Nonce {
	if strings.HasPrefix(nonce, numericNoncePrefix) {
		stripped := nonce[len(numericNoncePrefix):]
		parsed, err := strconv.ParseInt(stripped, 10, 64)
		if err != nil {
			return discord.NonceFromString(nonce)
		}

		return discord.NonceFromInt(parsed)
	} else {
		return discord.NonceFromString(nonce)
	}
}

func NonceToFluxer(nonce discord.Nonce) string {
	if str, ok := nonce.StringValue(); ok {
		return str
	} else if i, ok := nonce.IntValue(); ok {
		return numericNoncePrefix + strconv.FormatInt(i, 10)
	} else {
		panic("unexpected nonce variant")
	}
}
