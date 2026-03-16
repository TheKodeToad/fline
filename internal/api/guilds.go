package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/go-chi/chi/v5"
)

func guildsRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Get("/{guild_id}", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (resp any, err error) {
		fluxerResp, err := performFluxerRequest(w, r, client, &http.Request{
			URL: formatFluxerURL(conf, "/guilds/%s", r.PathValue("guild_id")),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to perform fluxer request: %w", err)
		}

		var guildIn fluxer.Guild
		err = json.NewDecoder(fluxerResp.Body).Decode(&guildIn)
		if err != nil {
			return nil, fmt.Errorf("failed to decode fluxer response: %w", err)
		}

		guildOut := convert.GuildToDiscord(guildIn)
		return guildOut, nil
	}))

	router.Put("/{guild_id}/bans/{user_id}", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		var inCreate discord.GuildBanCreate
		err := json.NewDecoder(r.Body).Decode(&inCreate)
		if !errors.Is(err, io.EOF) && err != nil {
			if errResp := makeUnmarshalErrorResponse(err); errResp != nil {
				return errResp, nil
			} else {
				return nil, fmt.Errorf("failed to decode payload: %w", err)
			}
		}

		outCreate := convert.GuildBanCreateToFluxer(inCreate)
		if auditLogReason := r.Header.Get("X-Audit-Log-Reason"); auditLogReason != "" {
			// NOTE: fluxer has a separate message for the audit log and ban list
			// this replicates the discord behaviour of them being the same
			unescaped, err := url.PathUnescape(auditLogReason)
			if err != nil {
				unescaped = auditLogReason
			}

			outCreate.Reason = unescaped
		}

		fluxerPayload, err := json.Marshal(outCreate)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal converted payload: %w", err)
		}

		_, err = performFluxerRequest(w, r, client, &http.Request{
			Method: "PUT",
			Body:   io.NopCloser(bytes.NewReader(fluxerPayload)),
			URL:    formatFluxerURL(conf, "/guilds/%s/bans/%s", r.PathValue("guild_id"), r.PathValue("user_id")),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to perform fluxer request: %w", err)
		}

		return apiNoContentResponse{}, nil
	}))

	router.Put("/{guild_id}/members/{user_id}/roles/{role_id}", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (resp any, err error) {
		_, err = performFluxerRequest(w, r, client, &http.Request{
			Method: "PUT",
			URL: formatFluxerURL(conf, "/guilds/%s/members/%s/roles/%s", r.PathValue("guild_id"), r.PathValue("user_id"), r.PathValue("role_id")),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to perform fluxer request: %w", err)
		}

		return apiNoContentResponse{}, nil
	}))

	router.Delete("/{guild_id}/members/{user_id}/roles/{role_id}", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (resp any, err error) {
		_, err = performFluxerRequest(w, r, client, &http.Request{
			Method: "DELETE",
			URL: formatFluxerURL(conf, "/guilds/%s/members/%s/roles/%s", r.PathValue("guild_id"), r.PathValue("user_id"), r.PathValue("role_id")),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to perform fluxer request: %w", err)
		}

		return apiNoContentResponse{}, nil
	}))

	return router
}
