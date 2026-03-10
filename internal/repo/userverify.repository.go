package repo

import (
	"context"
	"sync"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/database"
)

type IUserVerifyRepository interface {
	InsertOTPVerify(ctx context.Context, in database.CreateOTPParams) error
}

type userVerifyRepository struct {
	mu   sync.Mutex
	sqlc *database.Queries
}

// InsertOTPVerify implements [IUserVerifyRepository].
func (u *userVerifyRepository) InsertOTPVerify(ctx context.Context, in database.CreateOTPParams) error {
	return u.sqlc.CreateOTP(ctx, in)
}

func NewUserVerifyRepository() IUserVerifyRepository {
	return &userVerifyRepository{
		sqlc: database.New(global.Mdbc),
	}
}
