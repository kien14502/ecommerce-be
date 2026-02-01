package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/global"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		origin := c.Request.Header.Get("Origin")

		allowedOrigins := global.Config.Server.Cors

		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

		// Preflight request
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
