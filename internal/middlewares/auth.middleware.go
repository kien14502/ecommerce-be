package middlewares

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/services"
	"github.com/kien14502/ecommerce-be/pkg/response"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add authentication logic here
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.Error(response.ErrUnauthorized)
			c.Abort()
			return
		}

		accessToken = strings.TrimPrefix(accessToken, "Bearer ")
		blacklistKey := fmt.Sprintf("blacklist:%s", accessToken)
		exists, err := global.Rdb.Exists(c.Request.Context(), blacklistKey).Result()
		if err == nil && exists > 0 {
			c.Error(response.ErrUnauthorized)
			c.Abort()
			return
		}
		jwtService := services.NewJwtService()
		claims, err := services.IJwtService.ParseAccessToken(jwtService, accessToken)
		if err != nil {
			c.Error(response.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("deviceID", claims.DeviceID)

		c.Next()
	}
}
