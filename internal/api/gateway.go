package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	fine "github.com/TheKodeToad/fline/internal"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/go-chi/chi/v5"
)

func gatewayRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Get("/", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		return struct {
			URL string `json:"url"`
		}{fine.GatewayURL(r.Host)}, nil
	}))

	router.Get("/bot", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		fluxerHeaders, err := headersToFluxer(r.Header)
		if err != nil {
			return nil, fmt.Errorf("failed to convert header to fluxer: %w", err)
		}

		fluxerResp, err := client.Do(
			(&http.Request{
				Header: fluxerHeaders,
				URL:    formatFluxerURL(conf, "/gateway/bot"),
			}).WithContext(r.Context()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to perform fluxer request: %w", err)
		}
		writeDiscordHeaders(w.Header(), fluxerResp.Header)

		errResp, err := convFluxerErrorResponse(fluxerResp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert fluxer response: %w", err)
		} else if errResp != nil {
			return errResp, nil
		}

		var inInfo discord.GatewayBotInfo
		err = json.NewDecoder(fluxerResp.Body).Decode(&inInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to decode fluxer response: %w", err)
		}

		outInfo := inInfo
		outInfo.URL = fine.GatewayURL(r.Host)
		return outInfo, nil
	}))

	return router
}
