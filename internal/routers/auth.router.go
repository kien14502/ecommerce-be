package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/internal/middlewares"
)

func AuthRouter(g *gin.RouterGroup) {
	auths := g.Group("/auth")
	{
		// Define auth routes here
		auths.POST("/login", func(c *gin.Context) {
			// Login handler
		})
	}

	privateRouter := auths.Use(middlewares.AuthMiddleware())
	{
		privateRouter.POST("/logout", func(c *gin.Context) {
			// Logout handler
		})
	}

}
