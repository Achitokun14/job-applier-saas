package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

// New returns a configured zerolog.Logger.
// In production it outputs structured JSON; in development it uses a
// coloured, human-friendly console writer.
func New(appEnv string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	var lg zerolog.Logger

	if appEnv == "production" {
		lg = zerolog.New(os.Stdout).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Caller().
			Str("env", appEnv).
			Logger()
	} else {
		writer := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		lg = zerolog.New(writer).
			Level(zerolog.DebugLevel).
			With().
			Timestamp().
			Caller().
			Str("env", appEnv).
			Logger()
	}

	return lg
}

// responseWriter wraps http.ResponseWriter to capture status code and bytes written.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	bytes       int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

// RequestLogger returns a Chi-compatible middleware that logs every
// HTTP request with method, URL, status, duration, request_id, and
// remote_addr.
func RequestLogger(lg zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := wrapResponseWriter(w)

			reqID := middleware.GetReqID(r.Context())

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			lg.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", wrapped.status).
				Dur("duration", duration).
				Str("request_id", reqID).
				Str("remote_addr", r.RemoteAddr).
				Int("bytes", wrapped.bytes).
				Msg("request completed")
		})
	}
}
