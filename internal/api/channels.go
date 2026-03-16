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

	router.Post("/{id}/messages", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		var inCreate discord.MessageCreate
		err := json.NewDecoder(r.Body).Decode(&inCreate)
		if err != nil {
			if errResp := makeUnmarshalErrorResponse(err); errResp != nil {
				return errResp, nil
			} else {
				return nil, fmt.Errorf("failed to decode message payload: %w", err)
			}
		}

		outCreate := convert.MessageCreateToFluxer(inCreate)

		fluxerPayload, err := json.Marshal(outCreate)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal converted message payload: %w", err)
		}

		fluxerResp, err := client.Do(
			(&http.Request{
				Method: "POST",
				Body:   io.NopCloser(bytes.NewReader(fluxerPayload)),
				Header: headersToFluxer(r.Header),
				URL:    formatFluxerURL(conf, "/channels/%s/messages", r.PathValue("id")),
			}).WithContext(r.Context()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to post message: %w", err)
		}
		headersToDiscord(w.Header(), fluxerResp.Header)

		errResp, err := convFluxerErrorResponse(fluxerResp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert posted message error response: %w", err)
		} else if errResp != nil {
			return errResp, nil
		}

		var inMessage fluxer.Message
		err = json.NewDecoder(fluxerResp.Body).Decode(&inMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to decode posted message response: %w", err)
		}

		outMessage := convert.MessageToDiscord(inMessage)
		return outMessage, nil
	}))

	router.Post("/{id}/typing", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		fluxerResp, err := client.Do(
			(&http.Request{
				Method: "POST",
				Header: headersToFluxer(r.Header),
				URL:    formatFluxerURL(conf, "/channels/%s/typing", r.PathValue("id")),
			}).WithContext(r.Context()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to post message: %w", err)
		}
		headersToDiscord(w.Header(), fluxerResp.Header)

		errResp, err := convFluxerErrorResponse(fluxerResp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert typing error response: %w", err)
		} else if errResp != nil {
			return errResp, nil
		}

		return apiNoContentResponse{}, nil
	}))

	return router
}
