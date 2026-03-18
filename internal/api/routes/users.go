package apiroutes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/TheKodeToad/fline/internal/api"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/go-chi/chi/v5"
)

func usersRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Method("GET", "/@me", api.ProxyHandler[any, fluxer.UserPrivate]{
		Conf:   conf,
		Client: client,
		Path:   "/users/@me",
		MapResponse: func(user fluxer.UserPrivate) (any, error) {
			return convert.UserPrivateToDiscord(user), nil
		},
	})

	router.Method("POST", "/@me/channels", api.ProxyHandler[[]byte, fluxer.Channel]{
		Conf:   conf,
		Client: client,
		Path:   "/users/@me/channels",
		DecodeRequest: func(req *http.Request) ([]byte, error) {
			return io.ReadAll(req.Body)
		},
		MapRequest: func(body []byte) (any, error) {
			return body, nil
		},
		EncodeRequest: func(body any, req *http.Request) error {
			req.Body = io.NopCloser(bytes.NewReader(body.([]byte)))
			req.Header.Set("Content-Type", "application/json")
			return nil
		},
		MapResponse: func(inChannel fluxer.Channel) (any, error) {
			outChannel, ok := convert.ChannelToDiscord(inChannel)
			if !ok {
				return nil, fmt.Errorf("bug: dm channel cannot be converted to fluxer")
			}

			return outChannel, nil
		},
	})

	router.Method("GET", "/{user_id}", api.ProxyHandler[any, fluxer.UserPartial]{
		Conf:   conf,
		Client: client,
		Path:   "/users/{user_id}",
		MapResponse: func(user fluxer.UserPartial) (any, error) {
			// check for A deleted (non-existent) user instead of THE deleted user
			if user.Username == fluxer.DeletedUserUsername &&
				user.Discriminator == fluxer.DeletedUserDiscrim &&
				user.ID != fluxer.DeletedUserID {
				return nil, api.Error{
					APIError: discord.APIError{
						Code:    discord.APIErrorUnknownUser,
						Message: "Unknown User",
					},
					Status: http.StatusNotFound,
				}
			}

			return convert.UserPartialToDiscord(user), nil
		},
	})

	return router
}
