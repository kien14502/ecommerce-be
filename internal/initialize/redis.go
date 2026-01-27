package initialize

import (
	"context"
	"fmt"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func RedisInit() {
	redisConfig := global.Config.Redis
	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
		PoolSize: redisConfig.PoolSize,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Error("Redis initialization error", zap.Error(err))
	}

	global.Logger.Info("Redis connected successfully")
	global.Rdb = rdb
	redisExample()
}

func redisExample() {
	err := global.Rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		global.Logger.Error("Failed to set key in Redis", zap.Error(err))
	}

	val, err := global.Rdb.Get(ctx, "key").Result()
	if err != nil {
		global.Logger.Error("Failed to get key from Redis", zap.Error(err))
	} else {
		global.Logger.Info("Redis key value", zap.String("key", val))
	}
}
