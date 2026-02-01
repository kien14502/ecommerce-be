package models

import (
	"gorm.io/gorm"
)

const RoleTableName = "user_role_db"

type RoleModel struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	RoleName string `gorm:"type:varchar(255);not null;unique" json:"role_name"`
	RoleDesc string `gorm:"type:varchar(512);not null" json:"role_desc"`
}

func (r *RoleModel) TableName() string {
	return RoleTableName
}
