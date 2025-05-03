package redis

import "time"

type Option func(*Redis)

func ConnAttempts(attempts int) Option {
	return func(r *Redis) {
		r.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(r *Redis) {
		r.connTimeout = timeout
	}
}
