package models

type RegisterInput struct {
	Email    string `json:"verify_key"`
	Password string `json:"password"`
	Username string `json:"username"`
}
