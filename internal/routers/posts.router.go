package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/internal/wire"
)

func PostsRouter(rg *gin.RouterGroup) {

	postsController, _ := wire.InitPostRouterHandler()
	//router
	// posts := rg.Group("/posts")
	// {
	// 	posts.GET(":id", postsController.GetPosts)
	// }

	user := rg.Group("/user/:id")
	{
		user.GET("/posts", postsController.GetPosts)
	}
}
