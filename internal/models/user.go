package models

import (
	"gorm.io/gorm"
)

const UserTableName = "user_db"

type User struct {
	gorm.Model
	ID       uint        `gorm:"primaryKey;autoIncrement"`
	UserName string      `gorm:"type:varchar(255);not null" json:"user_name"`
	IsActive bool        `gorm:"type:boolean;default:true;not null" json:"is_active"`
	Roles    []RoleModel `gorm:"many2many:user_role_db;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"roles"`
	Password string      `gorm:"type:varchar(255);not null" json:"-"`
}

func (u *User) TableName() string {
	return UserTableName
}
