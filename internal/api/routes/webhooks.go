package apiroutes

import (
	"net/http"

	"github.com/TheKodeToad/fline/internal/api"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/go-chi/chi/v5"
)

func webhooksRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Method("GET", "/{webhook_id}", api.ProxyHandler[any, fluxer.Webhook]{
		Conf:   conf,
		Client: client,
		Path:   "/webhooks/{webhook_id}",
		MapResponse: func(webhook fluxer.Webhook) (any, error) {
			return convert.WebhookToDiscord(webhook), nil
		},
	})

	router.Method("GET", "/{webhook_id}/{webhook_token}", api.ProxyHandler[any, fluxer.Webhook]{
		Conf:   conf,
		Client: client,
		Path:   "/webhooks/{webhook_id}/{webhook_token}",
		MapResponse: func(inWebhook fluxer.Webhook) (any, error) {
			outWebhook := convert.WebhookToDiscord(inWebhook)
			outWebhook.User = nil

			return outWebhook, nil
		},
	})

	router.Method("DELETE", "/{webhook_id}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/webhooks/{webhook_id}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("DELETE", "/{webhook_id}/{webhook_token}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/webhooks/{webhook_id}/{webhook_token}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method(
		"POST", 
		"/{webhook_id}/{webhook_token}", 
		messageCreateHandler(conf, client, "/webhooks/{webhook_id}/{webhook_token}"),
	)

	return router
}
