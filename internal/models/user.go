package models

type UserModel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewUserModel(id, name string) *UserModel {
	return &UserModel{
		ID:   id,
		Name: name,
	}
}

