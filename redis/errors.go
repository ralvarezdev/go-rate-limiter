package redis

import "errors"

var (
	TooManyRequestsError = errors.New("too many requests")
	NilRateLimiterError  = errors.New("nil rate limiter")
)
