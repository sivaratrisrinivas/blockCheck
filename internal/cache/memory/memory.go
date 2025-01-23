package memory

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sivaratrisrinivas/web3/blockCheck/internal/cache/types"
)

type item struct {
	value      []byte
	expiration time.Time
}

type MemoryCache struct {
	items map[string]item
	mu    sync.RWMutex
	stats types.Stats
}

func NewMemoryCache(defaultTTL time.Duration) (*MemoryCache, error) {
	c := &MemoryCache{
		items: make(map[string]item),
		stats: types.Stats{},
	}

	// Start cleanup goroutine
	go c.cleanup(defaultTTL)

	return c, nil
}

func (c *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, exists := c.items[key]; exists {
		if item.expiration.IsZero() || time.Now().Before(item.expiration) {
			atomic.AddUint64(&c.stats.Hits, 1)
			return item.value, nil
		}
		// Remove expired item
		delete(c.items, key)
	}

	atomic.AddUint64(&c.stats.Misses, 1)
	return nil, nil
}

func (c *MemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration time.Time
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}

	c.items[key] = item{
		value:      value,
		expiration: expiration,
	}
	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]item)
	return nil
}

func (c *MemoryCache) Close() error {
	return nil
}

func (c *MemoryCache) GetStats() types.Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return types.Stats{
		Hits:   atomic.LoadUint64(&c.stats.Hits),
		Misses: atomic.LoadUint64(&c.stats.Misses),
		Keys:   uint64(len(c.items)),
	}
}

func (c *MemoryCache) cleanup(defaultTTL time.Duration) {
	ticker := time.NewTicker(defaultTTL / 2)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if !item.expiration.IsZero() && now.After(item.expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
