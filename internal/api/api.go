package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/TheKodeToad/fine/internal/config"
	"github.com/TheKodeToad/fine/internal/convert"
	"github.com/TheKodeToad/fine/internal/discord"
	"github.com/TheKodeToad/fine/internal/fluxer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func formatFluxerURL(conf *config.Config, format string, a ...any) *url.URL {
	return conf.FluxerAPIURL.JoinPath(
		"v"+conf.FluxerAPIVersion,
		fmt.Sprintf(format, a...),
	)
}

func headersToFluxer(header http.Header) http.Header {
	return header
}

func headersToDiscord(outHeader http.Header, inHeader http.Header) {
	passthrough := []string{
		"X-RateLimit-Limit",
		"X-RateLimit-Remaining",
		"X-RateLimit-Reset",
		"X-RateLimit-Reset-After",
		"X-RateLimit-Bucket",
		"X-RateLimit-Scope",
	}
	for _, key := range passthrough {
		if len(inHeader.Values(key)) != 0 {
			outHeader.Add(key, inHeader.Get(key))
		}
	}
}

func convFluxerErrorResponse(resp *http.Response) (any, error) {
	if resp.StatusCode < 400 || resp.StatusCode >= 600 {
		return nil, nil
	}

	var errObject fluxer.APIError
	err := json.NewDecoder(resp.Body).Decode(&errObject)
	if err != nil {
		return nil, fmt.Errorf("failed to decode error response: %w", err)
	}

	return apiError{
		APIError: convert.APIErrorToDiscord(errObject),
		status:   resp.StatusCode,
	}, nil
}

type apiError struct {
	discord.APIError
	status int
}

func makeLogger(r *http.Request) *slog.Logger {
	return slog.Default().With(slog.Any("url", r.URL))
}

func apiHandler(handler func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := makeLogger(r)

		formatStatus := func(status int) string {
			return fmt.Sprintf("%d %s", status, http.StatusText(status))
		}

		status := http.StatusOK
		respObject, err := handler(logger, w, r)
		if err != nil {
			slog.Warn("unexpected error", slog.Any("err", err))

			status = http.StatusInternalServerError
			respObject = apiError{
				APIError: discord.APIError{
					Code:    discord.APIErrorGeneral,
					Message: formatStatus(status),
				},
				status: status,
			}
		} else if apiErr, ok := respObject.(apiError); ok {
			if apiErr.Message == "" {
				apiErr.Message = formatStatus(apiErr.status)
			}

			status = apiErr.status
		}

		resp, err := json.Marshal(respObject)
		if err != nil {
			logger.Warn("failed to marshal response object", slog.Any("err", err))
			return
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(status)
		_, err = w.Write(resp)
		if err != nil {
			logger.Warn("error writing response", slog.Any("err", err))
			return
		}
	}
}

func Routes(conf *config.Config) chi.Router {
	var client http.Client

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Route("/v{version}", func(router chi.Router) {
		router.Mount("/channels", channelsRouter(conf, client))
		router.Mount("/gateway", gatewayRouter(conf, client))
		router.Mount("/oauth2", oauthRouter(conf, client))
		router.Mount("/users", usersRouter(conf, client))
	})

	return router
}
