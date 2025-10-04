package redis

type (
	// RateLimiter interface
	RateLimiter interface {
		Limit(ip string) error
	}
)
