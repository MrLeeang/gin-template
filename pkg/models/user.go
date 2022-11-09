package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Uuid         string `json:"uuid" form:"uuid" gorm:"type:varchar(64);comment:'uuid'"`
	Name         string `json:"name" form:"name" gorm:"type:varchar(64);comment:'用户名'"`
	Username     string `json:"username" form:"username" gorm:"type:varchar(64);comment:'登录名'"`
	Password     string `json:"password" form:"password" gorm:"type:varchar(64);comment:'密码'"`
	Introduction string `json:"introduction" form:"introduction" gorm:"type:text;comment:'个人介绍'"`
	Avatar       string `json:"avatar" form:"avatar" gorm:"type:varchar(255);comment:'头像'"`
	Roles        []Role `json:"roles" form:"roles" gorm:"-"`
}

func (User) TableName() string {
	return "user"
}

type User2Role struct {
	gorm.Model
	UserUuid string `json:"user_uuid" form:"user_uuid" gorm:"type:varchar(64);comment:'用户uuid'"`
	RoleUuid string `json:"role_uuid" form:"role_uuid" gorm:"type:varchar(64);comment:'角色uuid'"`
}

func (User2Role) TableName() string {
	return "user_2_role"
}

type UserLoginLog struct {
	gorm.Model
	UserUuid string `json:"user_uuid" form:"user_uuid" gorm:"type:varchar(64);comment:'用户uuid'"`
	OptType  string `json:"opt_type" form:"opt_type" gorm:"type:varchar(64);comment:'操作类型(login、logout)'"`
}

func (UserLoginLog) TableName() string {
	return "user_login_log"
}
