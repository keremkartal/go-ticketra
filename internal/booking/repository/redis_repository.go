package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	AcquireLock(key string, ttl time.Duration) (bool, error)
	ReleaseLock(key string) error
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) RedisRepository {
	return &redisRepository{client: client}
}

func (r *redisRepository) AcquireLock(key string, ttl time.Duration) (bool, error) {
	ctx := context.Background()
	

	success, err := r.client.SetNX(ctx, key, "locked", ttl).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

func (r *redisRepository) ReleaseLock(key string) error {
	ctx := context.Background()
	return r.client.Del(ctx, key).Err()
}