package repo

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/database"
	"github.com/kien14502/ecommerce-be/internal/models"
)

type IUserRepository interface {
	FindOne(userID string) (*models.User, error)
	FindAll() []*models.User
	GetUserByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	IsUserExisted(ctx context.Context, email string) (bool, error)
}

type userRepositoryType struct {
	mu   sync.RWMutex
	sqlc *database.Queries
}

// Create implements [IUserRepository].
func (u *userRepositoryType) Create(user *models.User) error {
	panic("unimplemented")
}

// FindAll implements [IUserRepository].
func (u *userRepositoryType) FindAll() []*models.User {
	panic("unimplemented")
}

// FindOne implements [IUserRepository].
func (u *userRepositoryType) FindOne(userID string) (*models.User, error) {
	panic("unimplemented")
}

// GetUserByEmail implements [IUserRepository].
func (u *userRepositoryType) GetUserByEmail(email string) (*models.User, error) {
	panic("unimplemented")
}

// IsUserExisted implements [IUserRepository].
func (u *userRepositoryType) IsUserExisted(ctx context.Context, email string) (bool, error) {
	emailNullString := sql.NullString{
		String: email,
		Valid:  false,
	}
	user, err := u.sqlc.GetUserByEmail(ctx, emailNullString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return user.UserID != 0, nil
}

func NewUserRepository() IUserRepository {
	return &userRepositoryType{
		sqlc: database.New(global.Mdbc),
	}
}
