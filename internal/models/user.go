package models

import (
	"github.com/kien14502/ecommerce-be/pkg/constants"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	UserName string `gorm:"type:varchar(255);not null" json:"user_name"`
	IsActive bool   `gorm:"type:boolean;default:true;not null" json:"is_active"`
	Roles    []Role `gorm:"many2many:user_role_db;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"roles"`
}

func (u *UserModel) TableName() string {
	return constants.UserTableName
}
