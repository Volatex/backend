package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Redis struct {
	Client       *redis.Client
	connAttempts int
	connTimeout  time.Duration
}

func New(addr, password string, db int, opts ...Option) (*Redis, error) {
	r := &Redis{
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	for _, opt := range opts {
		opt(r)
	}

	r.Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	var err error
	for r.connAttempts > 0 {
		_, err = r.Client.Ping(context.Background()).Result()
		if err == nil {
			break
		}

		log.Printf("Redis is trying to connect, attempts left: %d", r.connAttempts)
		time.Sleep(r.connTimeout)
		r.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("redis - NewRedis - connAttempts == 0: %w", err)
	}

	return r, nil
}

func (r *Redis) Close() error {
	return r.Client.Close()
}
