package repo

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/database"
	"github.com/kien14502/ecommerce-be/internal/models"
)

type IUserRepository interface {
	FindOne(ctx context.Context, userID string) (*models.User, error)
	FindAll() []*models.User
	GetUserByEmail(ctx context.Context, email string) (*database.User, error)
	Create(ctx context.Context, in database.CreateUserParams) error
	IsUserExisted(ctx context.Context, email string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (*database.User, error)
	MarkEmailVerified(ctx context.Context, email string) error
}

type userRepositoryType struct {
	mu   sync.RWMutex
	sqlc *database.Queries
}

// MarkEmailVerified implements [IUserRepository].
func (u *userRepositoryType) MarkEmailVerified(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := u.sqlc.MarkEmailVerified(ctx, sql.NullString{String: email, Valid: true})
	return err
}

// GetUserByUsername implements [IUserRepository].
func (u *userRepositoryType) GetUserByUsername(ctx context.Context, username string) (*database.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	user, err := u.sqlc.GetUserByUsername(ctx, sql.NullString{String: username, Valid: true})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create implements [IUserRepository].
func (u *userRepositoryType) Create(ctx context.Context, in database.CreateUserParams) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err := u.sqlc.CreateUser(ctx, in)
	return err
}

// FindAll implements [IUserRepository].
func (u *userRepositoryType) FindAll() []*models.User {
	panic("unimplemented")
}

// FindOne implements [IUserRepository].
func (u *userRepositoryType) FindOne(ctx context.Context, userID string) (*models.User, error) {
	panic("unimplemented")
}

// GetUserByEmail implements [IUserRepository].
func (u *userRepositoryType) GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	user, err := u.sqlc.GetUserByEmail(ctx, sql.NullString{String: email, Valid: true})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// IsUserExisted implements [IUserRepository].
func (u *userRepositoryType) IsUserExisted(ctx context.Context, email string) (bool, error) {
	user, err := u.sqlc.GetUserByEmail(ctx, sql.NullString{String: email, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return user.ID != "", nil
}

func NewUserRepository() IUserRepository {
	return &userRepositoryType{
		sqlc: database.New(global.Mdbc),
	}
}
