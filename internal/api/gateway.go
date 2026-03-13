package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	fine "github.com/TheKodeToad/fine/internal"
	"github.com/TheKodeToad/fine/internal/config"
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/go-chi/chi/v5"
)

func gatewayRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Get("/bot", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		fluxerResp, err := client.Do(
			(&http.Request{
				Header: headersToFluxer(r.Header),
				URL:    formatFluxerURL(conf, "/gateway/bot"),
			}).WithContext(r.Context()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to request gateway info: %w", err)
		}
		headersToDiscord(w.Header(), fluxerResp.Header)

		errResp, err := convFluxerErrorResponse(fluxerResp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert gateway info response: %w", err)
		} else if errResp != nil {
			return errResp, nil
		}

		var inInfo discord.GatewayInfo
		err = json.NewDecoder(fluxerResp.Body).Decode(&inInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to decode gateway info response: %w", err)
		}

		outInfo := inInfo
		outInfo.URL = fine.GatewayURL(r.Host)
		return outInfo, nil
	}))

	return router
}
