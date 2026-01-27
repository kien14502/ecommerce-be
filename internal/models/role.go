package models

import (
	"github.com/kien14502/ecommerce-be/pkg/constants"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	RoleName string `gorm:"type:varchar(255);not null;unique" json:"role_name"`
	RoleDesc string `gorm:"type:varchar(512);not null" json:"role_desc"`
}

func (r *Role) TableName() string {
	return constants.RoleTableName
}
