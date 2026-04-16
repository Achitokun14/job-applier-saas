package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"job-applier-backend/internal/cache"
	"job-applier-backend/internal/models"
)

// SubscriptionTierKey is the context key used to store the user's subscription tier.
const SubscriptionTierKey ContextKey = "subscriptionTier"

// tierOrder defines the ordering of subscription tiers for comparison.
var tierOrder = map[string]int{
	"free":       0,
	"pro":        1,
	"enterprise": 2,
}

// SubscriptionMiddleware looks up the user's subscription tier and injects it into the
// request context. It caches the tier in Redis for 5 minutes to avoid repeated DB lookups.
func SubscriptionMiddleware(db *gorm.DB, c cache.Cache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(ContextKey("userID")).(uint)
			if !ok {
				// Fall back to the handlers' context key type.
				type handlerContextKey string
				userID, ok = r.Context().Value(handlerContextKey("userID")).(uint)
				if !ok {
					next.ServeHTTP(w, r)
					return
				}
			}

			tier := getUserTierCached(r.Context(), db, c, userID)
			ctx := context.WithValue(r.Context(), SubscriptionTierKey, tier)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireTier returns middleware that blocks access if the user's subscription tier
// is below the specified minimum tier. Tier ordering: free < pro < enterprise.
func RequireTier(minTier string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tier, _ := r.Context().Value(SubscriptionTierKey).(string)
			if tier == "" {
				tier = "free"
			}

			minLevel, ok := tierOrder[minTier]
			if !ok {
				minLevel = 0
			}

			currentLevel, ok := tierOrder[tier]
			if !ok {
				currentLevel = 0
			}

			if currentLevel < minLevel {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusPaymentRequired)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": fmt.Sprintf("This feature requires a %s subscription or higher.", minTier),
					"tier":  tier,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getUserTierCached retrieves the user's subscription tier, using Redis cache with a
// 5-minute TTL to reduce database lookups.
func getUserTierCached(ctx context.Context, db *gorm.DB, c cache.Cache, userID uint) string {
	cacheKey := fmt.Sprintf("subscription_tier:%d", userID)

	// Try cache first.
	if c != nil {
		cached, err := c.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			return cached
		}
	}

	// Look up from DB.
	tier := "free"
	var sub models.Subscription
	if err := db.Where("user_id = ? AND status IN ?", userID, []string{"active", "past_due"}).First(&sub).Error; err == nil {
		if sub.Tier != "" {
			tier = sub.Tier
		}
	}

	// Cache the result for 5 minutes.
	if c != nil {
		_ = c.Set(ctx, cacheKey, tier, 5*time.Minute)
	}

	return tier
}
