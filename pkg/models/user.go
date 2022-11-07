package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model

	Uuid        string `json:"uuid" form:"uuid" gorm:"type:varchar(64);comment:'uuid'"`
	Name        string `json:"name" form:"name" gorm:"type:varchar(64);comment:'用户名'"`
	DisplayName string `json:"display_name" form:"display_name" gorm:"type:varchar(64);comment:'姓名'"`
	Password    string `json:"password" form:"password" gorm:"type:varchar(64);comment:'密码'"`
	RoleUuid    string `json:"role_uuid" form:"role_uuid" gorm:"type:varchar(64);comment:'角色uuid'"`
	Role        Role   `json:"role" form:"role" gorm:"-"`
}
