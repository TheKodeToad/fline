package multipartx

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
)

type InMemoryFile struct {
	FieldName string
	FileName  string
	Data      []byte
}

type InMemoryForm struct {
	Value map[string][]string
	Files []InMemoryFile
}

var (
	ErrTooManyFiles = errors.New("multipart: too many files")
	ErrFileTooLarge = errors.New("multipart: file to large")
)

// ReadInMemory reads a multipart form purely in memory without any temporary files.
// No limits are enforced on fields besides files.
func ReadInMemory(reader *multipart.Reader, maxFiles int, maxFileSize int64) (InMemoryForm, error) {
	result := InMemoryForm{
		Value: map[string][]string{},
	}

	remainingFiles := maxFiles
	for {
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return InMemoryForm{}, err
		}

		formName := part.FormName()
		if formName == "" {
			continue
		}

		fileName := part.FileName()
		if fileName == "" {
			data, err := io.ReadAll(part)
			if err != nil {
				return InMemoryForm{}, fmt.Errorf("error reading value '%s': %s", formName, err)
			}

			result.Value[formName] = append(result.Value[formName], string(data))
		} else {
			remainingFiles--
			if remainingFiles < 0 {
				return InMemoryForm{}, ErrTooManyFiles
			}

			var data bytes.Buffer
			_, err := io.CopyN(&data, part, maxFileSize+1)
			if err == nil {
				// we DID reach the end, so it's too big
				return InMemoryForm{}, ErrFileTooLarge
			} else if !errors.Is(err, io.EOF) {
				return InMemoryForm{}, fmt.Errorf("error reading file '%s'-'%s': %w", formName, fileName, err)
			}

			result.Files = append(result.Files, InMemoryFile{
				FieldName: formName,
				FileName:  fileName,
				Data:      data.Bytes(),
			})
		}
	}

	return result, nil
}
