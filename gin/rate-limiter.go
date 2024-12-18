package gin

import (
	"github.com/gin-gonic/gin"
)

// RateLimiter interface
type RateLimiter interface {
	Limit() gin.HandlerFunc
}
