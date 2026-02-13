package response

const (
	ErrCodeSuccess             = 2000
	ErrCodeBadRequest          = 400
	ErrCodeUnauthorized        = 401
	ErrCodeForbidden           = 403
	ErrCodeNotFound            = 404
	ErrCodeInternalServerError = 500
	ErrInvalidToken            = 498
	ErrUserExisted             = 1001
	ErrInvalidOTP              = 1002
	ErrSendEmailOTP            = 1003
)

var msg = map[int]string{
	ErrCodeBadRequest:          "Bad Request",
	ErrCodeUnauthorized:        "Unauthorized",
	ErrCodeForbidden:           "Forbidden",
	ErrCodeNotFound:            "Not Found",
	ErrCodeInternalServerError: "Internal Server Error",
	ErrInvalidToken:            "Invalid Token",
	ErrUserExisted:             "User existed",
	ErrInvalidOTP:              "Invalid OTP",
	ErrSendEmailOTP:            "Send email failed",
	ErrCodeSuccess:             "Success",
}

func GetMessage(code int) string {
	if message, exists := msg[code]; exists {
		return message
	}
	return "Unknown Error"
}
