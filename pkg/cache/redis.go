package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/BrondoL/captcha/config"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	cache *cache.Cache
}

func NewRedis(cfg config.Config) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.REDIS_HOST, cfg.REDIS_PORT),
		Password: cfg.REDIS_PASS,
		DB:       0,
	})

	client := cache.New(&cache.Options{
		Redis: rdb,
	})

	return &redisCache{
		cache: client,
	}
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	i := cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	}

	return c.cache.Set(&i)
}

func (c *redisCache) Get(ctx context.Context, key string, value interface{}) error {
	err := c.cache.Get(ctx, key, value)
	if err == cache.ErrCacheMiss {
		return nil
	}

	return err
}
