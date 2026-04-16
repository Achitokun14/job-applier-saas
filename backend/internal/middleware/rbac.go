package middleware

import (
	"encoding/json"
	"net/http"
)

// ContextKey is a type used for context value keys shared across middleware and handlers.
type ContextKey string

// UserRoleKey is the context key used to store the authenticated user's role.
const UserRoleKey ContextKey = "userRole"

// RequireRole returns middleware that checks if the authenticated user has one of the allowed roles.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user role from context (set by auth middleware)
			userRole, ok := r.Context().Value(UserRoleKey).(string)
			if !ok || userRole == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"error": "Access denied: no role found"})
				return
			}

			// Check if the user's role is in the allowed roles
			for _, role := range roles {
				if userRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Access denied: insufficient permissions"})
		})
	}
}
