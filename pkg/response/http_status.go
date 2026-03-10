package response

// Err code
const (
	// Common
	ErrInternalServer = "system"
	ErrInvalidParam   = "system"

	// Authentication code
	ErrOTPExisted     = "au0001"
	ErrInvalidEmail   = "au0002"
	ErrUserExisted    = "au0003"
	ErrRegisterFailed = "au0004"

	// Kafka
	ErrSendTopicFailed = "kaf0001"
)

// Success code
const (
	RegisterSuccess = "Register successful. Check email to get verify code!!!"
)

var msg = map[string]string{
	ErrOTPExisted:      "Otp already exists",
	ErrInvalidEmail:    "Invalid OTP",
	ErrUserExisted:     "User existed",
	ErrSendTopicFailed: "Failed to send kafka message",
	ErrInvalidParam:    "InvalidParam",
	ErrRegisterFailed:  "Register failed",
}

func GetMessage(code string) string {
	if message, exists := msg[code]; exists {
		return message
	}
	return "Unknown Error"
}
