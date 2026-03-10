package utils

import (
	"context"
	"fmt"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/pkg/response"
	"github.com/redis/go-redis/v9"
)

func OtpHandler(ctx context.Context, hashEmail string) (string, error) {
	otpFound, err := global.Rdb.Get(ctx, hashEmail).Result()
	switch {
	case err == redis.Nil:
		fmt.Println("Key does not exist")
		return "", nil
	case err != nil:
		fmt.Println("Get failed:", err)
		return response.ErrInvalidEmail, err
	case otpFound != "":
		return response.ErrOTPExisted, fmt.Errorf("")
	}

	return "", nil
}
