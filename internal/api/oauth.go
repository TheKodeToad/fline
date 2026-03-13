package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TheKodeToad/fine/internal/config"
	"github.com/TheKodeToad/fine/internal/convert"
	"github.com/TheKodeToad/fine/internal/fluxer"
	"github.com/go-chi/chi/v5"
)

func oauthRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Get("/applications/@me", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		fluxerResp, err := client.Do(
			(&http.Request{
				Header: headersToFluxer(r.Header),
				URL:    formatFluxerURL(conf, "/applications/@me"),
			}).WithContext(r.Context()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to request application: %w", err)
		}
		headersToDiscord(w.Header(), fluxerResp.Header)

		errResp, err := convFluxerErrorResponse(fluxerResp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert application error response: %w", err)
		} else if errResp != nil {
			return errResp, nil
		}

		var inApp fluxer.Application
		err = json.NewDecoder(fluxerResp.Body).Decode(&inApp)
		if err != nil {
			return nil, fmt.Errorf("failed to decode user response: %w", err)
		}

		outApp := convert.ApplicationToDiscord(inApp)
		return outApp, nil
	}))

	return router
}
