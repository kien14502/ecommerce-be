package response

const (
	ErrCodeBadRequest          = 400
	ErrCodeUnauthorized        = 401
	ErrCodeForbidden           = 403
	ErrCodeNotFound            = 404
	ErrCodeInternalServerError = 500
	ErrInvalidToken            = 498
)

var msg = map[int]string{
	ErrCodeBadRequest:          "Bad Request",
	ErrCodeUnauthorized:        "Unauthorized",
	ErrCodeForbidden:           "Forbidden",
	ErrCodeNotFound:            "Not Found",
	ErrCodeInternalServerError: "Internal Server Error",
	ErrInvalidToken:            "Invalid Token",
}

func GetMessage(code int) string {
	if message, exists := msg[code]; exists {
		return message
	}
	return "Unknown Error"
}
