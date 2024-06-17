package limiter

import (
	"context"
	"time"
)

type RepositoryInterface interface {
	Get(ctx context.Context, key string) (int, error)
	Set(ctx context.Context, key string, amount int, ttl time.Duration) error
	Increment(ctx context.Context, key string) error
	IsExpired(ctx context.Context, key string) (bool, error)
	SetExpired(ctx context.Context, key string, ttl time.Duration) error
}
