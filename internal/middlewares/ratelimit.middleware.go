package middlewares

import "github.com/gin-gonic/gin"

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add rate limiting logic here
		c.Next()
	}
}
