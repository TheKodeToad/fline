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
		fluxerResp, err := client.Do(
			(&http.Request{
				Header: headersToFluxer(r.Header),
				URL:    formatFluxerURL(conf, "/users/@me"),
			}).WithContext(r.Context()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to request user: %w", err))
		}
		headersToDiscord(w.Header(), fluxerResp.Header)

		errResp, err := convFluxerErrorResponse(fluxerResp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert user error response: %w", err)
		} else if errResp != nil {
			return errResp, nil
		}

		var inUser fluxer.UserPrivate
		err = json.NewDecoder(fluxerResp.Body).Decode(&inUser)
		if err != nil {
			return nil, fmt.Errorf("failed to decode user response: %w", err)
		}

		outUser := convert.UserPrivateToDiscord(inUser)
		return outUser, nil
	}))

	router.Get("/{id}", apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		fluxerResp, err := client.Do(
			(&http.Request{
				Header: headersToFluxer(r.Header),
				URL:    formatFluxerURL(conf, "/users/%s", r.PathValue("id")),
			}).WithContext(r.Context()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to request user: %w", err))
		}
		headersToDiscord(w.Header(), fluxerResp.Header)

		errResp, err := convFluxerErrorResponse(fluxerResp)
		if err != nil {
			return nil, fmt.Errorf("failed to convert user error response: %w", err)
		} else if errResp != nil {
			return errResp, nil
		}

		var inUser fluxer.UserPartial
		err = json.NewDecoder(fluxerResp.Body).Decode(&inUser)
		if err != nil {
			return nil, fmt.Errorf("failed to decode user response: %w", err)
		}

		// check for A deleted (non-existent) user instead of THE deleted user
		if inUser.Username == fluxer.DeletedUserUsername &&
			inUser.Discriminator == fluxer.DeletedUserDiscrim &&
			inUser.ID != fluxer.DeletedUserID {
			return apiError{
				APIError: discord.APIError{
					Code:    discord.APIErrorUnknownUser,
					Message: "Unknown User",
				},
				status: http.StatusNotFound,
			}, nil
		}

		outUser := convert.UserPartialToDiscord(inUser)
		return outUser, nil
	}))

	return router
}
