package slogx

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

// LegibleHandler is a [slog.LegibleHandler] which outputs Records in a nicely formatted manner.
// It is designed for development and may be innefficient.
type LegibleHandler struct {
	base      slog.Handler
	output    io.Writer
	bytes     *bytes.Buffer
	bytesLock *sync.Mutex
}

// should implement slog.Handler
var _ slog.Handler = new(LegibleHandler)

func NewLegibleHandler(w io.Writer, opts *slog.HandlerOptions) *LegibleHandler {
	bytes := new(bytes.Buffer)

	return &LegibleHandler{
		base: slog.NewJSONHandler(bytes, opts),
		output: w,
		bytes: bytes,
		bytesLock: new(sync.Mutex),
	}
}

func (h *LegibleHandler) Enabled(ctx context.Context, level slog.Level) bool {
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

func (h *LegibleHandler) Handle(ctx context.Context, record slog.Record) error {
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

func (h *LegibleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var clone = *h
	clone.base = h.base.WithAttrs(attrs)
	return &clone
}

func (h *LegibleHandler) WithGroup(name string) slog.Handler {
	return h.base.WithGroup(name)
}
