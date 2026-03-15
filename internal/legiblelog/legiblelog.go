package legiblelog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
)

// inspired by approach in https://dusted.codes/creating-a-pretty-console-logger-using-gos-slog-package
// FIXME: maybe there is something less innefficient

// Handler is a [slog.Handler] which outputs Records in a nicely formatted manner.
// It is designed for development and may be innefficient.
type Handler struct {
	base      *slog.JSONHandler
	output    io.Writer
	bytes     bytes.Buffer
	bytesLock sync.Mutex
}

// should implement slog.Handler
var _ slog.Handler = new(Handler)

func NewHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	var h Handler
	h.base = slog.NewJSONHandler(&h.bytes, opts)
	h.output = w
	return &h
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.base.Enabled(ctx, level)
}

func formatObj(output *strings.Builder, obj map[string]any, prefix string) {
	prefix += "  "

	for key, value := range obj {
		if value == "" {
			continue
		}

		fmt.Fprintf(output, "%s%s: ", prefix, key)

		if inner, ok := value.(map[string]any); ok {
			output.WriteRune('\n')
			formatObj(output, inner, prefix)
		} else {
			line := fmt.Sprintf("%+v", value)
			line = strings.ReplaceAll(line, "\n", "\n"+prefix)
			fmt.Fprintln(output, line)
		}
	}
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	h.bytesLock.Lock()
	defer h.bytesLock.Unlock()
	defer h.bytes.Reset()

	err := h.base.Handle(ctx, record)
	if err != nil {
		return fmt.Errorf("error in base.Handle: %w", err)
	}

	var obj map[string]any
	err = json.Unmarshal(h.bytes.Bytes(), &obj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal log object: %w", err)
	}

	delete(obj, "time")
	delete(obj, "level")
	delete(obj, "msg")

	time := record.Time.Format("15:04:05")

	var output strings.Builder
	fmt.Fprintf(&output, "[%s] [%s] %s\n", time, record.Level, record.Message)
	formatObj(&output, obj, "|")

	h.output.Write([]byte(output.String()))
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h.base.WithAttrs(attrs)
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h.base.WithGroup(name)
}
