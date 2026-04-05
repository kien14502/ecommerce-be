package initialize

import (
	"context"
	"encoding/json"

	"github.com/kien14502/ecommerce-be/consts"
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/services"
	"go.uber.org/zap"
)

// initialize/pubsub.go
func InitPubSub(redisService services.IRedisService) {
	ctx := context.Background()

	// Subscribe channel OTP
	err := redisService.Subscribe(ctx, consts.ChannelOTP, func(message string) {
		global.Logger.Info("received otp message", zap.String("message", message))

		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(message), &payload); err != nil {
			global.Logger.Error("unmarshal otp message failed", zap.Error(err))
			return
		}
		// Xử lý gửi email OTP
		email := payload["email"].(string)
		otp := payload["otp"].(float64)
		global.Logger.Info("sending otp email",
			zap.String("email", email),
			zap.Float64("otp", otp),
		)
	})
	if err != nil {
		global.Logger.Error("subscribe otp channel failed", zap.Error(err))
	}

	// Subscribe channel logout
	err = redisService.Subscribe(ctx, consts.ChannelUserLogout, func(message string) {
		global.Logger.Info("received logout message", zap.String("message", message))
		// Xử lý logout logic
	})
	if err != nil {
		global.Logger.Error("subscribe logout channel failed", zap.Error(err))
	}
}
