package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/internal/controllers"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	userController := controllers.NewUserController()
	//router
	users := rg.Group("/users")
	{
		users.GET("", userController.GetUser)
	}
}