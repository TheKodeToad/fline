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

	router.Method("POST", "/", api.ProxyHandler[discord.MessageCreate, fluxer.Message] {
		Conf: conf,
		Client: client,
		Path: "/channels/{channel_id}/messages",
		MapRequest: func(create discord.MessageCreate) (any, error) {
			return convert.MessageCreateToFluxer(create), nil
		},
		MapResponse: func(message fluxer.Message) (any, error) {
			return convert.MessageToDiscord(message), nil
		},
	})

	router.Method("DELETE", "/{message_id}", api.ProxyHandler[any, api.NoContentResponse]{
		Conf:           conf,
		Client:         client,
		Path:           "/channels/{channel_id}/messages/{message_id}",
		DecodeResponse: api.ExpectNoContentResponse,
	})

	router.Method("PUT", "/{message_id}/reactions/{emoji_id}/@me", api.ProxyHandler[any, api.NoContentResponse]{
		Conf:           conf,
		Client:         client,
		Path:           "/channels/{channel_id}/messages/{message_id}/reactions/{emoji_id}/@me",
		DecodeResponse: api.ExpectNoContentResponse,
	})

	return router
}
