package convert

import (
	"strconv"
	"strings"

	"github.com/TheKodeToad/fline/internal/discord"
)

// since Fluxer only allows string nonce values, we will append this prefix 
// when converting it to Fluxer and strip it when converting to Discord
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
	if nonce.IsString() {
		return nonce.StringValue()
	} else {
		return numericNoncePrefix + strconv.FormatInt(nonce.IntValue(), 10)
	}
}
