package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtService struct {
	accessSecret  string
	refreshSecret string
}

// GenerateAccessToken implements [IJwtService].
func (j *jwtService) GenerateAccessToken(userID string, deviceID string) (string, string, error) {
	jti := uuid.New().String()
	claims := Claims{
		UserID:   userID,
		DeviceID: deviceID,
		JTI:      jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.accessSecret))
	return signed, jti, err
}

// GenerateRefreshToken implements [IJwtService].
func (j *jwtService) GenerateRefreshToken(userID string, deviceID string) (string, string, error) {
	jti := uuid.New().String()
	claims := Claims{
		UserID:   userID,
		DeviceID: deviceID,
		JTI:      jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.refreshSecret))
	return signed, jti, err
}

// ParseAccessToken implements [IJwtService].
func (j *jwtService) ParseAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*Claims), nil
}

// ParseRefreshToken implements [IJwtService].
func (j *jwtService) ParseRefreshToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.refreshSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*Claims), nil
}

type Claims struct {
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
	JTI      string `json:"jti"`
	jwt.RegisteredClaims
}

type IJwtService interface {
	GenerateAccessToken(userID, deviceID string) (string, string, error)
	GenerateRefreshToken(userID, deviceID string) (string, string, error)
	ParseAccessToken(tokenStr string) (*Claims, error)
	ParseRefreshToken(tokenStr string) (*Claims, error)
}

func NewJwtService(accessSecret, refreshSecret string) IJwtService {
	return &jwtService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}
