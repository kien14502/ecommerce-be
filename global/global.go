package global

import (
	"github.com/kien14502/ecommerce-be/pkg/logger"
	"github.com/kien14502/ecommerce-be/pkg/settings"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config settings.Config
	Logger *logger.LoggerZap
	Mdb    *gorm.DB
	Rdb    *redis.Client // Redis client instance
)
