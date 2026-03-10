package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/kien14502/ecommerce-be/consts"
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/database"
	"github.com/kien14502/ecommerce-be/internal/dto"
	"github.com/kien14502/ecommerce-be/internal/models"
	"github.com/kien14502/ecommerce-be/internal/repo"
	"github.com/kien14502/ecommerce-be/pkg/otp"
	"github.com/kien14502/ecommerce-be/pkg/response"
	"github.com/kien14502/ecommerce-be/pkg/utils"
	"github.com/kien14502/ecommerce-be/pkg/utils/crypto"
)

type IUserService interface {
	GetUserName(userID string) string
	Register(ctx context.Context, in *models.RegisterInput) (resCode string, err error)
	Login(ctx context.Context, body dto.LoginRequest) (dto.LoginResponse, error)
	// verify opt
	VerifyOTP(ctx context.Context, body dto.LoginRequest) error
	// confirm password
	// logout
}

type userService struct {
	userRepo       repo.IUserRepository
	userAuthRepo   repo.IUserAuthRepository
	userVerifyRepo repo.IUserVerifyRepository
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
func (us *userService) Register(ctx context.Context, in *models.RegisterInput) (resCode string, err error) {
	// check email exist
	isExistUser, err := us.userRepo.IsUserExisted(ctx, in.Email)
	if err != nil {
		fmt.Print("[Register]", err.Error())
		return response.ErrInternalServer, errors.New("Internal server")
	}
	if isExistUser {
		return response.ErrUserExisted, errors.New("User existed")
	}

	// Hash email
	hashEmail := crypto.GetHash(in.Email)
	resCode, err = utils.OtpHandler(ctx, hashEmail)
	if err != nil {
		return resCode, errors.New("To many request")
	}

	// check otp is available
	// otpExisted, _ := us.userAuthRepo.GetOTP(ctx, hashKey)
	// if otpExisted != nil {
	// 	return response.ErrTooManyRequest, errors.New("please wait before requesting a new OTP")
	// }
	// user spam ...

	// new otp
	otp, err := otp.GenerateOTP(ctx, hashEmail)
	if err != nil {
		return response.ErrInternalServer, errors.New("Internal server")
	}

	body := make(map[string]interface{})
	body["otp"] = otp
	body["email"] = in.Email
	bodyRequest, _ := json.Marshal(body)

	msg := &sarama.ProducerMessage{
		Topic: consts.TopicOTP,
		Value: sarama.ByteEncoder(bodyRequest),
	}

	_, _, err = global.KafkaProducer.SendMessage(msg)
	if err != nil {
		global.Logger.Error("failed to send kafka message: " + err.Error())
		return response.ErrSendTopicFailed, err
	}

	err = us.userVerifyRepo.InsertOTPVerify(ctx, database.CreateOTPParams{})

	if err != nil {
		return response.ErrRegisterFailed, err
	}

	return response.RegisterSuccess, nil
}

func NewUserService(userRepo repo.IUserRepository, userAuthRepo repo.IUserAuthRepository) IUserService {
	return &userService{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
	}
}
