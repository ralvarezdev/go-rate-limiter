package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	godatabases "github.com/ralvarezdev/go-databases"
	gostringsadd "github.com/ralvarezdev/go-strings/add"
)

type (
	// DefaultRateLimiter struct
	DefaultRateLimiter struct {
		redisClient *redis.Client
		limit       int
		period      time.Duration
	}
)

// NewDefaultRateLimiter creates a new rate limiter
//
// Parameters:
//
//   - redisClient: Redis client
//   - limit: Maximum number of requests allowed within the specified period
//   - period: Time duration for the rate limit window
//
// Returns:
//
//   - *DefaultRateLimiter: Pointer to the created DefaultRateLimiter instance
//   - error: Error if the Redis client is nil
func NewDefaultRateLimiter(
	redisClient *redis.Client,
	limit int,
	period time.Duration,
) (*DefaultRateLimiter, error) {
	// Check if the Redis client is nil
	if redisClient == nil {
		return nil, godatabases.ErrNilConnection
	}

	return &DefaultRateLimiter{
		redisClient,
		limit,
		period,
	}, nil
}

// GetKey gets the rate limiter key
//
// Parameters:
//
//   - ip: IP address of the client
//
// Returns:
//
//   - string: Rate limiter key
func (d DefaultRateLimiter) GetKey(ip string) string {
	return gostringsadd.Prefixes(ip, KeySeparator, KeyPrefix)
}

// SetInitialValue sets the initial value for the given key
//
// Parameters:
//
//   - key: Key to set the initial value for
//
// Returns:
//
//   - error: Error if the operation fails
func (d DefaultRateLimiter) SetInitialValue(key string) error {
	_, err := d.redisClient.Set(context.Background(), key, 1, d.period).Result()
	return err
}

// Limit limits the rate of requests
//
// Parameters:
//
//   - ip: IP address of the client
//
// Returns:
//
//   - error: Error if the rate limit is exceeded or if the operation fails
func (d DefaultRateLimiter) Limit(ip string) error {
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
