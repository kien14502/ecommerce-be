package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/internal/wire"
)

func UserRouter(rg *gin.RouterGroup) {

	userController, _ := wire.InitUserRouterHandler()
	//router
	users := rg.Group("/users")
	{
		users.GET(":id", userController.GetUser)
	}
}
