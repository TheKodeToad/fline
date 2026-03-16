package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(conf *config.Config) chi.Router {
	var client http.Client

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)

	router.Route("/v{version}", func(router chi.Router) {
		router.Mount("/channels", channelsRouter(conf, client))
		router.Mount("/gateway", gatewayRouter(conf, client))
		router.Mount("/guilds", guildsRouter(conf, client))
		router.Mount("/oauth2", oauthRouter(conf, client))
		router.Mount("/users", usersRouter(conf, client))
	})

	router.NotFound(apiHandler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (resp any, err error) {
		logger.Debug("api route not found")
		return nil, apiError{status: http.StatusNotFound}
	}))

	return router
}

func formatFluxerURL(conf *config.Config, format string, a ...any) *url.URL {
	return conf.FluxerAPIURL.JoinPath(
		"v"+conf.FluxerAPIVersion,
		fmt.Sprintf(format, a...),
	)
}

// performFluxerRequest performs a HTTP request to Fluxer applying headers from the passed [http.Request] and to the passed [http.ResponseWriter].
func performFluxerRequest(
	serverW http.ResponseWriter,
	serverR *http.Request,
	client http.Client,
	req *http.Request,
) (*http.Response, error) {
	// it might be a bit confusing to have so many things done by one function
	// TODO: maybe something more declarative (where you just specify transformation functions) would be better
	req = req.WithContext(serverR.Context())

	if req.Header == nil {
		req.Header = http.Header{}
	}
	requestHeadersToFluxer(req.Header, serverR.Header)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 600 {
		// it's an error response
		if resp.Header.Get("Content-Type") != "application/json" {
			return nil, apiError{status: resp.StatusCode}
		}

		var errObject fluxer.APIError
		err := json.NewDecoder(resp.Body).Decode(&errObject)
		if err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}

		return nil, apiError{
			APIError: convert.APIErrorToDiscord(errObject),
			status:   resp.StatusCode,
		}
	}

	responseHeadersToDiscord(serverW.Header(), resp.Header)
	return resp, err
}

func removeIllegalHeaderValueChars(val string) string {
	var builder strings.Builder
	builder.Grow(len(val))

	for _, ch := range []byte(val) {
		// NOTE: based on httpguts.ValidHeaderFieldValue which appears to be what is used internally
		// if the behaviour of this function changes or this does not match, it might make the request fail
		if (ch < ' ' && ch != '\t') || ch == 0x7F {
			builder.WriteByte(' ')
		} else {
			builder.WriteByte(ch)
		}
	}

	return builder.String()
}

func requestHeadersToFluxer(out http.Header, headers http.Header) {
	passthrough := []string{
		"Content-Type",
		"Authorization",
	}
	for _, key := range passthrough {
		if len(headers.Values(key)) != 0 {
			out.Add(key, headers.Get(key))
		}
	}

	if auditLogReason := headers.Get("X-Audit-Log-Reason"); auditLogReason != "" {
		unescaped, err := url.PathUnescape(auditLogReason)
		if err != nil {
			// NOTE: Discord ignores invalid escape sequences
			// this approximates the behaviour
			unescaped = auditLogReason
		}

		unescaped = removeIllegalHeaderValueChars(unescaped)

		out.Set("X-Audit-Log-Reason", unescaped)
	}
}

func responseHeadersToDiscord(out http.Header, headers http.Header) {
	passthrough := []string{
		"X-RateLimit-Limit",
		"X-RateLimit-Remaining",
		"X-RateLimit-Reset",
		"X-RateLimit-Reset-After",
		"X-RateLimit-Bucket",
		"X-RateLimit-Scope",
	}
	for _, key := range passthrough {
		if len(headers.Values(key)) != 0 {
			out.Add(key, headers.Get(key))
		}
	}
}

// mapUnmarshalError creates an approapriate response if the passed error indicates malformed JSON.
func mapUnmarshalError(err error) error {
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

	// FIXME: handle io.EOF, as well as unexpected EOF (which is a different error of type *errors.errorString :/)
	return err
}

type apiError struct {
	discord.APIError
	status int
}

func (e apiError) Error() string {
	return fmt.Sprintf(
		"API Error: '%s' (code: %d; status: %d)",
		e.Message,
		e.Code,
		e.status,
	)
}

type apiNoContentResponse struct{}

type apiHandlerFunc func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (resp any, err error)

func apiHandler(handler apiHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := slog.Default().With(slog.Any("url", r.URL.String()))

		formatStatus := func(status int) string {
			return fmt.Sprintf("%d: %s", status, http.StatusText(status))
		}

		status := http.StatusOK
		respObject, err := handler(logger, w, r)
		if err != nil {
			var apiErr apiError
			if errors.As(err, &apiErr) {
				if apiErr.Message == "" {
					apiErr.Message = formatStatus(apiErr.status)
				}

				status = apiErr.status
				respObject = apiErr
			} else {
				logger.Warn("unexpected error in handler", slog.Any("err", err))

				status = http.StatusInternalServerError
				respObject = apiError{
					APIError: discord.APIError{
						Code:    discord.APIErrorGeneral,
						Message: formatStatus(status),
					},
					status: status,
				}
			}
		} else if _, ok := respObject.(apiNoContentResponse); ok {
			w.WriteHeader(http.StatusNoContent)
			return
		} else if _, ok := respObject.(apiError); ok {
			panic("handler returned apiError as response. it should be returned as the latter value (error).")
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
