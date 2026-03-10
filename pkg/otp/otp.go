package otp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"math/rand"

	"github.com/kien14502/ecommerce-be/global"
	"github.com/redis/go-redis/v9"
)

const (
	otpTTL      = 10 * time.Minute
	resendTTL   = 10 * time.Minute
	maxAttempts = 5
)

func hashOTP(otp int) string {
	str := strconv.Itoa(otp)
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}

func generateSixDigitOtp() int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := 100000 + rng.Intn(999999)
	return otp
}

func GenerateOTP(ctx context.Context, hashEmail string) (int, error) {
	otpKey := fmt.Sprintf("otp:user:%s", hashEmail)
	limitKey := fmt.Sprintf("otp:limit:%s", hashEmail)

	// resend cooldown
	exists, err := global.Rdb.Exists(ctx, limitKey).Result()
	if err != nil {
		return 0, err
	}

	if exists == 1 {
		return 0, fmt.Errorf("please wait before requesting another otp")
	}

	otp := generateSixDigitOtp()
	hashed := hashOTP(otp)

	err = global.Rdb.Set(ctx, otpKey, hashed, otpTTL).Err()
	if err != nil {
		return 0, err
	}

	// rate limit resend
	global.Rdb.Set(ctx, limitKey, "1", resendTTL)

	return otp, nil
}

func VerifyOTP(ctx context.Context, hashEmail string, input int) error {
	otpKey := fmt.Sprintf("otp:user:%s", hashEmail)
	attemptKey := fmt.Sprintf("otp:attempt:%s", hashEmail)

	val, err := global.Rdb.Get(ctx, otpKey).Result()

	if err == redis.Nil {
		return fmt.Errorf("otp expired")
	}

	if err != nil {
		return err
	}

	// check attempts
	attempts, _ := global.Rdb.Incr(ctx, attemptKey).Result()
	if attempts == 1 {
		global.Rdb.Expire(ctx, attemptKey, otpTTL)
	}

	if attempts > maxAttempts {
		return fmt.Errorf("too many attempts")
	}

	hashedInput := hashOTP(input)

	if hashedInput != val {
		return fmt.Errorf("invalid otp")
	}

	// success → delete keys
	global.Rdb.Del(ctx, otpKey)
	global.Rdb.Del(ctx, attemptKey)

	return nil
}
