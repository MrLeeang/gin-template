package db

import (
	"fmt"
	"gin-template/models"
	"gin-template/pkg/config"
	"gin-template/pkg/logger"
	"gin-template/pkg/utils"
	"strings"

	gormlogger "gorm.io/gorm/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func initDb() (*gorm.DB, error) {
	dbString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Global.Mysql.UserName,
		config.Global.Mysql.Password,
		config.Global.Mysql.Host,
		config.Global.Mysql.Port,
		config.Global.Mysql.DBName,
	)

	// 日志级别
	logLevel := gormlogger.Error

	if config.Global.Server.Debug {
		logLevel = gormlogger.Info
	}

	zapLogger := NewZapLogger(logger.Logger)
	zapLogger.SetAsDefault() // 可选：将 zapgorm2 设置为 GORM 的默认日志记录器
	zapLogger.LogMode(logLevel)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                      dbString,
		DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名禁用复数
		},
		Logger: zapLogger,
	})

	if err != nil {
		return db, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDb.SetMaxIdleConns(100) //设置最大连接数
	sqlDb.SetMaxOpenConns(20)  //设置最大的空闲连接数

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
			config.Global.Mysql.UserName,
			config.Global.Mysql.Password,
			config.Global.Mysql.Host,
			config.Global.Mysql.Port,
		)

		rootDb, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                      rootstring,
			DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		}), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				// TablePrefix:   "bc_", // 指定表前缀，修改默认表名
				SingularTable: true, // 表面禁用复数
			},
		})

		if err != nil {
			panic(err)
		}

		err = rootDb.Exec("create database " + config.Global.Mysql.DBName).Error

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
