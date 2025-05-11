package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type CodeRepo struct {
	client *redis.Client
}

func NewCodeRepo(r *redis.Client) *CodeRepo {
	return &CodeRepo{
		client: r,
	}
}

func (r *CodeRepo) Set(ctx context.Context, email, code string, ttlSeconds int) error {
	return r.client.Set(ctx, email, code, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *CodeRepo) Get(ctx context.Context, email string) (string, error) {
	code, err := r.client.Get(ctx, email).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("code not found")
	}
	return code, err
}

func (r *CodeRepo) Delete(ctx context.Context, email string) error {
	return r.client.Del(ctx, email).Err()
}
