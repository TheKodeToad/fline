package apiroutes

import (
	"net/http"

	"github.com/TheKodeToad/fline/internal/api"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/go-chi/chi/v5"
)

func channelsRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Mount("/{channel_id}/messages", messagesRouter(conf, client))

	router.Method("POST", "/{channel_id}/typing", api.ProxyHandler[any, api.NoContentResponse]{
		Conf: conf,
		Client: client,
		Path: "/channels/{channel_id}/typing",
		DecodeResponse: api.ExpectNoContentResponse,
	})

	return router
}
