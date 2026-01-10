package cache

import (
	"Blog-Backend/core"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheDAO struct {
	rdb *redis.Client
}

func NewCacheDAO() *CacheDAO {
	return &CacheDAO{rdb: core.RDB}
}

func (c *CacheDAO) SetJSON(ctx context.Context, key string, v any, ttl time.Duration) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, b, ttl).Err()
}

func (c *CacheDAO) GetJSON(ctx context.Context, key string, out any) (bool, error) {
	b, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal(b, out)
}
