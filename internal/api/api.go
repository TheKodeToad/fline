package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	fline "github.com/TheKodeToad/fline/internal"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
)

// NoContentResponse simply represents a 204 response with no body.
// It can be returned as the former value of a [HandlerFunc].
type NoContentResponse struct{}

// Error simply contains a Discord API error with a Status.
// It can be returned as the latter value of a [HandlerFunc] to yield an error response.
type Error struct {
	discord.APIError
	Status int `json:"-"`
}

func (e Error) Error() string {
	return fmt.Sprintf(
		"API Error: '%s' (code: %d; status: %d)",
		e.Message,
		e.Code,
		e.Status,
	)
}

type Handler func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := slog.Default().With(slog.Any("route", r.URL.String()))

	formatStatus := func(status int) string {
		return fmt.Sprintf("%d: %s", status, http.StatusText(status))
	}

	status := http.StatusOK
	respObject, err := h(logger, w, r)
	if err != nil {
		var apiErr Error
		if errors.As(err, &apiErr) {
			if apiErr.Message == "" {
				apiErr.Message = formatStatus(apiErr.Status)
			}

			status = apiErr.Status
			respObject = apiErr
		} else {
			logger.Warn("unexpected error in handler", slog.Any("err", err))

			status = http.StatusInternalServerError
			respObject = Error{
				APIError: discord.APIError{
					Code:    discord.APIErrorGeneral,
					Message: formatStatus(status),
				},
				Status: status,
			}
		}
	} else if _, ok := respObject.(NoContentResponse); ok {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if _, ok := respObject.(Error); ok {
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

// MapUnmarshalError maps the error to an appropriate error response if the passed error indicates malformed JSON.
func MapUnmarshalError(err error) error {
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return Error{
			APIError: discord.APIError{
				Code:    discord.APIErrorRequestBodyHasInvalidJSON,
				Message: "The request body contains invalid JSON.",
			},
			Status: http.StatusBadRequest,
		}
	}

	var fieldErr *json.UnmarshalTypeError
	if errors.As(err, &fieldErr) {
		return Error{
			APIError: discord.APIError{
				Code:    discord.APIErrorInvalidFormBody,
				Message: "Invalid Form Body",
			},
			Status: http.StatusBadRequest,
		}
	}

	// FIXME: handle io.EOF, as well as unexpected EOF (which is a different error of type *errors.errorString :/)
	return err
}

// FormatPathValues formats a string replacing {key} placeholders with r.PathValue("key").
func FormatPathValues(r *http.Request, format string) (string, error) {
	var result strings.Builder
	result.Grow(len(format))

	var withinBrace bool
	var chunkStart int

	makeErr := func(msg string, pos int) error {
		return fmt.Errorf("%s at pos %d", msg, pos+1)
	}

	for i, ch := range []byte(format) {
		switch ch {
		case '{':
			// NOTE: appending whole chunks at once seems to be faster than appending bytes one by one
			result.WriteString(format[chunkStart:i])

			if withinBrace {
				return "", makeErr("excessive opening braces", i)
			}

			withinBrace = true
			// FIXME: could this overflow?
			// I mean probably not but I feel uncomfortable just looking at it...
			chunkStart = i + 1
		case '}':
			if !withinBrace {
				return "", makeErr("excessive closing braces", i)
			}

			key := format[chunkStart:i]
			if key == "" {
				return "", makeErr("no key specified in placeholder", chunkStart-1)
			}

			result.WriteString(r.PathValue(key))

			withinBrace = false
			chunkStart = i + 1
		}
	}

	if withinBrace {
		return "", makeErr("expected close brace", len(format))
	}
	result.WriteString(format[chunkStart:])

	return result.String(), nil
}

func FormatFluxerURL(conf *config.Config, r *http.Request, format string) (*url.URL, error) {
	formatted, err := FormatPathValues(r, format)
	if err != nil {
		return nil, err
	}

	return conf.FluxerAPIURL.JoinPath(
		"v"+fline.FluxerAPIVersion,
		formatted,
	), nil
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

func decodeRequestJSON[T any](req *http.Request, optional bool) (T, error) {
	var result T

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		if optional {
			return result, nil
		}

		return result, Error{
			APIError: discord.APIError{
				Code:    discord.APIErrorInvalidFormBody,
				Message: "Invalid Form Body",
			},
			Status: http.StatusBadRequest,
		}
	}

	err := json.NewDecoder(req.Body).Decode(&result)
	if err != nil {
		return result, MapUnmarshalError(err)
	}

	return result, nil
}

// DecodeOptionalRequestJSON decodes the JSON content of req if the content type is properly set.
// Otherwise, an appropriate invalid form body error is return.
func DecodeRequestJSON[T any](req *http.Request) (T, error) {
	return decodeRequestJSON[T](req, false)
}

// DecodeOptionalRequestJSON decodes the JSON content of req if the content type is properly set.
// Otherwise, the zero value of T is returned.
func DecodeOptionalRequestJSON[T any](req *http.Request) (T, error) {
	return decodeRequestJSON[T](req, true)
}

func DecodeResponseJSON[T any](resp *http.Response) (T, error) {
	var result T

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		return result, fmt.Errorf("expected content type to be JSON but got '%s'", contentType)
	}

	err := json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func ExpectNoContentResponse(resp *http.Response) (NoContentResponse, error) {
	if resp.StatusCode != http.StatusNoContent {
		return NoContentResponse{}, fmt.Errorf(
			"expected status %d %s but got %s",
			http.StatusNoContent,
			http.StatusText(http.StatusNoContent),
			resp.Status,
		)
	}

	return NoContentResponse{}, nil
}

// ProxyHandler forwards the request to a Fluxer URL and allows transformation of the response.
type ProxyHandler[ReqBody any, RespBody any] struct {
	Conf   *config.Config
	Client http.Client
	// Path contains a format for the Fluxer route.
	// {name} placeholders can be used to expand path parameters from the original request.
	Path string

	// DecodeRequest is called with the request if MapRequest is not nil.
	// The returned value is passed to MapRequest if decoding is successful.
	// By default [DecodeRequestJSON] is used.
	DecodeRequest func(req *http.Request) (ReqBody, error)
	// EncodeRequest is called with the mapped request if MapRequest is not nil.
	// The returned value is passed to the Fluxer request if encoding is successful.
	// By default [json.MarshalJSON] is used.
	EncodeRequest func(body any) ([]byte, error)
	MapRequest    func(body ReqBody) (any, error)

	// DecodeResponse is called with the response if the status does not represent an error.
	// The returned value is passed to MapResponse if decoding is successful.
	// By default [DecodeResponseJSON] is used.
	DecodeResponse func(resp *http.Response) (RespBody, error)
	MapResponse    func(body RespBody) (any, error)
}

func (opts ProxyHandler[ReqBody, RespBody]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := Handler(func(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (any, error) {
		fluxerURL, err := FormatFluxerURL(opts.Conf, r, opts.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to format fluxer URL: %w", err)
		}

		fluxerReq := &http.Request{
			Method: r.Method,
			URL:    fluxerURL,
			Header: http.Header{},
		}
		fluxerReq = fluxerReq.WithContext(r.Context())
		requestHeadersToFluxer(fluxerReq.Header, r.Header)

		if opts.MapRequest != nil {
			var body ReqBody
			if opts.DecodeRequest != nil {
				body, err = opts.DecodeRequest(r)
			} else {
				body, err = DecodeRequestJSON[ReqBody](r)
			}
			if err != nil {
				return nil, fmt.Errorf("failed to decode request body: %w", err)
			}

			mappedBody, err := opts.MapRequest(body)
			if err != nil {
				return nil, fmt.Errorf("failed to map request body: %w", err)
			}

			var mappedBodyBytes []byte
			if opts.EncodeRequest != nil {
				mappedBodyBytes, err = opts.EncodeRequest(mappedBody)
			} else {
				if logger.Enabled(context.Background(), slog.LevelDebug) {
					// FIXME: maybe changing behaviour with debug log level is a bad idea...
					mappedBodyBytes, err = json.MarshalIndent(mappedBody, "", "  ")

					logger.Debug(
						"sending JSON body", 
						slog.String("fluxerURL", fluxerURL.String()), 
						slog.String("body", string(mappedBodyBytes)),
					)
				} else {
					mappedBodyBytes, err = json.Marshal(mappedBody)
				}
			}
			if err != nil {
				return nil, fmt.Errorf("failed to encode mapped request body: %w", err)
			}

			fluxerReq.Body = io.NopCloser(bytes.NewReader(mappedBodyBytes))
		}

		fluxerResp, err := opts.Client.Do(fluxerReq)
		if err != nil {
			return nil, fmt.Errorf("failed to perform fluxer request: %w", err)
		}
		responseHeadersToDiscord(w.Header(), fluxerResp.Header)

		if fluxerResp.StatusCode >= 400 && fluxerResp.StatusCode < 600 {
			// it's an error response
			if fluxerResp.Header.Get("Content-Type") != "application/json" {
				return nil, Error{Status: fluxerResp.StatusCode}
			}

			var errObject fluxer.APIError
			err := json.NewDecoder(fluxerResp.Body).Decode(&errObject)
			if err != nil {
				return nil, fmt.Errorf("failed to decode error response: %w", err)
			}

			return nil, Error{
				APIError: convert.APIErrorToDiscord(errObject),
				Status:   fluxerResp.StatusCode,
			}
		}

		var decodedResp RespBody
		if opts.DecodeResponse != nil {
			decodedResp, err = opts.DecodeResponse(fluxerResp)
		} else {
			decodedResp, err = DecodeResponseJSON[RespBody](fluxerResp)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to decode non-error response: %w", err)
		}

		if opts.MapResponse != nil {
			mappedResp, err := opts.MapResponse(decodedResp)
			if err != nil {
				return nil, fmt.Errorf("failed to map response: %w", err)
			}

			return mappedResp, nil
		} else {
			return decodedResp, nil
		}
	})
	h.ServeHTTP(w, r)
}
