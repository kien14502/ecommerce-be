package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/dto"
	"github.com/kien14502/ecommerce-be/internal/repo"
	"github.com/kien14502/ecommerce-be/pkg/response"
	"github.com/kien14502/ecommerce-be/pkg/utils"
	"github.com/kien14502/ecommerce-be/pkg/utils/crypto"
	"github.com/segmentio/kafka-go"
)

type IUserService interface {
	GetUserName(userID string) string
	Register(ctx context.Context, email, password string) (int, error)
	Login(ctx context.Context, body dto.LoginRequest) (dto.LoginResponse, error)
	// verify opt
	VerifyOTP(ctx context.Context, body dto.LoginRequest) error
	// confirm password
	// logout
}

type userService struct {
	userRepo     repo.IUserRepository
	userAuthRepo repo.IUserAuthRepository
}

// VerifyOTP implements [IUserService].
func (u *userService) VerifyOTP(ctx context.Context, body dto.LoginRequest) error {
	panic("unimplemented")
}

// Login implements [IUserService].
func (u *userService) Login(ctx context.Context, body dto.LoginRequest) (dto.LoginResponse, error) {
	panic("unimplemented")
	// Check existed user (verified, password) - Hash and compare password
	// Generate access token & refresh token
	// Save refresh token to cookie
	// Store refresh token in redis ttl
}

// GetUserName implements [IUserService].
func (u *userService) GetUserName(userID string) string {
	panic("unimplemented")

}

// Register implements [IUserService].
func (us *userService) Register(ctx context.Context, email string, password string) (int, error) {
	// Hash email
	hashEmail := crypto.GetHash(email)

	// check otp is available
	// user spam ...
	// check email exist
	isExistUser, err := us.userRepo.IsUserExisted(ctx, email)
	if err != nil {
		fmt.Print("[Register]", err.Error())
		return response.ErrInternalServer, errors.New("Internal server")
	}
	if isExistUser {
		return response.ErrBadRequest, errors.New("User existed")
	}
	// new otp
	otp := utils.GenerateSixDigitOtp()
	// save otp in redis
	err = us.userAuthRepo.AddOTP(ctx, hashEmail, otp, int64(10*time.Minute))
	if err != nil {
		return response.ErrInternalServer, errors.New("Internal server")
	}

	body := make(map[string]interface{})
	body["otp"] = otp
	body["email"] = email
	bodyRequest, _ := json.Marshal(body)

	message := kafka.Message{
		Key:   []byte("otp-auth"),
		Value: []byte(bodyRequest),
		Time:  time.Now(),
	}
	err = global.KafkaProducer.PublishUserRegistered(ctx, message)
	if err != nil {
		fmt.Print("err[register]", err.Error())
		return response.ErrInternalServer, errors.New("Internal server")
	}
	return response.SuccessCode, nil
}

func NewUserService(userRepo repo.IUserRepository, userAuthRepo repo.IUserAuthRepository) IUserService {
	return &userService{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
	}
}
