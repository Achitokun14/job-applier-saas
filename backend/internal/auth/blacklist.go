package auth

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBlacklist struct {
	client *redis.Client
}

func NewTokenBlacklist(redisAddr string) *TokenBlacklist {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return &TokenBlacklist{client: client}
}

// Blacklist adds a token ID to the blacklist with a TTL matching the token's remaining expiry.
func (tb *TokenBlacklist) Blacklist(tokenID string, expiry time.Duration) error {
	ctx := context.Background()
	return tb.client.Set(ctx, "blacklist:"+tokenID, "1", expiry).Err()
}

// IsBlacklisted checks whether a token ID has been blacklisted.
func (tb *TokenBlacklist) IsBlacklisted(tokenID string) bool {
	ctx := context.Background()
	val, err := tb.client.Exists(ctx, "blacklist:"+tokenID).Result()
	if err != nil {
		return false
	}
	return val > 0
}
