package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/TheKodeToad/fine/internal/config"
	"github.com/go-chi/chi/v5"
)

func channelsRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Post("/{id}/messages", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		fluxerResp, err := client.Do(
			(&http.Request{
				Method: "POST",
				Body:   r.Body,
				Header: headersToFluxer(r.Header),
				URL:    formatFluxerURL(conf, "/channels/%s/messages", r.PathValue("id")),
			}).WithContext(r.Context()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to request message post")
		}
		headersToDiscord(w.Header(), fluxerResp.Header)

		errResp, err := convFluxerErrorResponse(fluxerResp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert message post error response: %w", err)
		} else if errResp != nil {
			return errResp, nil
		}

		fluxerBody, err := io.ReadAll(fluxerResp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read message post response: %w", err)
		}

		return json.RawMessage(fluxerBody), nil
	}))

	return router
}
