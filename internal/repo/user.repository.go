package repo

import (
	"errors"
	"strconv"
	"sync"

	"github.com/kien14502/ecommerce-be/internal/models"
)

type IUserRepository interface {
	FindOne(userID string) (*models.User, error)
	FindAll() []*models.User
	GetUserByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}

type userRepositoryType struct {
	users map[string]*models.User
	mu    sync.RWMutex
}

// Create implements [IUserRepository].
func (u *userRepositoryType) Create(user *models.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()
	userID := strconv.FormatUint(uint64(user.ID), 10)
	if _, exists := u.users[userID]; exists {
		return errors.New("user already exists")
	}
	u.users[userID] = user
	return nil
}

// FindAll implements [IUserRepository].
func (u *userRepositoryType) FindAll() []*models.User {
	u.mu.RLock()
	defer u.mu.RUnlock()

	users := make([]*models.User, 0, len(u.users))
	for _, user := range u.users {
		users = append(users, user)
	}
	return users
}

// FindOne implements [IUserRepository].
func (u *userRepositoryType) FindOne(userID string) (*models.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	user, exists := u.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByEmail implements [IUserRepository].
func (u *userRepositoryType) GetUserByEmail(email string) (*models.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	for _, user := range u.users {
		if user.UserName == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func NewUserRepository() IUserRepository {
	return &userRepositoryType{
		users: make(map[string]*models.User),
	}
}
