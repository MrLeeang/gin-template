package models

import "github.com/jinzhu/gorm"

type Role struct {
	gorm.Model

	Uuid        string `json:"uuid" form:"uuid" gorm:"type:varchar(64);comment:'uuid'"`
	Name        string `json:"name" form:"name" gorm:"type:varchar(64);comment:'角色名称'"`
	DisplayName string `json:"display_name" form:"display_name" gorm:"type:varchar(64);comment:'中文名称，显示名称'"`
}

func (Role) TableName() string {
	return "role"
}
