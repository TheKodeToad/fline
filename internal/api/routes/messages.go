package apiroutes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/TheKodeToad/fline/internal/api"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/multipartx"
	"github.com/go-chi/chi/v5"
)

func messagesRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Method("GET", "/", api.ProxyHandler[string, []fluxer.Message]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages",
		DecodeRequest: func(req *http.Request) (string, error) {
			return req.URL.RawQuery, nil
		},
		MapRequest: func(query string) (any, error) {
			return query, nil
		},
		EncodeRequest: func(query any, req *http.Request) error {
			req.URL.RawQuery = query.(string)
			return nil
		},
		MapResponse: func(inMessages []fluxer.Message) (any, error) {
			outMessages := make([]discord.Message, 0, len(inMessages))
			for _, message := range inMessages {
				outMessages = append(outMessages, convert.MessageToDiscord(message))
			}

			return outMessages, nil
		},
	})

	router.Method("POST", "/", api.ProxyHandler[discord.MessageCreate, fluxer.Message]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages",
		DecodeRequest: func(req *http.Request) (discord.MessageCreate, error) {
			contentType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
			if err != nil {
				return discord.MessageCreate{}, api.ErrInvalidFormBody
			}

			switch contentType {
			case "application/json":
				var create discord.MessageCreate

				err := json.NewDecoder(req.Body).Decode(&create)
				if err != nil {
					return discord.MessageCreate{}, err
				}

				return create, nil
			case "multipart/form-data":
				var create discord.MessageCreate

				reader, err := req.MultipartReader()
				if err != nil {
					return discord.MessageCreate{}, fmt.Errorf("failed to create multipart reader: %w", err)
				}

				form, err := multipartx.ReadInMemory(reader, conf.MaxUploadFiles, conf.MaxUploadFileSize)
				if errors.Is(err, multipartx.ErrTooManyFiles) ||
					errors.Is(err, multipartx.ErrFileTooLarge) {
					return discord.MessageCreate{}, api.ErrInvalidFormBody
				} else if err != nil {
					return discord.MessageCreate{}, fmt.Errorf("failed to read multipart form: %w", err)
				}

				err = create.UnmarshalForm(form)
				if err != nil {
					return discord.MessageCreate{}, api.ErrInvalidFormBody
				}

				return create, nil
			default:
				// TODO: application/x-www-form-urlencoded is also supported... for some reason
				return discord.MessageCreate{}, api.ErrInvalidFormBody
			}
		},
		MapRequest: func(inCreate discord.MessageCreate) (any, error) {
			outCreate, ok := convert.MessageCreateToFluxer(inCreate)
			if !ok {
				return nil, api.ErrInvalidFormBody
			}

			return outCreate, nil
		},
		EncodeRequest: func(body any, req *http.Request) error {
			create := body.(fluxer.MessageCreate)
			if len(create.Files) != 0 {
				var data bytes.Buffer
				writer := multipart.NewWriter(&data)

				err := create.EncodeForm(writer)
				if err != nil {
					return fmt.Errorf("failed to encode multipart form: %w", err)
				}

				err = writer.Close()
				if err != nil {
					return fmt.Errorf("failed to finish multipart form: %w", err)
				}

				contentType := mime.FormatMediaType("multipart/form-data", map[string]string{
					"boundary": writer.Boundary(),
				})
				req.Header.Set("Content-Type", contentType)
				req.Body = io.NopCloser(bytes.NewReader(data.Bytes()))
				return nil
			} else {
				data, err := json.Marshal(create)
				if err != nil {
					return err
				}

				req.Header.Set("Content-Type", "application/json")
				req.Body = io.NopCloser(bytes.NewReader(data))
				return nil
			}
		},
		MapResponse: func(message fluxer.Message) (any, error) {
			return convert.MessageToDiscord(message), nil
		},
	})

	router.Method("GET", "/{message_id}", api.ProxyHandler[any, fluxer.Message]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}",
		MapResponse: func(message fluxer.Message) (any, error) {
			return convert.MessageToDiscord(message), nil
		},
	})

	router.Method("PATCH", "/{message_id}", api.ProxyHandler[discord.MessageEdit, fluxer.Message]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}",
		DecodeRequest: func(req *http.Request) (discord.MessageEdit, error) {
			contentType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
			if err != nil {
				return discord.MessageEdit{}, api.ErrInvalidFormBody
			}

			switch contentType {
			case "application/json":
				var edit discord.MessageEdit

				err := json.NewDecoder(req.Body).Decode(&edit)
				if err != nil {
					return discord.MessageEdit{}, err
				}

				return edit, nil
			case "multipart/form-data":
				var edit discord.MessageEdit

				reader, err := req.MultipartReader()
				if err != nil {
					return discord.MessageEdit{}, fmt.Errorf("failed to create multipart reader: %w", err)
				}

				form, err := multipartx.ReadInMemory(reader, conf.MaxUploadFiles, conf.MaxUploadFileSize)
				if errors.Is(err, multipartx.ErrTooManyFiles) ||
					errors.Is(err, multipartx.ErrFileTooLarge) {
					return discord.MessageEdit{}, api.ErrInvalidFormBody
				} else if err != nil {
					return discord.MessageEdit{}, fmt.Errorf("failed to read multipart form: %w", err)
				}

				err = edit.UnmarshalForm(form)
				if err != nil {
					return discord.MessageEdit{}, api.ErrInvalidFormBody
				}

				return edit, nil
			default:
				// TODO: application/x-www-form-urlencoded is also supported... for some reason
				return discord.MessageEdit{}, api.ErrInvalidFormBody
			}
		},
		MapRequest: func(inEdit discord.MessageEdit) (any, error) {
			outCreate, ok := convert.MessageEditToFluxer(inEdit)
			if !ok {
				return nil, api.ErrInvalidFormBody
			}

			return outCreate, nil
		},
		EncodeRequest: func(body any, req *http.Request) error {
			create := body.(discord.MessageEdit)
			if len(create.Files) != 0 {
				var data bytes.Buffer
				writer := multipart.NewWriter(&data)

				err := create.EncodeForm(writer)
				if err != nil {
					return fmt.Errorf("failed to encode multipart form: %w", err)
				}

				err = writer.Close()
				if err != nil {
					return fmt.Errorf("failed to finish multipart form: %w", err)
				}

				fmt.Println(data.String())

				contentType := mime.FormatMediaType("multipart/form-data", map[string]string{
					"boundary": writer.Boundary(),
				})
				req.Header.Set("Content-Type", contentType)
				req.Body = io.NopCloser(bytes.NewReader(data.Bytes()))
				return nil
			} else {
				data, err := json.Marshal(create)
				if err != nil {
					return err
				}

				req.Header.Set("Content-Type", "application/json")
				req.Body = io.NopCloser(bytes.NewReader(data))
				return nil
			}
		},
		MapResponse: func(message fluxer.Message) (any, error) {
			return convert.MessageToDiscord(message), nil
		},
	})

	router.Method("DELETE", "/{message_id}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("POST", "/bulk-delete", api.ProxyHandler[discord.MessageBulkDelete, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/bulk-delete",
		MapRequest: func(body discord.MessageBulkDelete) (any, error) {
			return convert.MessageBulkDeleteToFluxer(body), nil
		},
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("PUT", "/{message_id}/reactions/{emoji_id}/@me", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}/reactions/{emoji_id}/@me",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("DELETE", "/{message_id}/reactions", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}/reactions",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("DELETE", "/{message_id}/reactions/{emoji_id}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}/reactions/{emoji_id}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("DELETE", "/{message_id}/reactions/{emoji_id}/{user_id}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/channels/{channel_id}/messages/{message_id}/reactions/{emoji_id}/{user_id}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	return router
}
