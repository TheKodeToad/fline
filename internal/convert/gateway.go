package convert

import "github.com/TheKodeToad/fine/internal/discord"

var validGatewayCloseCodes = [...]bool{
	discord.GatewayClosedUnknownError:         true,
	discord.GatewayClosedUnknownOpcode:        true,
	discord.GatewayClosedDecodeError:          true,
	discord.GatewayClosedNotAuthenticated:     true,
	discord.GatewayClosedAuthenicationFailed:  true,
	discord.GatewayClosedAlreadyAuthenticated: true,
	discord.GatewayClosedInvalidSeq:           true,
	discord.GatewayClosedRateLimited:          true,
	discord.GatewayClosedSessionTimedOut:      true,
	discord.GatewayClosedInvalidShard:         true,
	discord.GatewayClosedShardingRequired:     true,
	discord.GatewayClosedInvalidAPIVersion:    true,
	discord.GatewayClosedInvalidIntents:       true,
	discord.GatewayClosedDisallowedIntents:    true,
}

func isValidGatewayCloseCode(code int) bool {
	if code < 0 || code > len(validGatewayCloseCodes) {
		return false
	}

	return validGatewayCloseCodes[code]
}

func GatewayCloseToDiscord(code int, text string) (int, string) {
	if isValidGatewayCloseCode(code) {
		return code, text
	} else {
		return discord.GatewayClosedUnknownError, text
	}
}
