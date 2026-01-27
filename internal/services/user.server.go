package services

import "github.com/kien14502/ecommerce-be/internal/repo"

type IUserService interface {
	GetUserName(userID string) string
	Register(username, password string) error
}

type userService struct {
	userRepo repo.IUserRepository
}

// GetUserName implements [IUserService].
func (u *userService) GetUserName(userID string) string {
	return u.userRepo.FindOne(userID)
}

// Register implements [IUserService].
func (u *userService) Register(username string, password string) error {
	panic("unimplemented")
}

func NewUserService(userRepo repo.IUserRepository) IUserService {
	return &userService{
		userRepo: userRepo,
	}
}
