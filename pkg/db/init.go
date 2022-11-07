package db

import (
	"fmt"
	"gin-template/pkg/config"
	"gin-template/pkg/logger"
	"gin-template/pkg/models"
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

	// 指定表前缀，修改默认表名
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "gintemplate_" + defaultTableName
	}

	db.AutoMigrate(new(models.User))
	db.AutoMigrate(new(models.Role))

	return db, nil
}

func InitData() {

	adminRoleUuid := utils.GetUuid()

	DB.Create(&models.Role{
		Uuid:        adminRoleUuid,
		Name:        "admin",
		DisplayName: "管理员",
	})

	DB.Create(&models.Role{
		Uuid:        utils.GetUuid(),
		Name:        "user",
		DisplayName: "用户",
	})

	adminUser := models.User{
		Uuid:        utils.GetUuid(),
		Name:        "admin",
		DisplayName: "管理员",
		RoleUuid:    adminRoleUuid,
		Password:    "e10adc3949ba59abbe56e057f20f883e",
	}

	DB.Create(&adminUser)

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
