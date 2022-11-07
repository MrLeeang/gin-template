package config

import (
	"gopkg.in/ini.v1"
)

type config struct {
	Mysql
	Server
}

var Config = &config{}

func init() {

	cfg, err := ini.Load("config.ini")

	if err != nil {
		// 从服务器配置路径获取
		cfg, err = ini.Load("/etc/gin-template/config.ini")
		if err != nil {
			panic("Fail to read file: config.ini")
		}
	}

	// mysql
	Config.Mysql.UserName = cfg.Section("mysql").Key("username").String()
	Config.Mysql.Password = cfg.Section("mysql").Key("password").String()
	Config.Mysql.Host = cfg.Section("mysql").Key("host").String()
	Config.Mysql.Port = cfg.Section("mysql").Key("port").String()
	Config.Mysql.DBName = cfg.Section("mysql").Key("dbname").String()

	// server
	Config.Server.ServerPort = cfg.Section("server").Key("serverPort").String()
	Config.Server.UploadDir = cfg.Section("server").Key("uploadDir").String()
	Config.Server.MaxRequest, _ = cfg.Section("server").Key("maxRequest").Int64()
	Config.Server.Debug, _ = cfg.Section("server").Key("debug").Bool()
	Config.Server.Encrypt, _ = cfg.Section("server").Key("encrypt").Bool()

}
