package unit

import (
	"context"
	"log"
	"log/slog"
	"sync"
	"testing"
)

type LogRecord struct {
	Level   slog.Level
	Message string
	Attrs   map[string]any // groups are flattened, e.g. "details.request_id"
}

type LogCapture struct {
	mu      sync.Mutex
	records []LogRecord
}

func (c *LogCapture) Records() []LogRecord {
	c.mu.Lock()
	defer c.mu.Unlock()

	return append([]LogRecord(nil), c.records...)
}

// CaptureLogs routes slog's default logger into memory for the duration of the test.
func CaptureLogs(t *testing.T) *LogCapture {
	t.Helper()

	capture := &LogCapture{}

	previousLogger := slog.Default()
	previousWriter, previousFlags := log.Writer(), log.Flags()

	slog.SetDefault(slog.New(captureHandler{capture}))

	t.Cleanup(func() {
		slog.SetDefault(previousLogger)
		log.SetOutput(previousWriter)
		log.SetFlags(previousFlags)
	})

	return capture
}

// captureHandler ignores WithAttrs and WithGroup, which this codebase does not use.
type captureHandler struct{ capture *LogCapture }

func (captureHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h captureHandler) Handle(_ context.Context, r slog.Record) error {
	attrs := map[string]any{}

	r.Attrs(func(a slog.Attr) bool {
		flatten(attrs, "", a)
		return true
	})

	h.capture.mu.Lock()
	defer h.capture.mu.Unlock()

	h.capture.records = append(h.capture.records, LogRecord{
		Level:   r.Level,
		Message: r.Message,
		Attrs:   attrs,
	})

	return nil
}

func (h captureHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h captureHandler) WithGroup(string) slog.Handler      { return h }

func flatten(out map[string]any, prefix string, a slog.Attr) {
	key := prefix + a.Key

	if a.Value.Kind() == slog.KindGroup {
		for _, nested := range a.Value.Group() {
			flatten(out, key+".", nested)
		}

		return
	}

	out[key] = a.Value.Any()
}
