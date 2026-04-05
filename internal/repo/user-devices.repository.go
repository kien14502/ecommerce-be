package repo

import (
	"context"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/database"
)

type IUserDevicesRepository interface {
	CreateUserDevice(ctx context.Context, body database.CreateDeviceParams) error
	GetUserDevice(ctx context.Context, id string) (*database.UserDevice, error)
	GetListDeviceByUserId(ctx context.Context, userId string) ([]database.UserDevice, error)
	UpdateUserDevice(ctx context.Context, body database.UpdateDeviceByIDAndUserIDParams) error
}

type userDevicesRepository struct {
	sqlc *database.Queries
}

// UpdateUserDevice implements [IUserDevicesRepository].
func (u *userDevicesRepository) UpdateUserDevice(ctx context.Context, body database.UpdateDeviceByIDAndUserIDParams) error {
	err := u.sqlc.UpdateDeviceByIDAndUserID(ctx, body)
	if err != nil {
		return err
	}
	return nil
}

// GetListDeviceByUserId implements [IUserDevicesRepository].
func (u *userDevicesRepository) GetListDeviceByUserId(ctx context.Context, userId string) ([]database.UserDevice, error) {
	listDevice, err := u.sqlc.ListUserDevices(ctx, userId)
	if err != nil {
		return nil, err
	}

	return listDevice, nil
}

// CreateUserDevice implements [IUserDevicesRepository].
func (u *userDevicesRepository) CreateUserDevice(ctx context.Context, body database.CreateDeviceParams) error {
	err := u.sqlc.CreateDevice(ctx, body)

	return err
}

// GetUserDevice implements [IUserDevicesRepository].
func (u *userDevicesRepository) GetUserDevice(ctx context.Context, id string) (*database.UserDevice, error) {
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
