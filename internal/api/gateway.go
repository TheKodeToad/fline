package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TheKodeToad/fine/internal/config"
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/go-chi/chi/v5"
)

func gatewayRoutes(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Get("/bot", func(w http.ResponseWriter, r *http.Request) {
		fluxerResp, err := client.Do(
			(&http.Request{
				Header: forwardHeader(&r.Header),
				URL:    makeFluxerURL("/gateway/bot", conf),
			}).WithContext(r.Context()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to request fluxer gateway info"))
		}

		var info discord.GatewayInfo
		err = json.NewDecoder(fluxerResp.Body).Decode(&info)
		if err != nil {
			panic(fmt.Errorf("failed to decode fluxer gateway info response: %w", err))
		}

		info.URL = "ws://" + r.Host + "/gateway"

		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(info)
		if err != nil {
			panic(fmt.Errorf("failed to write discord user response: %w", err))
		}
	})

	return router
}
