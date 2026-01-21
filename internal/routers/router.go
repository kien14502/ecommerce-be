package routers

import (
	"github.com/gin-gonic/gin"
)

func AppRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		RegisterUserRoutes(api)
	}
	return r
}