package apiroutes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/TheKodeToad/fline/internal/api"
	"github.com/TheKodeToad/fline/internal/config"
	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/fluxer"
	"github.com/TheKodeToad/fline/internal/misc"
	"github.com/disgoorg/snowflake/v2"
	"github.com/go-chi/chi/v5"
)

func decodeDiscordMessageCreateField(result *discord.MessageCreate, key, value string) error {
	switch key {
	case "content":
		result.Content = misc.New(value)
	case "nonce":
		result.Nonce = misc.New(discord.NonceFromString(value))
	case "tts":
		switch value {
		case "false":
			result.TTS = misc.New(false)
		case "true":
			result.TTS = misc.New(true)
		default:
			return api.ErrInvalidFormBody
		}
	case "sticker_ids":
		id, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return api.ErrInvalidFormBody
		}

		result.StickerIDs = append(result.StickerIDs, snowflake.ID(id))
	case "payload_json":
		err := json.Unmarshal([]byte(value), result)
		if err != nil {
			return api.ErrInvalidFormBody
		}
	case "flags":
		flags, err := strconv.Atoi(value)
		if err != nil {
			return api.ErrInvalidFormBody
		}

		result.Flags = flags
	case "enforce_nonce":
		switch value {
		case "false":
			result.EnforceNonce = misc.New(false)
		case "true":
			result.EnforceNonce = misc.New(true)
		default:
			return api.ErrInvalidFormBody
		}
	}

	return nil
}

func decodeDiscordMessageCreateFile(result *discord.MessageCreate, part *multipart.Part) error {
	// FIXME: probably denial of service
	data, err := io.ReadAll(part)
	if err != nil {
		return err
	}

	result.Files = append(result.Files, discord.MessageFile{
		FieldName: part.FormName(),
		Filename:  part.FileName(),
		Content:   data,
	})
	return nil
}

func decodeDiscordMessageCreateForm(req *http.Request) (discord.MessageCreate, error) {
	var result discord.MessageCreate

	reader, err := req.MultipartReader()
	if err != nil {
		return discord.MessageCreate{}, err
	}

	for {
		part, err := reader.NextRawPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			// FIXME: probably should return an API response if this is actually a parse error
			return discord.MessageCreate{}, err
		}

		if part.FormName() == "" {
			continue
		}

		if part.FileName() == "" {
			bytes, err := io.ReadAll(part)
			if err != nil {
				return discord.MessageCreate{}, fmt.Errorf("error reading part body: %w", err)
			}

			err = decodeDiscordMessageCreateField(&result, part.FormName(), string(bytes))
			if err != nil {
				return discord.MessageCreate{}, fmt.Errorf("failed to parse field '%s': %w", part.FormName(), err)
			}
		} else {
			err := decodeDiscordMessageCreateFile(&result, part)
			if err != nil {
				return discord.MessageCreate{}, fmt.Errorf("failed to parse file '%s': %w", part.FileName(), err)
			}
		}
	}

	return result, nil
}

