package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TheKodeToad/fline/internal/config"
	"github.com/go-chi/chi/v5"
)

func channelsRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Mount("/{channel_id}/messages", messagesRouter(conf, client))

	router.Post("/{channel_id}/typing", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		_, err := performFluxerRequest(w, r, client, &http.Request{
			Method: "POST",
			URL:    formatFluxerURL(conf, "/channels/%s/typing", r.PathValue("channel_id")),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to perform fluxer request: %w", err)
		}

		return apiNoContentResponse{}, nil
	}))

	return router
}
