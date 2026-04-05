package services

import (
	"context"
	"fmt"
	"time"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type IRedisService interface {
	SaveRefreshToken(ctx context.Context, userID, deviceID, hash string) error
	GetRefreshToken(ctx context.Context, userID, deviceID string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID, deviceID string) error
	DeleteAllRefreshTokens(ctx context.Context, userID string) error
	BlackListToken(ctx context.Context, jti string, ttl time.Duration) error
	IsTokenBlackList(ctx context.Context, jti string) (bool, error)
	SaveOtp(ctx context.Context, email, hash string) error
	GetOtp(ctx context.Context, email string) (string, error)
	DeleteOtp(ctx context.Context, email string) error
	IncrementLoginAttempts(ctx context.Context, ip string) (int64, error)
	GetLoginAttempts(ctx context.Context, ip string) (int64, error)

	// Pub/Sub
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channel string, handler func(message string)) error
	Unsubscribe(ctx context.Context, channel ...string) error
}

type redisService struct{}

// Publish implements [IRedisService].
func (r *redisService) Publish(ctx context.Context, channel string, message interface{}) error {
	return global.Rdb.Publish(ctx, channel, message).Err()
}

// Subscribe implements [IRedisService].
func (r *redisService) Subscribe(ctx context.Context, channel string, handler func(message string)) error {
	pubsub := global.Rdb.Subscribe(ctx, channel)

	// Kiểm tra subscribe thành công
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("subscribe to channel %s failed: %w", channel, err)
	}

	// Lắng nghe message trong goroutine
	go func() {
		defer pubsub.Close()
		ch := pubsub.Channel()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					global.Logger.Info("pubsub channel closed", zap.String("channel", channel))
					return
				}
				handler(msg.Payload)
			case <-ctx.Done():
				global.Logger.Info("pubsub context done", zap.String("channel", channel))
				return
			}
		}
	}()

	return nil
}

// Unsubscribe implements [IRedisService].
func (r *redisService) Unsubscribe(ctx context.Context, channel ...string) error {
	pubsub := global.Rdb.Subscribe(ctx, channel...)
	defer pubsub.Close()
	return pubsub.Unsubscribe(ctx, channel...)
}

// BlackListToken implements [IRedisService].
func (r *redisService) BlackListToken(ctx context.Context, jti string, ttl time.Duration) error {
	key := fmt.Sprintf("auth:blacklist:%s", jti)
	return global.Rdb.Set(ctx, key, "1", ttl).Err()
}

// DeleteAllRefreshToken implements [IRedisService].
func (r *redisService) DeleteAllRefreshTokens(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("auth:refresh:%s:*", userID)
	keys, err := global.Rdb.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return global.Rdb.Del(ctx, keys...).Err()
}

// DeleteOtp implements [IRedisService].
func (r *redisService) DeleteOtp(ctx context.Context, email string) error {
	key := fmt.Sprintf("auth:otp:%s", email)
	return global.Rdb.Del(ctx, key).Err()
}

// DeleteRefreshToken implements [IRedisService].
func (r *redisService) DeleteRefreshToken(ctx context.Context, userID string, deviceID string) error {
	key := fmt.Sprintf("auth:refresh:%s:%s", userID, deviceID)
	return global.Rdb.Del(ctx, key).Err()
}

// GetLoginAttempts implements [IRedisService].
func (r *redisService) GetLoginAttempts(ctx context.Context, ip string) (int64, error) {
	key := fmt.Sprintf("auth:ratelimit:login:%s", ip)
	val, err := global.Rdb.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// GetOtp implements [IRedisService].
func (r *redisService) GetOtp(ctx context.Context, email string) (string, error) {
	key := fmt.Sprintf("auth:otp:%s", email)
	return global.Rdb.Get(ctx, key).Result()
}

// SaveOtp implements [IRedisService].
func (r *redisService) SaveOtp(ctx context.Context, email string, hash string) error {
	key := fmt.Sprintf("auth:otp:%s", email)
	return global.Rdb.Set(ctx, key, hash, 5*time.Minute).Err()
}

// GetRefreshToken implements [IRedisService].
func (r *redisService) GetRefreshToken(ctx context.Context, userID string, deviceID string) (string, error) {
	key := fmt.Sprintf("auth:refresh:%s:%s", userID, deviceID)
	return global.Rdb.Get(ctx, key).Result()
}

// IncrementLoginAttempts implements [IRedisService].
func (r *redisService) IncrementLoginAttempts(ctx context.Context, ip string) (int64, error) {
	key := fmt.Sprintf("auth:ratelimit:login:%s", ip)
	pipe := global.Rdb.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 15*time.Minute)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

// IsTokenBlackList implements [IRedisService].
func (r *redisService) IsTokenBlackList(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("auth:blacklist:%s", jti)
	result, err := global.Rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// SaveRefreshToken implements [IRedisService].
func (r *redisService) SaveRefreshToken(ctx context.Context, userID string, deviceID string, hash string) error {
	key := fmt.Sprintf("auth:refresh:%s:%s", userID, deviceID)
	return global.Rdb.Set(ctx, key, hash, 7*24*time.Hour).Err()
}

func NewRedisService() IRedisService {
	return &redisService{}
}
