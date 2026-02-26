package response

const (
	UnknownCode              = 500
	UnknownReason            = ""
	SupportPackageIsVersion1 = true
	SuccessCode              = 200
)

const (
	ErrBadRequest         = 400
	ErrUnauthorized       = 401
	ErrForbidden          = 403
	ErrNotFound           = 404
	ErrConflict           = 409
	ErrTooManyRequest     = 429
	ErrClientClosed       = 499
	ErrInternalServer     = 500
	ErrServiceUnavailable = 503
	ErrGatewayTimeOut     = 504
)

var msg = map[int]string{
	// HTTP errors
	ErrBadRequest:         "Bad Request",
	ErrUnauthorized:       "Unauthorized",
	ErrForbidden:          "Forbidden",
	ErrNotFound:           "Not Found",
	ErrConflict:           "Conflict",
	ErrTooManyRequest:     "Too Many Requests",
	ErrClientClosed:       "Client Closed Request",
	ErrInternalServer:     "Internal Server Error",
	ErrServiceUnavailable: "Service Unavailable",
	ErrGatewayTimeOut:     "Gateway Timeout",
}

func GetMessage(code int) string {
	if message, exists := msg[code]; exists {
		return message
	}
	return "Unknown Error"
}
