package main

import (
	"strconv"

	_ "github.com/kien14502/ecommerce-be/docs" // source swagger docs
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/initialize"
	swaggerFiles "github.com/swaggo/files"     // gin-swagger middleware
	ginSwagger "github.com/swaggo/gin-swagger" // swagger embed files
)

// @title           Ecommerce API
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Joker
// @contact.url    http://www.swagger.io/support
// @contact.email  stylishjoker.epu@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api/1.0.0
// @schema http

func main() {
	r := initialize.Run()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := global.Config.Server.Port
	r.Run(":" + strconv.Itoa(port))
}
