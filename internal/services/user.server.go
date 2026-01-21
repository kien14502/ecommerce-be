package services

import "github.com/kien14502/ecommerce-be/internal/repo"

type UserServiceType struct{
	userRepo *repo.UserRepositoryType
}

func NewUserService() *UserServiceType {
	return &UserServiceType{
		userRepo: repo.NewUserRepository(),
	}
}

func (us *UserServiceType) GetUserName(userID string) string {
	return us.userRepo.GetUserByID(userID)
}