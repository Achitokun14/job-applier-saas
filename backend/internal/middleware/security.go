package middleware

import (
	"net/http"

	"github.com/unrolled/secure"

	"job-applier-backend/internal/config"
)

// SecurityHeaders returns middleware that sets security-related HTTP headers.
func SecurityHeaders(cfg *config.Config) func(http.Handler) http.Handler {
	isDev := cfg.AppEnv != "production"

	secureMiddleware := secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		STSSeconds:            31536000, // 1 year
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		IsDevelopment:         isDev,
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := secureMiddleware.Process(w, r)
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// If a redirect was issued by the secure middleware, don't continue.
			if w.Header().Get("Location") != "" {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
