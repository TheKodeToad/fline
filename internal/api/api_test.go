package api

import (
	"fmt"
	"net/http"
	"testing"
)

func TestFormatPathValues(t *testing.T) {
	var r http.Request
	r.SetPathValue("joe", "biden")
	r.SetPathValue("yin", "yang")

	const expected = "/joe/biden/yin/yang"
	const format = "/joe/{joe}/yin/{yin}"

	got, err := FormatPathValues(&r, format)
	if err != nil {
		t.Fatalf("formatPathValues failed: %v", err)
	}
	if got != expected {
		t.Errorf("with format '%s': expected '%s' but got '%s'", format, expected, got)
	}
}

func TestFormatPathValuesErrors(t *testing.T) {
	var r http.Request

	expectErr := func(format string, expected string) {
		_, err := FormatPathValues(&r, format)
		errStr := fmt.Sprint(err)
		if errStr != expected {
			t.Errorf("with format '%s': expected error '%s' but got '%s'", format, expected, errStr)
		}
	}

	expectErr("/{}/{}", "no key specified in placeholder at pos 2")
	expectErr("/{joe}{{biden", "excessive opening braces at pos 8")
	expectErr("/{joe}}", "excessive closing braces at pos 7")
	expectErr("/{joe", "expected close brace at pos 6")
}

func BenchmarkFormatPathValues(b *testing.B) {
	var r http.Request
	r.SetPathValue("channel_id", "1474750010096888327")
	r.SetPathValue("message_id", "1483247239782003626")
	r.SetPathValue("emoji_id", "🦋")

	for b.Loop() {
		FormatPathValues(&r, "/channels/{channel_id}/messages/{message_id}/reactions/{emoji_id}/@me")
	}
}

func BenchmarkSprintf(b *testing.B) {
	for b.Loop() {
		_ = fmt.Sprintf(
			"/channels/%s/messages/%s/reactions/%s/@me", 
			"1474750010096888327", 
			"1483247239782003626", 
			"🦋",
		)
	}
}
