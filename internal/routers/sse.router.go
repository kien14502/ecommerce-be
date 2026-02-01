package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/controllers"
	"github.com/kien14502/ecommerce-be/internal/services"
	"github.com/kien14502/ecommerce-be/internal/sse"
)

func SseRouter(rg *gin.RouterGroup) {
	sseManager := sse.NewManager(global.Rdb)

	// Initialize Services
	notificationService := services.NewNotificationService(sseManager)
	go func() {
		// Simulate sending notifications
		// notificationService.NotifyUser("user123", "new_message", map[string]interface{}{
		// 	"message": "You have a new message!",
		// })
		_ = notificationService
	}()
	// Initialize Controllers
	sseController := controllers.NewSSEController(sseManager)

	sseGroup := rg.Group("sse")
	{
		sseGroup.GET("/events", sseController.StreamEvents)
	}
}
