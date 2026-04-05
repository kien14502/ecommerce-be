package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/database"
	"github.com/kien14502/ecommerce-be/pkg/response"
)

type IUserSessionRepository interface {
	CreateUserSession(ctx context.Context, body database.CreateSessionParams) error
	GetUserSessionByToken(ctx context.Context, refreshTokenHash string) (*database.UserSession, error)
	DeleteUserSession(ctx context.Context, body database.DeleteSessionParams) error
	DeleteAllByUserID(ctx context.Context, userID string) error
}

type userSessionRepository struct {
	sqlc *database.Queries
}

// DeleteAllByUserID implements [IUserSessionRepository].
func (u *userSessionRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	err := u.sqlc.DeleteAllDevicesByUserID(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserSession implements [IUserSessionRepository].
func (u *userSessionRepository) DeleteUserSession(ctx context.Context, body database.DeleteSessionParams) error {
	if err := u.sqlc.DeleteSession(ctx, body); err != nil {
		return err
	}
	return nil
}

// GetUserSessionByToken implements [IUserSessionRepository].
func (u *userSessionRepository) GetUserSessionByToken(
	ctx context.Context,
	refreshTokenHash string,
) (*database.UserSession, error) {
	session, err := u.sqlc.GetSessionByToken(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrUnauthorized
		}
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, response.ErrUnauthorized
	}

	return &session, nil
}

// CreateUserSession implements [IUserSessionRepository].
func (u *userSessionRepository) CreateUserSession(ctx context.Context, body database.CreateSessionParams) error {
	if err := u.sqlc.CreateSession(ctx, body); err != nil {
		return err
	}
	return nil
}

func NewUserSessionRepository() IUserSessionRepository {
	return &userSessionRepository{
		sqlc: database.New(global.Mdbc),
	}
}
