package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/middlewares"
	"github.com/kien14502/ecommerce-be/internal/routers"
)

func RouterInit() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	// r.Use(middlewares.LoggerMiddleware())
	r.Use(middlewares.ErrorHandlerMiddleware())
	r.Use(middlewares.RateLimitMiddleware())
	r.Use(middlewares.TrimBodyMiddleware())

	version := global.Config.Server.Version
	api := r.Group("/api/" + version)

	routers.AuthRouter(api)

	privateRouter := api.Group("/")
	privateRouter.Use(middlewares.AuthMiddleware())

	routers.UserRouter(privateRouter)
	routers.PostsRouter(privateRouter)

	return r
}
