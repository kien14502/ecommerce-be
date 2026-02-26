package global

import (
	"database/sql"

	"github.com/kien14502/ecommerce-be/pkg/kafka"
	"github.com/kien14502/ecommerce-be/pkg/logger"
	"github.com/kien14502/ecommerce-be/pkg/settings"
	"github.com/redis/go-redis/v9"
)

var (
	Config        settings.Config
	Logger        *logger.LoggerZap
	Mdbc          *sql.DB
	Rdb           *redis.Client // Redis client instance
	KafkaProducer *kafka.Producer
)
