package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	godatabasesredis "github.com/ralvarezdev/go-databases/redis"
	"strconv"
	"time"
)

type (
	// RateLimiter interface
	RateLimiter interface {
		Limit(ip string) error
	}

	// DefaultRateLimiter struct
	DefaultRateLimiter struct {
		redisClient *redis.Client
		limit       int
		period      time.Duration
	}
)

// NewDefaultRateLimiter creates a new rate limiter
func NewDefaultRateLimiter(
	redisClient *redis.Client,
	limit int,
	period time.Duration,
) (*DefaultRateLimiter, error) {
	// Check if the Redis client is nil
	if redisClient == nil {
		return nil, godatabasesredis.NilClientError
	}

	return &DefaultRateLimiter{
		redisClient: redisClient,
		limit:       limit,
		period:      period,
	}, nil
}

// GetKey gets the rate limiter key
func (d *DefaultRateLimiter) GetKey(ip string) string {
	return godatabasesredis.GetKey(ip, RateLimiterPrefix)
}

// SetInitialValue sets the initial value for the given key
func (d *DefaultRateLimiter) SetInitialValue(key string) error {
	_, err := d.redisClient.Set(context.Background(), key, 1, d.period).Result()
	return err
}

// Limit limits the rate of requests
func (d *DefaultRateLimiter) Limit(ip string) error {
	key := d.GetKey(ip)

	// Check the current rate limit
	value, err := d.redisClient.Get(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	// Parse value
	var count int64
	if err == nil {
		count, _ = strconv.ParseInt(value, 10, 64)
	} else {
		// Set the initial value
		return d.SetInitialValue(key)
	}

	// If the rate limit is exceeded, return an error
	if count >= int64(d.limit) {
		return ErrTooManyRequests
	}

	// Increment the request count
	err = d.redisClient.Incr(context.Background(), key).Err()
	if err != nil {
		return err
	}

	// Set the expiration time
	err = d.redisClient.Expire(context.Background(), key, d.period).Err()
	if err != nil {
		return err
	}

	return nil
}
