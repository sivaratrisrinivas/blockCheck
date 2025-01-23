package factory

import (
	"fmt"

	"github.com/sivaratrisrinivas/web3/blockCheck/config"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/cache"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/cache/memory"
	"github.com/sivaratrisrinivas/web3/blockCheck/internal/cache/redis"
)

// NewCache creates a new cache instance based on configuration
func NewCache(cfg *config.Config) (cache.Cache, error) {
	switch cfg.Cache.Type {
	case "redis":
		return redis.NewRedisCache(redis.Config{
			Host:     cfg.Redis.Host,
			Port:     cfg.Redis.Port,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
	case "memory":
		return memory.NewMemoryCache(cfg.Cache.TTL)
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cfg.Cache.Type)
	}
}
