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

func webhooksRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Method("GET", "/{webhook_id}", api.ProxyHandler[any, fluxer.Webhook]{
		Conf:   conf,
		Client: client,
		Path:   "/webhooks/{webhook_id}",
		MapResponse: func(webhook fluxer.Webhook) (any, error) {
			return convert.WebhookToDiscord(webhook), nil
		},
	})

	router.Method("GET", "/{webhook_id}/{webhook_token}", api.ProxyHandler[any, fluxer.Webhook]{
		Conf:   conf,
		Client: client,
		Path:   "/webhooks/{webhook_id}/{webhook_token}",
		MapResponse: func(inWebhook fluxer.Webhook) (any, error) {
			outWebhook := convert.WebhookToDiscord(inWebhook)
			outWebhook.User = nil

			return outWebhook, nil
		},
	})

	router.Method("DELETE", "/{webhook_id}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/webhooks/{webhook_id}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("DELETE", "/{webhook_id}/{webhook_token}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:   conf,
		Client: client,
		Path:   "/webhooks/{webhook_id}/{webhook_token}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("POST", "/{webhook_id}/{webhook_token}", api.ProxyHandler[discord.WebhookExecute, *fluxer.Message]{
		Conf:   conf,
		Client: client,
		Path:   "/webhooks/{webhook_id}/{webhook_token}",
		DecodeRequest: func(req *http.Request) (discord.WebhookExecute, error) {
			contentType, _, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
			if err != nil {
				return discord.WebhookExecute{}, api.ErrInvalidFormBody
			}

			var exec discord.WebhookExecute

			query := req.URL.Query() 
			if wait := query["wait"]; len(wait) != 0 {
				switch wait[0] {
				case "true":
					exec.Wait = true
				case "false":
					exec.Wait = false
				default:
					return discord.WebhookExecute{}, api.ErrInvalidFormBody
				}
			}

			switch contentType {
			case "application/json":
				err := json.NewDecoder(req.Body).Decode(&exec)
				if err != nil {
					return discord.WebhookExecute{}, err
				}

				return exec, nil
			case "multipart/form-data":
				reader, err := req.MultipartReader()
				if err != nil {
					return discord.WebhookExecute{}, fmt.Errorf("failed to create multipart reader: %w", err)
				}

				form, err := multipartx.ReadInMemory(reader, conf.MaxUploadFiles, conf.MaxUploadFileSize)
				if errors.Is(err, multipartx.ErrTooManyFiles) ||
					errors.Is(err, multipartx.ErrFileTooLarge) {
					return discord.WebhookExecute{}, api.ErrInvalidFormBody
				} else if err != nil {
					return discord.WebhookExecute{}, fmt.Errorf("failed to read multipart form: %w", err)
				}

				err = exec.UnmarshalForm(form)
				if err != nil {
					return discord.WebhookExecute{}, api.ErrInvalidFormBody
				}

				return exec, nil
			default:
				// TODO: application/x-www-form-urlencoded is also supported... for some reason
				return discord.WebhookExecute{}, api.ErrInvalidFormBody
			}
		},
		MapRequest: func(inExec discord.WebhookExecute) (any, error) {
			outCreate, ok := convert.WebhookExecuteToDiscord(inExec)
			if !ok {
				return nil, api.ErrInvalidFormBody
			}

			return outCreate, nil
		},
		EncodeRequest: func(body any, req *http.Request) error {
			exec := body.(discord.WebhookExecute)

			// NOTE: this is just safe because this is a boolean
			// don't copy this for other cases!
			req.URL.RawQuery = fmt.Sprintf("wait=%t", exec.Wait)

			if len(exec.Files) != 0 {
				var data bytes.Buffer
				writer := multipart.NewWriter(&data)

				err := exec.EncodeForm(writer)
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
				data, err := json.Marshal(exec)
				if err != nil {
					return err
				}

				req.Header.Set("Content-Type", "application/json")
				req.Body = io.NopCloser(bytes.NewReader(data))
				return nil
			}
		},
		DecodeResponse: func(resp *http.Response) (*fluxer.Message, error) {
			if resp.StatusCode == http.StatusNoContent {
				return nil, nil
			} else {
				return api.DecodeResponseJSON[*fluxer.Message](resp)
			}
		},
		MapResponse: func(message *fluxer.Message) (any, error) {
			if message == nil {
				return api.EmptyResponse{Status: http.StatusNoContent}, nil
			} else {
				return convert.MessageToDiscord(*message), nil
			}
		},
	})

	return router
}
