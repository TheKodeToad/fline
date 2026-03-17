package apiroutes

import (
	"net/http"

	"github.com/TheKodeToad/fline/internal/api"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/go-chi/chi/v5"
)

func oauthRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Method("GET", "/applications/@me", api.ProxyHandler[any, fluxer.Application]{
		Conf: conf,
		Client: client,
		Path: "/applications/@me",
		MapResponse: func(app fluxer.Application) (any, error) {
			return convert.ApplicationToDiscord(app), nil
		},
	})

	return router
}
