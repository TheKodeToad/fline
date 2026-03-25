package convert

import (
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
)

func APIErrorCodeToDiscord(code fluxer.APIErrorCode) discord.APIErrorCode {
	switch code {
	case fluxer.APIErrorAccessDenied:
		return discord.APIErrorMissingAccess
	case fluxer.APIErrorAccountTooNewForGuild:
		return discord.APIErrorChannelVerificationToHighToGainAccess
	case fluxer.APIErrorApplicationNotFound:
		return discord.APIErrorUnknownApplication
	case fluxer.APIErrorApplicationNotOwned:
		return discord.APIErrorUnauthorizedApplicationAction
	case fluxer.APIErrorUnknownApplication:
		return discord.APIErrorUnknownApplication
	case fluxer.APIErrorUnknownChannel:
		return discord.APIErrorUnknownChannel
	case fluxer.APIErrorUnknownGuild:
		return discord.APIErrorUnknownGuild
	case fluxer.APIErrorUnknownMember:
		return discord.APIErrorUnknownMember
	case fluxer.APIErrorUnknownWebhook:
		return discord.APIErrorUnknownWebhook
	case fluxer.APIErrorCannotSendEmptyMessage:
		return discord.APIErrorCannotSendEmptyMessage
	case fluxer.APIErrorMissingPermissions:
		return discord.APIErrorLackingPermissions
	// TODO: ...draw the rest of the owl
	default:
		return discord.APIErrorGeneral
	}
}

func APIErrorToDiscord(apiErr fluxer.APIError) discord.APIError {
	return discord.APIError{
		Code: APIErrorCodeToDiscord(apiErr.Code),
		Message: apiErr.Message,
	}

}
