package apiroutes

import (
	"bytes"
	"io"
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

	router.Method("PATCH", "/{channel_id}", api.ProxyHandler[any, fluxer.Channel]{
		Conf:   conf,
		Client: client,
		DecodeRequest: func(req *http.Request) (any, error) {
			return io.ReadAll(req.Body)
		},
		MapRequest: func(body any) (any, error) {
			return body, nil
		},
		EncodeRequest: func(body any, req *http.Request) error {
			req.Body = io.NopCloser(bytes.NewReader(body.([]byte)))
			return nil
		},
		MapResponse: func(inChannel fluxer.Channel) (any, error) {
			outChannel, ok := convert.ChannelToDiscord(inChannel)
			if !ok {
				// FIXME: better behaviour
				return nil, api.Error{
					APIError: discord.APIError{
						Code:    discord.APIErrorUnknownChannel,
						Message: "This Channel Uh... Totally Never Existed",
					},
					Status: http.StatusNotFound,
				}
			}

			return outChannel, nil
		},
	})

	router.Method("DELETE", "/{channel_id}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("GET", "/{channel_id}/webhooks", api.ProxyHandler[any, []fluxer.Webhook]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/webhooks",
		MapResponse: func(inWebhooks []fluxer.Webhook) (any, error) {
			outWebhooks := make([]discord.Webhook, 0, len(inWebhooks))
			for _, webhook := range inWebhooks {
				outWebhooks = append(outWebhooks, convert.WebhookToDiscord(webhook))
			}

			return outWebhooks, nil
		},
	})

	router.Method("POST", "/{channel_id}/webhooks", api.ProxyHandler[discord.WebhookCreate, fluxer.Webhook]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/webhooks",
		MapRequest: func(create discord.WebhookCreate) (any, error) {
			return create, nil
		},
		MapResponse: func(webhook fluxer.Webhook) (any, error) {
			return convert.WebhookToDiscord(webhook), nil
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
