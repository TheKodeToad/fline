package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/go-chi/chi/v5"
)

func usersRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Get("/@me", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		fluxerResp, err := performFluxerRequest(w, r, client, &http.Request{
			URL: formatFluxerURL(conf, "/users/@me"),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to perform fluxer request: %w", err)
		}

		var inUser fluxer.UserPrivate
		err = json.NewDecoder(fluxerResp.Body).Decode(&inUser)
		if err != nil {
			return nil, fmt.Errorf("failed to decode fluxer response: %w", err)
		}

		outUser := convert.UserPrivateToDiscord(inUser)
		return outUser, nil
	}))

	router.Get("/{id}", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		fluxerResp, err := performFluxerRequest(w, r, client, &http.Request{
			URL: formatFluxerURL(conf, "/users/%s", r.PathValue("id")),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to perform fluxer request: %w", err)
		}

		var inUser fluxer.UserPartial
		err = json.NewDecoder(fluxerResp.Body).Decode(&inUser)
		if err != nil {
			return nil, fmt.Errorf("failed to decode fluxer response: %w", err)
		}

		// check for A deleted (non-existent) user instead of THE deleted user
		if inUser.Username == fluxer.DeletedUserUsername &&
			inUser.Discriminator == fluxer.DeletedUserDiscrim &&
			inUser.ID != fluxer.DeletedUserID {
			return nil, apiError{
				APIError: discord.APIError{
					Code:    discord.APIErrorUnknownUser,
					Message: "Unknown User",
				},
				status: http.StatusNotFound,
			}
		}

		outUser := convert.UserPartialToDiscord(inUser)
		return outUser, nil
	}))

	return router
}
