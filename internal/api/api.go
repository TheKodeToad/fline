package api

import (
	"net/http"
	"net/url"

	"github.com/TheKodeToad/fine/internal/config"
	"github.com/go-chi/chi/v5"
)

func makeFluxerURL(path string, conf *config.Config) *url.URL {
	return conf.FluxerAPIURL.JoinPath("v"+conf.FluxerAPIVersion, path)
}

func forwardHeader(header *http.Header) http.Header {
	return http.Header{
		"Authorization": {header.Get("Authorization")},
	}
}

func isErrorResponse(resp http.Response) bool {
	return resp.StatusCode >= 400 && resp.StatusCode < 600
}

func Routes(conf *config.Config) chi.Router {
	var client http.Client

	router := chi.NewRouter()

	router.Route("/v{version}", func(router chi.Router) {
		router.Mount("/users", usersRoutes(conf, client))
		router.Mount("/oauth2", oauthRoutes(conf, client))
		router.Mount("/gateway", gatewayRoutes(conf, client))
	})

	return router
}
