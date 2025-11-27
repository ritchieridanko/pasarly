package cache

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ritchieridanko/pasarly/backend/services/auth/configs"
)

type Cache struct {
	config *configs.Cache
	client *redis.Client
}

func NewCache(cfg *configs.Cache, c *redis.Client) *Cache {
	return &Cache{config: cfg, client: c}
}

func (c *Cache) Set(ctx context.Context, key string, value any, d time.Duration) error {
	var e error
	for attempt := 0; attempt < c.config.MaxRetries; attempt++ {
		var err error
		if d <= 0 {
			err = c.client.Set(ctx, key, value, 0).Err()
		} else {
			err = c.client.Set(ctx, key, value, d).Err()
		}

		if err == nil {
			return nil
		}

		e = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, c.config.BaseDelay, attempt); err != nil {
			return err
		}
	}

	return e
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	var e error
	for attempt := 0; attempt < c.config.MaxRetries; attempt++ {
		res, err := c.client.Get(ctx, key).Result()
		if err == nil {
			return res, nil
		}

		e = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, c.config.BaseDelay, attempt); err != nil {
			return "", err
		}
	}

	return "", e
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	var e error
	for attempt := 0; attempt < c.config.MaxRetries; attempt++ {
		res, err := c.client.Exists(ctx, key).Result()
		if err == nil {
			return res > 0, nil
		}

		e = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, c.config.BaseDelay, attempt); err != nil {
			return false, err
		}
	}

	return false, e
}

func (c *Cache) Evaluate(ctx context.Context, hashKey, script string, keys []string, args ...any) (any, error) {
	hash, err := c.Get(ctx, hashKey)
	if err != nil {
		hash, err = c.load(ctx, hashKey, script)
		if err != nil {
			return nil, err
		}
	}

	var e error
	for attempt := 0; attempt < c.config.MaxRetries; attempt++ {
		res, err := c.client.EvalSha(ctx, hash, keys, args...).Result()
		if err == nil {
			return res, nil
		}

		if strings.Contains(err.Error(), "NOSCRIPT") {
			hash, err = c.load(ctx, hashKey, script)
			if err != nil {
				return nil, err
			}
		}

		e = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, c.config.BaseDelay, attempt); err != nil {
			return nil, err
		}
	}

	return nil, e
}

func (c *Cache) load(ctx context.Context, key, script string) (string, error) {
	var e error
	for attempt := 0; attempt < c.config.MaxRetries; attempt++ {
		res, err := c.client.ScriptLoad(ctx, script).Result()
		if err == nil {
			if err := c.Set(ctx, key, res, -1); err != nil {
				return "", err
			}

			return res, nil
		}

		e = err
		if !isRetryable(err) {
			break
		}
		if err := backoffWait(ctx, c.config.BaseDelay, attempt); err != nil {
			return "", err
		}
	}

	return "", e
}
