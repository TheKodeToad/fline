package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/go-chi/chi/v5"
)

func oauthRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Get("/applications/@me", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		fluxerHeaders, err := headersToFluxer(r.Header)
		if err != nil {
			return nil, fmt.Errorf("failed to convert header to fluxer: %w", err)
		}

		fluxerResp, err := client.Do(
			(&http.Request{
				Header: fluxerHeaders,
				URL:    formatFluxerURL(conf, "/applications/@me"),
			}).WithContext(r.Context()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to request application: %w", err)
		}
		writeDiscordHeaders(w.Header(), fluxerResp.Header)

		errResp, err := convFluxerErrorResponse(fluxerResp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert fluxer error response: %w", err)
		} else if errResp != nil {
			return errResp, nil
		}

		var inApp fluxer.Application
		err = json.NewDecoder(fluxerResp.Body).Decode(&inApp)
		if err != nil {
			return nil, fmt.Errorf("failed to decode fluxer response: %w", err)
		}

		outApp := convert.ApplicationToDiscord(inApp)
		return outApp, nil
	}))

	return router
}
