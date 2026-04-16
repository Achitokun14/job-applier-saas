package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/httprate"
)

// GlobalRateLimit returns middleware that limits all requests to 100 per minute per IP.
func GlobalRateLimit() func(http.Handler) http.Handler {
	return httprate.LimitByIP(100, 1*time.Minute)
}

// AuthRateLimit returns middleware that limits auth requests to 5 per minute per IP.
func AuthRateLimit() func(http.Handler) http.Handler {
	return httprate.LimitByIP(5, 1*time.Minute)
}
