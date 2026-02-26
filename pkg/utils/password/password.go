package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong  = errors.New("password must be less than 72 characters")
	ErrHashFailed       = errors.New("failed to hash password")
	ErrInvalidPassword  = errors.New("invalid password")
)

const (
	// Cost của bcrypt (10-12 recommended, higher = more secure but slower)
	DefaultCost = 12
	MinCost     = 10
	MaxCost     = 14

	MinPasswordLength = 8
	MaxPasswordLength = 72 // bcrypt limit
)

// HashPassword mã hóa password thành bcrypt hash
func HashPassword(password string) (string, error) {
	// Validate password length
	if len(password) < MinPasswordLength {
		return "", ErrPasswordTooShort
	}
	if len(password) > MaxPasswordLength {
		return "", ErrPasswordTooLong
	}

	// Generate bcrypt hash
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrHashFailed, err)
	}

	return string(hashedBytes), nil
}

// ComparePassword so sánh plain password với hashed password
func ComparePassword(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}
		return fmt.Errorf("password comparison failed: %w", err)
	}
	return nil
}

// IsPasswordValid kiểm tra password có hợp lệ không (trước khi hash)
func IsPasswordValid(password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}
	if len(password) > MaxPasswordLength {
		return ErrPasswordTooLong
	}
	return nil
}

// HashPasswordWithCost hash password với custom cost
func HashPasswordWithCost(password string, cost int) (string, error) {
	if cost < MinCost || cost > MaxCost {
		return "", fmt.Errorf("cost must be between %d and %d", MinCost, MaxCost)
	}

	if err := IsPasswordValid(password); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrHashFailed, err)
	}

	return string(hashedBytes), nil
}
