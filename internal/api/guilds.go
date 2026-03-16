package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/go-chi/chi/v5"
)

func guildsRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Put("/{guild_id}/bans/{user_id}", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		fluxerHeaders, err := headersToFluxer(r.Header)
		if err != nil {
			return nil, fmt.Errorf("failed to convert header to fluxer: %w", err)
		}

		var inCreate discord.GuildBanCreate
		err = json.NewDecoder(r.Body).Decode(&inCreate)
		if !errors.Is(err, io.EOF) && err != nil {
			if errResp := makeUnmarshalErrorResponse(err); errResp != nil {
				return errResp, nil
			} else {
				return nil, fmt.Errorf("failed to decode payload: %w", err)
			}
		}

		outCreate := convert.GuildBanCreateToFluxer(inCreate)
		if len(fluxerHeaders.Values("X-Audit-Log-Reason")) != 0 {
			// NOTE: fluxer has a separate message for the audit log and ban list
			// this replicates the discord behaviour of them being the same
			outCreate.Reason = fluxerHeaders.Get("X-Audit-Log-Reason")
		}

		fluxerPayload, err := json.Marshal(outCreate)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal converted payload: %w", err)
		}

		fluxerResp, err := client.Do(
			(&http.Request{
				Method: "PUT",
				Body:   io.NopCloser(bytes.NewReader(fluxerPayload)),
				Header: fluxerHeaders,
				URL:    formatFluxerURL(conf, "/guilds/%s/bans/%s", r.PathValue("guild_id"), r.PathValue("user_id")),
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

		return apiNoContentResponse{}, nil
	}))

	return router
}
