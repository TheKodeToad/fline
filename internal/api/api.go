package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
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

	if resp.Header.Get("Content-Type") != "application/json" {
		return apiError{status: resp.StatusCode}, nil
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

// makeUnmarshalErrorResponse creates an approapriate response if the passed error indicates malformed JSON.
// If it does not, nil is returned and it is probably better to treat it as an internal server error.
func makeUnmarshalErrorResponse(err error) any {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return apiError{
			APIError: discord.APIError{
				Code:    discord.APIErrorRequestBodyHasInvalidJSON,
				Message: "The request body contains invalid JSON.",
			},
			status: http.StatusBadRequest,
		}
	}

	var fieldErr *json.UnmarshalTypeError
	if errors.As(err, &fieldErr) {
		return apiError{
			APIError: discord.APIError{
				Code:    discord.APIErrorInvalidFormBody,
				Message: "Invalid Form Body",
			},
			status: http.StatusBadRequest,
		}
	}

	return nil
}

type apiError struct {
	discord.APIError
	status int
}

type apiNoContentResponse struct{}

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

		if _, ok := respObject.(apiNoContentResponse); ok {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		resp, err := json.Marshal(respObject)
		if err != nil {
			logger.Warn("failed to marshal response object", slog.Any("err", err))
			return
		}

		w.Header().Add("Content-Type", "application/json")
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
