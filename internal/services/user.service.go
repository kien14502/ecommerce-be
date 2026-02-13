package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/repo"
	"github.com/kien14502/ecommerce-be/pkg/response"
	"github.com/kien14502/ecommerce-be/pkg/utils"
	"github.com/kien14502/ecommerce-be/pkg/utils/crypto"
	"github.com/segmentio/kafka-go"
)

type IUserService interface {
	GetUserName(userID string) string
	Register(email, password string) int
}

type userService struct {
	userRepo     repo.IUserRepository
	userAuthRepo repo.IUserAuthRepository
}

// GetUserName implements [IUserService].
func (u *userService) GetUserName(userID string) string {
	panic("unimplemented")

}

// Register implements [IUserService].
func (us *userService) Register(email string, password string) int {
	// Hash email
	hashEmail := crypto.GetHash(email)

	// check otp is available
	// user spam ...
	// check email exist
	isExistUser, err := us.userRepo.IsUserExisted(email)
	if err != nil {
		fmt.Print("[Register]", err.Error())
		return response.ErrCodeInternalServerError
	}
	if isExistUser {
		return response.ErrUserExisted
	}
	// new otp
	otp := utils.GenerateSixDigitOtp()
	// save otp in redis
	err = us.userAuthRepo.AddOTP(hashEmail, otp, int64(10*time.Minute))
	if err != nil {
		return response.ErrInvalidOTP
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
	err = global.Kafka.WriteMessages(context.Background(), message)
	if err != nil {
		fmt.Print("err[register]", err.Error())
		return response.ErrSendEmailOTP
	}
	// send email otp
	// err = sendto.SendTemplateEmailOTP([]string{email}, "phankien.epu@gmail.com", "verify-email.html", map[string]interface{}{
	// 	"otp":        strconv.Itoa(otp),
	// 	"name":       email,
	// 	"expMinutes": 10,
	// })
	// if err != nil {
	// 	return response.ErrSendEmailOTP
	// }
	return response.ErrCodeSuccess
}

func NewUserService(userRepo repo.IUserRepository, userAuthRepo repo.IUserAuthRepository) IUserService {
	return &userService{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
	}
}
