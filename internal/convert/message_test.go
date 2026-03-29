package convert_test

import (
	"reflect"
	"testing"

	"github.com/TheKodeToad/fline/internal/convert"
	"github.com/TheKodeToad/fline/internal/discord"
	"github.com/TheKodeToad/fline/internal/misc"
	"github.com/TheKodeToad/fline/internal/multipartx"
)

func TestUploadsToFluxer(t *testing.T) {
	expect := func(
		inAttachments []discord.Attachment,
		inFiles []multipartx.InMemoryFile,
		expectedOutAttachments []discord.Attachment,
		expectedOutFiles []multipartx.InMemoryFile,
	) {
		gotOutAttachments, gotOutFiles, ok := convert.UploadsToFluxer(inAttachments, inFiles)

		if !ok {
			t.Errorf("conversion failed")
		}

		if !reflect.DeepEqual(expectedOutAttachments, gotOutAttachments) {
			t.Errorf("converting attachments %+v: expected %+v but got %+v", inAttachments, expectedOutAttachments, gotOutAttachments)
		}

		if !reflect.DeepEqual(expectedOutFiles, gotOutFiles) {
			t.Errorf("converting files %+v: expected %+v but got %+v", inFiles, expectedOutFiles, gotOutFiles)
		}
	}

	expect(
		[]discord.Attachment{
			{ID: 0},
			{ID: 1},
		},
		[]multipartx.InMemoryFile{
			{FieldName: "files[0]", FileName: "abcd", Data: []byte("test")},
			{FieldName: "files[1]", FileName: "efgh", Data: []byte("test2")},
		},
		[]discord.Attachment{
			{ID: 0, Filename: misc.New("abcd")},
			{ID: 1, Filename: misc.New("efgh")},
		},
		[]multipartx.InMemoryFile{
			{FieldName: "files[0]", FileName: "abcd", Data: []byte("test")},
			{FieldName: "files[1]", FileName: "efgh", Data: []byte("test2")},
		},
	)

	expect(
		[]discord.Attachment{
			{ID: 0, Filename: misc.New("coolio.txt")},
		},
		[]multipartx.InMemoryFile{
			{FieldName: "files[0]", FileName: "abcd", Data: []byte("test")},
			{FieldName: "files[0]", FileName: "efgh", Data: []byte("test2")},
		},
		[]discord.Attachment{
			{ID: 0, Filename: misc.New("coolio.txt")},
			{ID: 1, Filename: misc.New("efgh")},
		},
		[]multipartx.InMemoryFile{
			{FieldName: "files[0]", FileName: "abcd", Data: []byte("test")},
			{FieldName: "files[1]", FileName: "efgh", Data: []byte("test2")},
		},
	)

	expect(
		[]discord.Attachment{},
		[]multipartx.InMemoryFile{
			{FieldName: "cheese", FileName: "cheese.txt", Data: []byte("foo")},
			{FieldName: "on_toast", FileName: "on_toast.txt", Data: []byte("bar")},
		},
		[]discord.Attachment{
			{ID: 0, Filename: misc.New("cheese.txt")},
			{ID: 1, Filename: misc.New("on_toast.txt")},
		},
		[]multipartx.InMemoryFile{
			{FieldName: "files[0]", FileName: "cheese.txt", Data: []byte("foo")},
			{FieldName: "files[1]", FileName: "on_toast.txt", Data: []byte("bar")},
		},
	)
}
