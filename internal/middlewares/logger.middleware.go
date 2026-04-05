package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/global"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		userID, _ := c.Get("userID")

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if query != "" {
			path = path + "?" + query
		}

		switch {
		case statusCode >= 500:
			global.Logger.Error("Server error",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", statusCode),
				zap.Duration("latency", latency),
				zap.String("ip", ip),
				zap.Any("userID", userID),
				zap.String("userAgent", userAgent),
				zap.String("error", errorMessage),
			)

		case statusCode >= 400:
			global.Logger.Warn("Client error",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", statusCode),
				zap.Duration("latency", latency),
				zap.String("ip", ip),
				zap.Any("userID", userID),
				zap.String("error", errorMessage),
			)

		default:
			global.Logger.Info("Request",
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", statusCode),
				zap.Duration("latency", latency),
				zap.String("ip", ip),
				zap.Any("userID", userID),
			)
		}
	}
}
