package repo

import (
	"context"
	"time"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/database"
)

type IOtpVerificationRepository interface {
	CreateOtp(ctx context.Context, in database.CreateOTPParams) error
	DeleteOtp(ctx context.Context, otp database.DeleteOTPParams) error
}

type otpVerificationRepository struct {
	sqlc *database.Queries
}

// DeleteOtp implements [IOtpVerificationRepository].
func (o *otpVerificationRepository) DeleteOtp(ctx context.Context, in database.DeleteOTPParams) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return o.sqlc.DeleteOTP(ctx, in)
}

// CreateOtp implements [IOtpVerificationRepository].
func (o *otpVerificationRepository) CreateOtp(ctx context.Context, in database.CreateOTPParams) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return o.sqlc.CreateOTP(ctx, in)
}

func NewOtpVerificationRepository() IOtpVerificationRepository {
	return &otpVerificationRepository{
		sqlc: database.New(global.Mdbc),
	}
}
