//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/kien14502/ecommerce-be/internal/controllers"
	"github.com/kien14502/ecommerce-be/internal/repo"
	"github.com/kien14502/ecommerce-be/internal/services"

	"github.com/google/wire"
)

func InitUserRouterHandler() (*controllers.UserController, error) {
	wire.Build(
		// repo
		repo.NewUserRepository,
		repo.NewUserVerifyRepository,
		repo.NewUserDeviceRepository,

		// service
		services.NewUserService,
		services.NewRedisService,
		services.NewJwtService,

		// controller
		controllers.NewUserController,
	)
	return new(controllers.UserController), nil
}
