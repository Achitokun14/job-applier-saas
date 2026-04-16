package errors

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
)

// InitSentry initialises the Sentry SDK.  Call this once at startup.
// If dsn is empty the function is a no-op and returns nil.
func InitSentry(dsn, appEnv string) error {
	if dsn == "" {
		return nil
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		Environment:      appEnv,
		TracesSampleRate: 0.2,
		EnableTracing:    true,
	})
	if err != nil {
		return fmt.Errorf("sentry init: %w", err)
	}
	return nil
}

// FlushSentry drains buffered events.  Call before application exit.
func FlushSentry() {
	sentry.Flush(2 * time.Second)
}

// SentryMiddleware returns a Chi-compatible middleware that captures panics
// and reports them to Sentry.
func SentryMiddleware() func(next http.Handler) http.Handler {
	handler := sentryhttp.New(sentryhttp.Options{
		Repanic: true, // re-panic so Chi's Recoverer can still log the error
	})
	return handler.Handle
}

// CaptureError reports an error to Sentry.  If the context carries a Sentry
// hub (e.g. injected by the middleware) it is used; otherwise the current hub.
func CaptureError(err error, ctx context.Context) {
	if err == nil {
		return
	}

	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
	}
	hub.CaptureException(err)
}
