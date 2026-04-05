package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/internal/middlewares"
	"github.com/kien14502/ecommerce-be/internal/wire"
)

func AuthRouter(g *gin.RouterGroup) {
	userController, _ := wire.InitUserRouterHandler()

	auth := g.Group("/auth")
	{
		// Define auth routes here
		auth.POST("/register", userController.Register)
		auth.POST("/verify-otp", userController.VerifyOtp)
		auth.POST("/resend-verify-email", userController.VerifyOtp)
		auth.POST("/login", userController.Login)
		// auth.POST("/oauth-google")
		// Login with Google
	}

	privateRouter := auth.Use(middlewares.AuthMiddleware())
	{
		privateRouter.GET("/me", userController.GetMe)
		privateRouter.POST("/logout", userController.Logout)
	}

}
