package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/kien14502/ecommerce-be/global"
)

type IUserAuthRepository interface {
	AddOTP(ctx context.Context, email string, otp int, expirationTime int64) error
}

type userAuthRepository struct{}

// AddOTP implements [IUserAuthRepository].
func (u *userAuthRepository) AddOTP(ctx context.Context, email string, otp int, expirationTime int64) error {
	key := fmt.Sprintf("user:%s:otp", email)
	return global.Rdb.SetEx(ctx, key, otp, time.Duration(expirationTime)).Err()
}

func NewUserAuthRepository() IUserAuthRepository {
	return &userAuthRepository{}
}
