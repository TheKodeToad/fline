package apiroutes

import (
	"log/slog"
	"net/http"

	"github.com/TheKodeToad/fline/internal/api"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Router(conf *config.Config) chi.Router {
	var client http.Client

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Route("/v{version}", func(router chi.Router) {
		router.Mount("/channels", channelsRouter(conf, client))
		router.Mount("/gateway", gatewayRouter(conf, client))
		router.Mount("/guilds", guildsRouter(conf, client))
		router.Mount("/oauth2", oauthRouter(conf, client))
		router.Mount("/users", usersRouter(conf, client))
		router.Mount("/webhooks", webhooksRouter(conf, client))
	})

	nf := api.Handler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (resp any, err error) {
		logger.Debug("api route not found")
		return nil, api.Error{Status: http.StatusNotFound}
	})
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		nf.ServeHTTP(w, r)
	})

	return router
}
