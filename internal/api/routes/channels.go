package apiroutes

import (
	"net/http"

	"github.com/TheKodeToad/fline/internal/api"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/go-chi/chi/v5"
)

func channelsRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Method("GET", "/{channel_id}", api.ProxyHandler[any, fluxer.Channel]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}",
		MapResponse: func(inChannel fluxer.Channel) (any, error) {
			outChannel, ok := convert.ChannelToDiscord(inChannel)
			if !ok {
				return nil, api.Error{
					APIError: discord.APIError{
						Code:    discord.APIErrorUnknownChannel,
						Message: "Unknown Channel",
					},
					Status: http.StatusNotFound,
				}
			}

			return outChannel, nil
		},
	})

	router.Mount("/{channel_id}/messages", messagesRouter(conf, client))

	router.Method("POST", "/{channel_id}/typing", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/typing",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	return router
}
