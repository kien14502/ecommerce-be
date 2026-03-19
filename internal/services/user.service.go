package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetUserName(userID string) string
	Register(ctx context.Context, in dto.RegisterRequest) error
	Login(ctx context.Context, body dto.LoginRequest, ip, userAgent string) (*dto.LoginResponse, error)
	VerifyOTP(ctx context.Context, body dto.VerifyOtpRequest) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error)
	// confirm password
	// logout
}

type userService struct {
	userRepo       repo.IUserRepository
	userVerifyRepo repo.IUserVerifyRepository
	redisService   IRedisService
	jwtService     IJwtService
	userDevice     repo.IUserDevicesRepository
}

// RefreshToken implements [IUserService].
func (u *userService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	// Parse refresh token (JWT)
	claims, err := u.jwtService.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, response.ErrInvalidEmail
	}
	// Get hash from Redis
	storedHash, err := u.redisService.GetRefreshToken(ctx, claims.UserID, claims.DeviceID)
	if err != nil {
		return nil, errors.New("session expired, please login again")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(refreshToken)); err != nil {
		// Possible token reuse attack
		u.redisService.DeleteAllRefreshTokens(ctx, claims.UserID)
		return nil, errors.New("invalid refresh token, all sessions revoked")
	}
	u.redisService.DeleteRefreshToken(ctx, claims.UserID, claims.DeviceID)
	// Extract userId + deviceId
	// Get hash from Redis
	// Compare hash (bcrypt)
	// Delete old token (rotate)
	// Generate new Access Token - Refresh Token
	// Save new hash to Redis
	return u.generateAndSaveTokens(ctx, claims.UserID, claims.DeviceID)
}

// VerifyOTP implements [IUserService].
func (u *userService) VerifyOTP(ctx context.Context, body dto.VerifyOtpRequest) (*dto.LoginResponse, error) {
	storedHash, err := u.redisService.GetOtp(ctx, body.Email)
	if err != nil {
		return nil, response.ErrOtpExpiredOrNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(body.Otp)); err != nil {
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
	userExisted, err := u.userRepo.GetUserByUsername(ctx, body.Username)
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
	deviceName, deviceType := utils.ParseDevice(userAgent)

	u.userDevice.CreateUserDevice(ctx, database.CreateDeviceParams{
		ID:         body.DeviceID,
		UserID:     userExisted.ID,
		DeviceName: sql.NullString{String: deviceName, Valid: true},
		DeviceType: sql.NullString{String: deviceType, Valid: true},
		UserAgent:  sql.NullString{String: userAgent, Valid: true},
		IpAddress:  sql.NullString{String: ip, Valid: true},
	})
	// Store refresh token in redis ttl
	refreshTokenHashed := crypto.GetHash(refreshToken)
	u.redisService.SaveRefreshToken(ctx, userExisted.ID, body.DeviceID, refreshTokenHashed)
	// Save device info to DB
	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GetUserName implements [IUserService].
func (u *userService) GetUserName(userID string) string {
	panic("unimplemented")
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

	if err := tx.Commit(); err != nil {
		return response.ErrInternalServer
	}

	// redis otp
	us.redisService.SaveOtp(ctx, in.Email, otpHash)

	// kafka email
	body := map[string]interface{}{
		"email": in.Email,
		"otp":   otpCode,
	}

	bodyBytes, _ := json.Marshal(body)

	msg := &sarama.ProducerMessage{
		Topic: consts.TopicOTP,
		Value: sarama.ByteEncoder(bodyBytes),
	}

	_, _, err = global.KafkaProducer.SendMessage(msg)
	if err != nil {
		global.Logger.Error("send kafka failed: " + err.Error())
	}

	return nil
}

func NewUserService(userRepo repo.IUserRepository, redisService IRedisService, userVerifyRepo repo.IUserVerifyRepository) IUserService {
	return &userService{
		userRepo:       userRepo,
		redisService:   redisService,
		userVerifyRepo: userVerifyRepo,
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
