package cache

import (
	"context"
	"time"
)

// Cache defines an interface for key-value caching operations.
type Cache interface {
	// Get retrieves a value by key. Returns empty string and nil error on cache miss.
	Get(ctx context.Context, key string) (string, error)

	// Set stores a value with the given key and TTL.
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Delete removes a single key from the cache.
	Delete(ctx context.Context, key string) error

	// DeletePattern removes all keys matching the given glob pattern.
	// Uses SCAN + DEL internally to avoid blocking the Redis server.
	DeletePattern(ctx context.Context, pattern string) error
}
