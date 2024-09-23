package db

import (
	"fmt"
	"gin-template/models"
	"gin-template/pkg/config"
	"gin-template/pkg/logger"
	"gin-template/pkg/utils"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func initDb() (*gorm.DB, error) {
	dbString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Config.Mysql.UserName,
		config.Config.Mysql.Password,
		config.Config.Mysql.Host,
		config.Config.Mysql.Port,
		config.Config.Mysql.DBName,
	)

	db, err := gorm.Open("mysql", dbString)

	if err != nil {
		return db, err
	}

	sqlDb := db.DB()

	sqlDb.SetMaxIdleConns(100) //设置最大连接数
	sqlDb.SetMaxOpenConns(20)  //设置最大的空闲连接数

	// 设置全局表名禁用复数
	db.SingularTable(true)

	if config.Config.Debug {
		db.LogMode(true)
		db.SetLogger(logger.Logger)
	}

	db.AutoMigrate(new(models.User))
	db.AutoMigrate(new(models.Role))
	db.AutoMigrate(new(models.User2Role))
	db.AutoMigrate(new(models.UserLoginLog))

	return db, nil
}

func InitData() {

	adminRoleUuid := utils.GetUuid()
	userRoleUuid := utils.GetUuid()

	DB.Create(&models.Role{
		Uuid:        adminRoleUuid,
		Name:        "admin",
		DisplayName: "管理员",
	})

	DB.Create(&models.Role{
		Uuid:        userRoleUuid,
		Name:        "user",
		DisplayName: "用户",
	})

	adminUser := models.User{
		Uuid:         utils.GetUuid(),
		Username:     "admin",
		Name:         "管理员",
		Password:     "e10adc3949ba59abbe56e057f20f883e",
		Introduction: "平台管理员",
	}

	DB.Create(&adminUser)

	DB.Create(&models.User2Role{
		UserUuid: adminUser.Uuid,
		RoleUuid: adminRoleUuid,
	})

	DB.Create(&models.User2Role{
		UserUuid: adminUser.Uuid,
		RoleUuid: userRoleUuid,
	})

}

func init() {

	var err error

	DB, err = initDb()

	if err != nil {

		if !strings.Contains(err.Error(), "Unknown database") {
			panic(err)
		}

		// 创建数据库
		rootstring := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/?charset=utf8",
			config.Config.Mysql.UserName,
			config.Config.Mysql.Password,
			config.Config.Mysql.Host,
			config.Config.Mysql.Port,
		)

		rootDb, err := gorm.Open("mysql", rootstring)

		if err != nil {
			panic(err)
		}

		err = rootDb.Exec("create database " + config.Config.Mysql.DBName).Error

		if err != nil {
			panic(err)
		}

		DB, err = initDb()

		if err != nil {
			panic(err)
		}

		// 初始化数据
		InitData()

	}

}

func Add(model interface{}) error {
	return DB.Create(model).Error
}

func Delete(model interface{}, where ...interface{}) error {
	return DB.Delete(model, where...).Error
}

func Unscoped(model interface{}, where ...interface{}) error {
	return DB.Unscoped().Delete(model, where...).Error
}

func Save(model interface{}) error {
	return DB.Save(model).Error
}
