package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/kien14502/ecommerce-be/consts"
	"github.com/kien14502/ecommerce-be/global"
	"github.com/kien14502/ecommerce-be/internal/database"
	"github.com/kien14502/ecommerce-be/internal/dto"
	"github.com/kien14502/ecommerce-be/internal/repo"
	"github.com/kien14502/ecommerce-be/pkg/otp"
	"github.com/kien14502/ecommerce-be/pkg/response"
	"github.com/kien14502/ecommerce-be/pkg/utils"
	"github.com/kien14502/ecommerce-be/pkg/utils/crypto"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// TODO implement all goroutine
type IUserService interface {
	Register(ctx context.Context, in dto.RegisterRequest) error
	Login(ctx context.Context, body dto.LoginRequest, ip, userAgent string) (*dto.LoginResponse, error)
	VerifyOTP(ctx context.Context, body dto.VerifyOtpRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error)
	GetMe(ctx context.Context, userID string) (*dto.UserResponse, error)
	ResendVerifyEmail(ctx context.Context, body dto.ResendVerifyOtpRequest) error
	// confirm password
	Logout(ctx context.Context, accessToken string) error
}

type userService struct {
	userRepo       repo.IUserRepository
	userVerifyRepo repo.IUserVerifyRepository
	redisService   IRedisService
	jwtService     IJwtService
	userDevice     repo.IUserDevicesRepository
	userSession    repo.IUserSessionRepository
}

