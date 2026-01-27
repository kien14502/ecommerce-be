package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/internal/routers"
)

func RouterInit() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		routers.RegisterUserRoutes(api)
	}
	return r
}
