package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TheKodeToad/fine/internal/config"
	"github.com/TheKodeToad/fine/internal/convert"
	"github.com/TheKodeToad/fine/internal/fluxer"
	"github.com/go-chi/chi/v5"
)

func usersRoutes(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Get("/@me", func(w http.ResponseWriter, r *http.Request) {
		fluxerResp, err := client.Do(
			(&http.Request{
				Header: forwardHeader(&r.Header),
				URL:    makeFluxerURL("/users/@me", conf),
			}).WithContext(r.Context()),
		)
		if err != nil {
			panic(fmt.Errorf("failed to request fluxer user"))
		}

		var user fluxer.PrivateUser
		err = json.NewDecoder(fluxerResp.Body).Decode(&user)
		if err != nil {
			panic(fmt.Errorf("failed to decode fluxer user response: %w", err))
		}

		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(convert.PrivateUserToDiscord(user))
		if err != nil {
			panic(fmt.Errorf("failed to write discord user response: %w", err))
		}
	})

	return router
}
