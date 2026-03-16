package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/go-chi/chi/v5"
)

func channelsRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Post("/{channel_id}/messages", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		var inCreate discord.MessageCreate
		err := json.NewDecoder(r.Body).Decode(&inCreate)
		if err != nil {
			return nil, fmt.Errorf("failed to decode payload: %w", mapUnmarshalError(err))
		}

		outCreate := convert.MessageCreateToFluxer(inCreate)

		fluxerPayload, err := json.Marshal(outCreate)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal converted payload: %w", err)
		}

		fluxerResp, err := performFluxerRequest(w, r, client, &http.Request{
			Method: "POST",
			URL:    formatFluxerURL(conf, "/channels/%s/messages", r.PathValue("channel_id")),
			Body:   io.NopCloser(bytes.NewReader(fluxerPayload)),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to perform fluxer request: %w", err)
		}

		var inMessage fluxer.Message
		err = json.NewDecoder(fluxerResp.Body).Decode(&inMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to decode fluxer response: %w", err)
		}

		outMessage := convert.MessageToDiscord(inMessage)
		return outMessage, nil
	}))

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
