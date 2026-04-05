package dto

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
	Username string `json:"username" validate:"required,min=3,max=30,alphanum"`
	FullName string `json:"full_name" validate:"required,min=3,max=100,alphaunicode"`
}
type LoginRequest struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	DeviceID string `json:"device_id"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required,len=6,numeric"`
}

type UserResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	FullName  string  `json:"full_name"`
	AvatarUrl *string `json:"avatar_url,omitempty"`
}

type ResendVerifyOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
}
