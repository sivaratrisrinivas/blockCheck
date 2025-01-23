package redis

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/cache/types"
)

type RedisCache struct {
	client *redis.Client
	stats  types.Stats
}

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func NewRedisCache(cfg Config) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		stats:  types.Stats{},
	}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		atomic.AddUint64(&c.stats.Misses, 1)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	atomic.AddUint64(&c.stats.Hits, 1)
	return val, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Clear(ctx context.Context) error {
	return c.client.FlushAll(ctx).Err()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

func (c *RedisCache) GetStats() types.Stats {
	return types.Stats{
		Hits:   atomic.LoadUint64(&c.stats.Hits),
		Misses: atomic.LoadUint64(&c.stats.Misses),
		Keys:   uint64(c.client.DBSize(context.Background()).Val()),
	}
}
