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

func messagesRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Method("GET", "/", api.ProxyHandler[string, []fluxer.Message]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages",
		DecodeRequest: func(req *http.Request) (string, error) {
			return req.URL.RawQuery, nil
		},
		MapRequest: func(query string) (any, error) {
			return query, nil
		},
		EncodeRequest: func(query any, req *http.Request) error {
			req.URL.RawQuery = query.(string)
			return nil
		},
		MapResponse: func(inMessages []fluxer.Message) (any, error) {
			outMessages := make([]discord.Message, 0, len(inMessages))
			for _, message := range inMessages {
				outMessages = append(outMessages, convert.MessageToDiscord(message))
			}

			return outMessages, nil
		},
	})

	router.Method("POST", "/", messageCreateHandler(conf, client, "/channels/{channel_id}/messages"))

	router.Method("GET", "/{message_id}", api.ProxyHandler[any, fluxer.Message]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}",
		MapResponse: func(message fluxer.Message) (any, error) {
			return convert.MessageToDiscord(message), nil
		},
	})

	router.Method("DELETE", "/{message_id}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("POST", "/bulk-delete", api.ProxyHandler[discord.MessageBulkDelete, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/bulk-delete",
		MapRequest: func(body discord.MessageBulkDelete) (any, error) {
			return convert.MessageBulkDeleteToFluxer(body), nil
		},
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("PUT", "/{message_id}/reactions/{emoji_id}/@me", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}/reactions/{emoji_id}/@me",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("DELETE", "/{message_id}/reactions", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}/reactions",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("DELETE", "/{message_id}/reactions/{emoji_id}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}/reactions/{emoji_id}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("DELETE", "/{message_id}/reactions/{emoji_id}/{user_id}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}/reactions/{emoji_id}/{user_id}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	return router
}
