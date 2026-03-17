package apiroutes

import (
	"log/slog"
	"net/http"

	fline "github.com/TheKodeToad/fline/internal"
	"github.com/TheKodeToad/fline/internal/api"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/go-chi/chi/v5"
)

func gatewayRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Method("GET", "/", api.Handler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		type gatewayInfo struct {
			URL string `json:"url"`
		}
		return gatewayInfo{fline.GatewayURL(r.Host)}, nil
	}))

	router.Method("GET", "/bot", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := api.ProxyHandler[any, discord.GatewayBotInfo]{
			Conf:   conf,
			Client: client,
			Path: "/gateway/bot",
			MapResponse: func(inInfo discord.GatewayBotInfo) (any, error) {
				outInfo := inInfo
				outInfo.URL = fline.GatewayURL(r.Host)
				return outInfo, nil
			},
		}
		h.ServeHTTP(w, r)
	}))

	return router
}
