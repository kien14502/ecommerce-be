package models

type RegisterInput struct {
	Email         string `json:"verify_key"`
	VerifyPurpose string `json:"verify_purpose"`
}
