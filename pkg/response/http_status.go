package response

import "net/http"

// Success code
const (
	RegisterSuccess = "Register successful. Check email to get verify code!!!"
)

var (

	// ================= COMMON =================

	ErrInternalServer = &AppError{
		Code:       "system",
		Message:    "Internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrInvalidParam = &AppError{
		Code:       "system",
		Message:    "Invalid parameters",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrTooManyRequests = &AppError{
		HTTPStatus: 429,
		Code:       "SYS0002",
		Message:    "Too many requests, please try again later",
	}
	ErrNotFound = &AppError{
		HTTPStatus: 404,
		Code:       "SYS0404",
		Message:    "Resource not found",
	}
	// ================= AUTH =================

	ErrAlreadyVerified = &AppError{
		HTTPStatus: 400,
		Code:       "AUTH0002",
		Message:    "Email already verified",
	}
	ErrOtpStillValid = &AppError{
		HTTPStatus: 400,
		Code:       "AUTH0003",
		Message:    "OTP is still valid. Please check your email inbox",
	}
	ErrOtpExpired = &AppError{
		HTTPStatus: 400,
		Code:       "AUTH0004",
		Message:    "OTP has expired. Please request a new one",
	}

	ErrOTPExisted = &AppError{
		Code:       "au0001",
		Message:    "OTP already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrInvalidEmail = &AppError{
		Code:       "au0002",
		Message:    "Invalid email or OTP",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidOTP = &AppError{
		Code:       "au0008",
		Message:    "Invalid OTP",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrInvalidRefreshToken = &AppError{
		Code:       "au0009",
		Message:    "Invalid refresh token",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrVerifyFailed = &AppError{
		Code:       "au0004",
		Message:    "Verify failed",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrInvalidPassword = &AppError{
		Code:       "AU_PW0001",
		Message:    "Invalid password",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrUserExisted = &AppError{
		Code:       "au0003",
		Message:    "User already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrUnauthorized = &AppError{
		Code:       "au0004",
		Message:    "Unauthorized",
		HTTPStatus: http.StatusUnauthorized,
	}

	ErrRegisterFailed = &AppError{
		Code:       "au0004",
		Message:    "Register failed",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrEmailNotVerified = &AppError{
		Code:       "AU_E001",
		Message:    "Email not verified",
		HTTPStatus: http.StatusBadRequest,
	}

	ErrUsernameIsExisted = &AppError{
		Code:       "au0005",
		Message:    "Username already exists",
		HTTPStatus: http.StatusConflict,
	}

	ErrCreateUserFailed = &AppError{
		Code:       "au0006",
		Message:    "Create user failed",
		HTTPStatus: http.StatusInternalServerError,
	}

	ErrUserNotFound = &AppError{
		Code:       "auu0004",
		Message:    "User not existed",
		HTTPStatus: http.StatusBadRequest,
	}

	// OTP expired or not found
	ErrOtpExpiredOrNotFound = &AppError{
		Code:       "au0007",
		Message:    "OTP expired or not found",
		HTTPStatus: http.StatusConflict,
	}
	// ================= KAFKA =================

	ErrSendTopicFailed = &AppError{
		Code:       "kaf0001",
		Message:    "Failed to send kafka message",
		HTTPStatus: http.StatusInternalServerError,
	}
)
