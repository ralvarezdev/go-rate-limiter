package redis

import (
	gostringsseparator "github.com/ralvarezdev/go-strings/separator"
)

var (
	// KeyPrefix is the prefix of the rate limiter key
	KeyPrefix = "rate_limiter"

	// KeySeparator is the separator used between the prefix of the rate limiter key
	KeySeparator = gostringsseparator.Dots
)
