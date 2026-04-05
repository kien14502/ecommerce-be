//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/kien14502/ecommerce-be/internal/controllers"
	"github.com/kien14502/ecommerce-be/internal/repo"
	"github.com/kien14502/ecommerce-be/internal/services"
)

func InitPostRouterHandler() (*controllers.PostController, error) {
	wire.Build(
		// repo
		repo.NewPostRepository,

		// service
		services.NewPostsService,

		// controller
		controllers.NewPostController,
	)
	return new(controllers.PostController), nil
}
