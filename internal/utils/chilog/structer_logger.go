package chilog

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/exp/slog"
)

func NewStructuredLogger(handler slog.Handler) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{Logger: handler})
}

type StructuredLogger struct {
	Logger slog.Handler
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	var logFields []slog.Attr
	logFields = append(logFields, slog.String("ts", time.Now().UTC().Format(time.RFC1123)))

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields = append(logFields, slog.String("req_id", reqID))
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	handler := l.Logger.WithAttrs(append(logFields,
		slog.String("http_scheme", scheme),
		slog.String("http_proto", r.Proto),
		slog.String("http_method", r.Method),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent()),
		slog.String("uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI))))

	entry := StructuredLoggerEntry{
		Logger:  slog.New(handler),
		Context: r.Context(),
	}

	return &entry
}

type StructuredLoggerEntry struct {
	Logger  *slog.Logger
	Context context.Context
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger.LogAttrs(l.Context, slog.LevelInfo, "request complete",
		slog.Int("resp_status", status),
		slog.Int("resp_byte_length", bytes),
		slog.Float64("resp_elapsed_ms", float64(elapsed.Nanoseconds())/1000000.0),
	)
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger.LogAttrs(l.Context, slog.LevelInfo, "",
		slog.String("stack", string(stack)),
		slog.String("panic", fmt.Sprintf("%+v", v)),
	)
}
