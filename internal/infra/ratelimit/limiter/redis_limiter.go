package limiter

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	rdb          *redis.Client
	keyExpiresIn time.Duration
	mu           sync.Mutex
}

func NewRedisLimiter(rdb *redis.Client) *RedisLimiter {
	return &RedisLimiter{
		rdb:          rdb,
		keyExpiresIn: time.Minute,
	}
}

func (r *RedisLimiter) Quota(ctx context.Context, key string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pipe := r.rdb.Pipeline()
	quota := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Second)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return quota.Val(), nil
}
