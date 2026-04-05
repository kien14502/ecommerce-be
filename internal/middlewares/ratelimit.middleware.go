package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/pkg/response"
)

const (
	maxRequests = 100
	windowTime  = 1 * time.Minute
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		count, err := global.Rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Error(response.ErrInternalServer)
			c.Abort()
			return
		}

		if count == 1 {
			global.Rdb.Expire(ctx, key, windowTime)
		}

		if count > maxRequests {

			ttl, _ := global.Rdb.TTL(ctx, key).Result()
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", int(ttl.Seconds())))
			c.Error(response.ErrTooManyRequests)
			c.Abort()
			return
		}

		remaining := maxRequests - int(count)
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		c.Next()
	}
}
