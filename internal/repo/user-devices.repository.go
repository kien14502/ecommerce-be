package repo

import (
	"context"
	"time"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/database"
)

type IUserDevicesRepository interface {
	CreateUserDevice(ctx context.Context, body database.CreateDeviceParams) error
	GetUserDevice(ctx context.Context, id string) (*database.UserDevice, error)
	GetListDeviceByUserId(ctx context.Context, userId string) ([]database.UserDevice, error)
}

type userDevicesRepository struct {
	sqlc *database.Queries
}

// GetListDeviceByUserId implements [IUserDevicesRepository].
func (u *userDevicesRepository) GetListDeviceByUserId(ctx context.Context, userId string) ([]database.UserDevice, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	listDevice, err := u.sqlc.ListUserDevices(ctx, userId)
	if err != nil {
		return nil, err
	}

	return listDevice, nil
}

// CreateUserDevice implements [IUserDevicesRepository].
func (u *userDevicesRepository) CreateUserDevice(ctx context.Context, body database.CreateDeviceParams) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := u.sqlc.CreateDevice(ctx, body)

	return err
}

// GetUserDevice implements [IUserDevicesRepository].
func (u *userDevicesRepository) GetUserDevice(ctx context.Context, id string) (*database.UserDevice, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userDevice, err := u.sqlc.GetDeviceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &userDevice, nil
}

func NewUserDeviceRepository() IUserDevicesRepository {
	return &userDevicesRepository{
		sqlc: database.New(global.Mdbc),
	}
}
