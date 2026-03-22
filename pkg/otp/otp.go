package otp

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"

	"math/rand"
)

const (
	otpTTL      = 10 * time.Minute
	resendTTL   = 10 * time.Minute
	maxAttempts = 5
)

func HashOTP(otp int) string {
	str := strconv.Itoa(otp)
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

func GenerateSixDigitOtp() int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := 100000 + rng.Intn(999999)
	return otp
}

func CompareOTPHashed(input string, hashedOTP string) bool {
	parsed, err := strconv.Atoi(input)
	if err != nil {
		return false
	}
	return HashOTP(parsed) == hashedOTP
}
