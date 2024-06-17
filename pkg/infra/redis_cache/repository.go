package redis_cache

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) Get(ctx context.Context, key string) (int, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	consumed, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}
	return consumed, nil
}

func (r *RedisRepository) Set(ctx context.Context, key string, amount int, ttl time.Duration) error {
	err := r.client.Set(ctx, key, amount, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisRepository) Increment(ctx context.Context, key string) error {
	err := r.client.Incr(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisRepository) IsExpired(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Get(ctx, key+":expired").Result()
	if err != nil {
		return false, err
	}
	expired, err := strconv.ParseBool(result)
	if err != nil {
		return false, err
	}
	return expired, nil
}

func (r *RedisRepository) SetExpired(ctx context.Context, key string, ttl time.Duration) error {
	err := r.client.Set(ctx, key+":expired", true, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
