//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/controllers"
	"github.com/kien14502/ecommerce-be/internal/services"
	"github.com/kien14502/ecommerce-be/internal/sse"
)

type SSEBundle struct {
	SSEController       *controllers.SSEController
	NotificationService *services.NotificationService
}

func InitSseRouterHandler() *SSEBundle {
	wire.Build(
		// external
		wire.Value(global.Rdb),

		// core
		sse.NewManager,

		// services
		services.NewNotificationService,

		// controllers
		controllers.NewSSEController,

		wire.Struct(new(SSEBundle), "*"),
	)
	return nil
}