func encodeFluxerMessageCreateForm(output io.Writer, create fluxer.MessageCreate) (string, error) {
	writer := multipart.NewWriter(output)

	parseFileID := func(name string) (snowflake.ID, bool) {
		const prefix = "files["
		const suffix = "]"

		if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, suffix) {
			return 0, false
		}

		idStr := name[len(prefix) : len(name)-len(suffix)]

		id, err := strconv.ParseUint(idStr, 10, 64)
		return snowflake.ID(id), err == nil
	}

	oldAttachments := map[snowflake.ID]discord.Attachment{}
	if create.Attachments != nil {
		for _, attachment := range create.Attachments {
			if _, ok := oldAttachments[attachment.ID]; ok {
				// duplicate attachment ID
				return "", api.ErrInvalidFormBody
			}

			oldAttachments[attachment.ID] = attachment
		}
		create.Attachments = []discord.Attachment{}
	}

	for newID, file := range create.Files {
		// NOTE: the field name of every file on Fluxer must match files[ID]
		// Discord, on the other hand allows a non-conforming field name
		// remapping every field name to a new set of sequential IDs is an easy way to conform to Fluxer's requirements without collisions
		// we just must make sure to remap the old IDs to the new IDs in the attachments array

		newAttachment := discord.Attachment{
			ID:       snowflake.ID(newID),
			Filename: &file.Filename,
		}

		if originalID, ok := parseFileID(file.FieldName); ok {
			if attachment, ok := oldAttachments[originalID]; ok {
				delete(oldAttachments, originalID)

				newAttachment = attachment

				newAttachment.ID = snowflake.ID(newID)
				if newAttachment.Filename == nil {
					newAttachment.Filename = &file.Filename
				}

			}
		}

		create.Attachments = append(create.Attachments, newAttachment)

		fileWriter, err := writer.CreateFormFile(fmt.Sprintf("files[%d]", newID), file.Filename)
		if err != nil {
			return "", fmt.Errorf("failed to create file reader for '%s': %w", file.Filename, err)
		}

		_, err = fileWriter.Write(file.Content)
		if err != nil {
			return "", fmt.Errorf("failed to write file '%s': %w", file.Filename, err)
		}
	}

	if len(oldAttachments) != 0 {
		// there are attachments which don't correspond to an actual file ID
		return "", api.ErrInvalidFormBody
	}

	jsonPayloadWriter, err := writer.CreateFormField("payload_json")
	if err != nil {
		return "", fmt.Errorf("failed to create JSON payload writer: %w", err)
	}

	err = json.NewEncoder(jsonPayloadWriter).Encode(create)
	if err != nil {
		return "", fmt.Errorf("failed to encode JSON payload field: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("failed to finish writing form data: %w", err)
	}

	contentType := mime.FormatMediaType("multipart/form-data", map[string]string{
		"boundary": writer.Boundary(),
	})
	return contentType, nil
}

func messagesRouter(conf *config.Config, client http.Client) chi.Router {
	router := chi.NewRouter()

	router.Method("GET", "/", api.ProxyHandler[string, []fluxer.Message]{
		Conf: conf,
		Client: client,
		Path: "/channels/{channel_id}/messages",
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
				create, err := decodeDiscordMessageCreateForm(req)
				if err != nil {
					return discord.MessageCreate{}, fmt.Errorf("failed to decode message form: %w", err)
				}

				return create, nil
			default:
				// TODO: application/x-www-form-urlencoded is also supported... for some reason
				return discord.MessageCreate{}, api.ErrInvalidFormBody
			}
		},
		MapRequest: func(create discord.MessageCreate) (any, error) {
			return convert.MessageCreateToFluxer(create), nil
		},
		EncodeRequest: func(body any, req *http.Request) error {
			create := body.(fluxer.MessageCreate)
			if len(create.Files) != 0 {
				var buf bytes.Buffer
				contentType, err := encodeFluxerMessageCreateForm(&buf, create)
				if err != nil {
					return fmt.Errorf("failed to encode message form: %w", err)
				}

				req.Header.Set("Content-Type", contentType)
				req.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
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
		Conf:           conf,
		Client:         client,
		Path:           "/channels/{channel_id}/messages/{message_id}",
		MapResponse: func(message fluxer.Message) (any, error) {
			return convert.MessageToDiscord(message), nil
		},
	})

	router.Method("DELETE", "/{message_id}", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:           conf,
		Client:         client,
		Path:           "/channels/{channel_id}/messages/{message_id}",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("POST", "/bulk-delete", api.ProxyHandler[discord.MessageBulkDelete, api.EmptyResponse]{
		Conf:           conf,
		Client:         client,
		Path:           "/channels/{channel_id}/messages/bulk-delete",
		MapRequest: func(body discord.MessageBulkDelete) (any, error) {
			return convert.MessageBulkDeleteToFluxer(body), nil
		},
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("PUT", "/{message_id}/reactions/{emoji_id}/@me", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:           conf,
		Client:         client,
		Path:           "/channels/{channel_id}/messages/{message_id}/reactions/{emoji_id}/@me",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	router.Method("DELETE", "/{message_id}/reactions", api.ProxyHandler[any, api.EmptyResponse]{
		Conf:           conf,
		Client:         client,
		Path:           "/channels/{channel_id}/messages/{message_id}/reactions",
		DecodeResponse: func(resp *http.Response) (api.EmptyResponse, error) {
			return api.ExpectEmptyResponse(resp, http.StatusNoContent)
		},
	})

	return router
}