// ResendVerifyEmail implements [IUserService].
func (u *userService) ResendVerifyEmail(ctx context.Context, body dto.ResendVerifyOtpRequest) error {
	userExisted, err := u.userRepo.GetUserByEmail(ctx, body.Email)
	if err != nil {
		return response.ErrUserNotFound
	}
	fmt.Println("User", userExisted)
	if userExisted.IsEmailVerified.Valid && userExisted.IsEmailVerified.Bool {
		return response.ErrAlreadyVerified
	}
	oldOtp, err := u.redisService.GetOtp(ctx, body.Email)
	if err != nil {
		return fmt.Errorf("Get otp failed:%w", err)
	}
	if oldOtp != "" {
		return response.ErrOtpStillValid
	}
	otpCode := otp.GenerateSixDigitOtp()
	otpHash := otp.HashOTP(otpCode)
	err = u.userVerifyRepo.InsertOTPVerify(ctx, database.CreateOTPParams{
		ID:      uuid.New().String(),
		Email:   body.Email,
		OtpHash: otpHash,
		Purpose: database.NullOtpVerificationsPurpose{
			OtpVerificationsPurpose: database.OtpVerificationsPurposeRegister,
			Valid:                   true,
		},
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	if err != nil {
		return fmt.Errorf("insert otp failed: %w", err)
	}

	err = u.redisService.SaveOtp(ctx, body.Email, otpHash)
	if err != nil {
		return fmt.Errorf("set otp redis failed: %w", err)
	}

	go u.sendOtpKafka(body.Email, otpCode)

	return nil
}

// Logout implements [IUserService].
func (u *userService) Logout(ctx context.Context, accessToken string) error {
	claims, err := u.jwtService.ParseAccessToken(accessToken)
	if err != nil {
		return response.ErrUnauthorized
	}
	err = u.userSession.DeleteUserSession(ctx, database.DeleteSessionParams{
		UserID:   claims.UserID,
		DeviceID: claims.DeviceID,
	})
	if err != nil {
		return response.ErrUnauthorized
	}
	if err := u.redisService.DeleteRefreshToken(ctx, claims.UserID, claims.DeviceID); err != nil {
		return response.ErrUnauthorized
	}
	return nil
}

// GetMe implements [IUserService].
func (u *userService) GetMe(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindOne(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email.String,
		Username:  user.Username.String,
		FullName:  user.FullName.String,
		AvatarUrl: &user.AvatarUrl.String,
	}, nil

}

// RefreshToken implements [IUserService].
func (u *userService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {

	// 1. Parse JWT
	claims, err := u.jwtService.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, response.ErrUnauthorized
	}

	// 2. Hash token (SHA256)
	hash := crypto.GetHash(refreshToken)

	// 3. Check DB (source of truth)
	_, err = u.userSession.GetUserSessionByToken(ctx, hash)
	if err != nil {
		return nil, response.ErrUnauthorized
	}

	// 4. (optional) check Redis
	storedHash, err := u.redisService.GetRefreshToken(ctx, claims.UserID, claims.DeviceID)
	if err != nil {
		return nil, response.ErrUnauthorized
	}

	// 5. Compare hash (simple string compare)
	if storedHash != hash {
		_ = u.userSession.DeleteAllByUserID(ctx, claims.UserID)
		_ = u.redisService.DeleteAllRefreshTokens(ctx, claims.UserID)
		return nil, errors.New("token reuse detected, all sessions revoked")
	}

	// 6. Rotate: delete old session
	err = u.userSession.DeleteUserSession(ctx, database.DeleteSessionParams{
		UserID:   claims.UserID,
		DeviceID: claims.DeviceID,
	})
	if err != nil {
		return nil, err
	}

	// 7. Generate new tokens
	loginResponse, err := u.generateAndSaveTokens(ctx, claims.UserID, claims.DeviceID)
	if err != nil {
		return nil, err
	}

	// 8. Save Redis
	newHash := crypto.GetHash(loginResponse.RefreshToken)
	if err := u.redisService.SaveRefreshToken(ctx, claims.UserID, claims.DeviceID, newHash); err != nil {
		return nil, err
	}

	return loginResponse, nil
}

// VerifyOTP implements [IUserService].
func (u *userService) VerifyOTP(ctx context.Context, body dto.VerifyOtpRequest) (*dto.LoginResponse, error) {
	storedHash, err := u.redisService.GetOtp(ctx, body.Email)
	fmt.Print("storedHash", storedHash)
	if err != nil {
		return nil, response.ErrOtpExpiredOrNotFound
	}
	if isOtpValid := otp.CompareOTPHashed(storedHash, body.Otp); isOtpValid {
		return nil, response.ErrInvalidOTP
	}
	if err := u.userRepo.MarkEmailVerified(ctx, body.Email); err != nil {
		return nil, response.ErrVerifyFailed
	}
	u.redisService.DeleteOtp(ctx, body.Email)

	user, err := u.userRepo.GetUserByEmail(ctx, body.Email)
	if err != nil {
		return nil, response.ErrUserNotFound
	}

	deviceID := uuid.New().String()
	// TODO implement create device & session
	accessToken, _, err := u.jwtService.GenerateAccessToken(user.ID, deviceID)
	if err != nil {
		return nil, err
	}

	refreshToken, _, err := u.jwtService.GenerateRefreshToken(user.ID, deviceID)
	if err != nil {
		return nil, err
	}
	refreshHash, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	u.redisService.SaveRefreshToken(ctx, user.ID, deviceID, string(refreshHash))

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login implements [IUserService].
func (u *userService) Login(ctx context.Context, body dto.LoginRequest, ip, userAgent string) (*dto.LoginResponse, error) {
	// Check existed user (verified, password) - Hash and compare password
	userExisted, err := u.userRepo.GetUserByEmail(ctx, body.Username)
	if err != nil {
		return nil, response.ErrUserNotFound
	}
	if !userExisted.IsEmailVerified.Bool {
		return nil, response.ErrEmailNotVerified
	}

	isValidPassword := crypto.ComparePassword(body.Password, userExisted.PasswordHash.String)
	if isValidPassword == false {
		return nil, response.ErrInvalidPassword
	}
	// Generate access token & refresh token
	accessToken, _, _ := u.jwtService.GenerateAccessToken(userExisted.ID, body.DeviceID)
	refreshToken, _, _ := u.jwtService.GenerateRefreshToken(userExisted.ID, body.DeviceID)
	// Get/Create device ID
	refreshTokenHashed := crypto.GetHash(refreshToken)
	u.redisService.SaveRefreshToken(ctx, userExisted.ID, body.DeviceID, refreshTokenHashed)
	deviceName, deviceType := utils.ParseDevice(userAgent)

	tx, err := global.Mdbc.BeginTx(ctx, nil)
	if err != nil {
		return nil, response.ErrInternalServer
	}
	queries := database.New(global.Mdbc)
	qtx := queries.WithTx(tx)
	defer tx.Rollback()

	err = qtx.CreateDevice(ctx, database.CreateDeviceParams{
		ID:         body.DeviceID,
		UserID:     userExisted.ID,
		DeviceName: sql.NullString{String: deviceName, Valid: true},
		DeviceType: sql.NullString{String: deviceType, Valid: true},
		UserAgent:  sql.NullString{String: userAgent, Valid: true},
		IpAddress:  sql.NullString{String: ip, Valid: true},
	})

	if err != nil {
		return nil, response.ErrInvalidParam
	}

	newUserSessionID := uuid.New().String()
	err = qtx.CreateSession(ctx, database.CreateSessionParams{
		ID:               newUserSessionID,
		UserID:           userExisted.ID,
		DeviceID:         body.DeviceID,
		RefreshTokenHash: refreshTokenHashed,
		ExpiresAt:        time.Now().Add(time.Duration(global.Config.Jwt.RefreshExp) * time.Hour),
	})
	if err != nil {
		return nil, &response.AppError{
			Code:       "uss0001",
			Message:    err.Error(),
			HTTPStatus: http.StatusBadRequest,
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, response.ErrInternalServer
	}
	// Save device info to DB
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Register implements [IUserService].
func (us *userService) Register(ctx context.Context, in dto.RegisterRequest) error {
	existed, err := us.userRepo.IsUserExisted(ctx, in.Email)
	if err != nil {
		return response.ErrInternalServer
	}
	if existed {
		return response.ErrUserExisted
	}

	hashPassword, _ := crypto.HashPassword(in.Password)
	userID := uuid.New().String()

	otpCode := otp.GenerateSixDigitOtp()
	otpHash := otp.HashOTP(otpCode)

	tx, err := global.Mdbc.BeginTx(ctx, nil)
	if err != nil {
		return response.ErrInternalServer
	}
	queries := database.New(global.Mdbc)
	qtx := queries.WithTx(tx)
	defer tx.Rollback()

	err = qtx.CreateUser(ctx, database.CreateUserParams{
		ID:           userID,
		Email:        sql.NullString{String: in.Email, Valid: true},
		PasswordHash: sql.NullString{String: hashPassword, Valid: true},
		Username:     sql.NullString{String: in.Username, Valid: true},
		FullName:     sql.NullString{String: in.FullName, Valid: true},
	})

	if err != nil {
		return response.ErrCreateUserFailed
	}

	err = qtx.CreateOTP(ctx, database.CreateOTPParams{
		ID:      uuid.New().String(),
		Email:   in.Email,
		OtpHash: otpHash,
		Purpose: database.NullOtpVerificationsPurpose{
			OtpVerificationsPurpose: database.OtpVerificationsPurposeRegister,
			Valid:                   true,
		},
		ExpiresAt: time.Now().Add(5 * time.Minute),
	})

	if err != nil {
		return response.ErrRegisterFailed
	}
	// redis otp
	err = us.redisService.SaveOtp(ctx, in.Email, otpHash)

	if err != nil {
		return response.ErrInternalServer
	}

	if err := tx.Commit(); err != nil {
		return response.ErrInternalServer
	}

	go us.sendOtpKafka(in.Email, otpCode)

	return nil
}

func NewUserService(userRepo repo.IUserRepository, redisService IRedisService, userVerifyRepo repo.IUserVerifyRepository, jwtService IJwtService, userDevices repo.IUserDevicesRepository, userSession repo.IUserSessionRepository) IUserService {
	return &userService{
		userRepo:       userRepo,
		redisService:   redisService,
		userVerifyRepo: userVerifyRepo,
		jwtService:     jwtService,
		userDevice:     userDevices,
		userSession:    userSession,
	}
}

// ─── Helper ───────────────────────────────────────────────

func (s *userService) generateAndSaveTokens(ctx context.Context, userID, deviceID string) (*dto.LoginResponse, error) {
	// Generate access token
	accessToken, _, err := s.jwtService.GenerateAccessToken(userID, deviceID)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, _, err := s.jwtService.GenerateRefreshToken(userID, deviceID)
	if err != nil {
		return nil, err
	}

	// Hash refresh token
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Save to Redis
	if err := s.redisService.SaveRefreshToken(ctx, userID, deviceID, string(hash)); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *userService) sendOtpKafka(email string, otpCode int) {
	payload := map[string]interface{}{
		"email": email,
		"otp":   otpCode,
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		global.Logger.Error("marshal kafka payload failed", zap.Error(err))
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: consts.TopicOTP,
		Value: sarama.ByteEncoder(bodyBytes),
	}
	_, _, err = global.KafkaProducer.SendMessage(msg)
	if err != nil {
		global.Logger.Error("send kafka failed", zap.Error(err))
	}
}
