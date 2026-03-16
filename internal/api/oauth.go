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
		fluxerResp, err := performFluxerRequest(w, r, client, &http.Request{
			URL:    formatFluxerURL(conf, "/applications/@me"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to request application: %w", err)
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
