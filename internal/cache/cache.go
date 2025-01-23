package cache

import (
	"context"
	"time"

	"github.com/sivaratrisrinivas/web3/blockCheck/internal/cache/types"
)

// Cache defines the interface that all cache implementations must satisfy
type Cache interface {
	// Get retrieves a value from the cache
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value in the cache with optional TTL
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete removes a value from the cache
	Delete(ctx context.Context, key string) error

	// Clear removes all values from the cache
	Clear(ctx context.Context) error

	// Close releases any resources used by the cache
	Close() error

	// GetStats returns cache statistics
	GetStats() types.Stats
}
