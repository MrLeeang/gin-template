package db

import (
	"context"
	"fmt"
	"gin-template/models"
	"gin-template/pkg/config"
	"gin-template/pkg/utils"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Session *gorm.DB

// 初始化数据库连接池
func createSession() error {
	dbString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Global.Mysql.UserName,
		config.Global.Mysql.Password,
		config.Global.Mysql.Host,
		config.Global.Mysql.Port,
		config.Global.Mysql.DBName,
	)

	// 日志级别
	logLevel := logger.Warn

	if config.Global.Debug {
		logLevel = logger.Info
	}

	var err error

	Session, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                      dbString,
		DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名禁用复数
		},
		Logger: zapLogger{
			level:                     logLevel,
			IgnoreRecordNotFoundError: false,           // 忽略not found
			SlowThreshold:             2 * time.Second, // 慢sql
		},
	})

	if err != nil {
		return err
	}

	sqlDb, err := Session.DB()
	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(20)   //设置最大空闲连接数
	sqlDb.SetMaxOpenConns(100)  //设置最大打开连接数
	sqlDb.SetConnMaxLifetime(0) // 设置连接的最大生命周期

	Session.AutoMigrate(new(models.User))
	Session.AutoMigrate(new(models.Role))
	Session.AutoMigrate(new(models.User2Role))
	Session.AutoMigrate(new(models.UserLoginLog))

	return nil
}

// 初始化数据
func initData() {

	adminRoleUuid := utils.GetUuid()
	userRoleUuid := utils.GetUuid()

	Session.Create(&models.Role{
		Uuid:        adminRoleUuid,
		Name:        "admin",
		DisplayName: "管理员",
	})

	Session.Create(&models.Role{
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

	Session.Create(&adminUser)

	Session.Create(&models.User2Role{
		UserUuid: adminUser.Uuid,
		RoleUuid: adminRoleUuid,
	})

	Session.Create(&models.User2Role{
		UserUuid: adminUser.Uuid,
		RoleUuid: userRoleUuid,
	})
}

// 创建数据库
func createDb() error {

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
	}))

	if err != nil {
		return err
	}

	err = rootDb.Exec("create database " + config.Global.Mysql.DBName).Error

	if err != nil {
		return err
	}

	return nil
}

func InitializeDatabase() {

	if Session != nil {
		// 关闭已经打开的数据库连接
		sqlDB, err := Session.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	if err := createSession(); err != nil {
		if !strings.Contains(err.Error(), "Unknown database") {
			panic(err)
		}

		// 数据库不存在,创建数据库
		if err := createDb(); err != nil {
			panic(err)
		}

		if err := createSession(); err != nil {
			panic(err)
		}

		// 初始化数据
		initData()
	}

}

func Add(ctx context.Context, model interface{}) error {
	return Session.WithContext(ctx).Create(model).Error
}

func Delete(ctx context.Context, model interface{}, where ...interface{}) error {
	return Session.WithContext(ctx).Delete(model, where...).Error
}

func Unscoped(ctx context.Context, model interface{}, where ...interface{}) error {
	return Session.WithContext(ctx).Unscoped().Delete(model, where...).Error
}

func First(ctx context.Context, model interface{}, where ...interface{}) error {
	return Session.WithContext(ctx).First(model, where...).Error
}

func Save(ctx context.Context, model interface{}) error {
	return Session.WithContext(ctx).Save(model).Error
}
