package redis

import "errors"

var (
	ErrTooManyRequests = errors.New("too many requests")
	ErrNilRateLimiter  = errors.New("nil rate limiter")
)
