package config

import (
	"gopkg.in/ini.v1"
)

type config struct {
	Mysql
	Server
	Consul
	Mail
	Service
	Alibaba
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

	// service
	Config.Service.Address = cfg.Section("service").Key("address").String()

	// consul
	Config.Consul.Address = cfg.Section("consul").Key("address").String()

	// mail
	Config.Mail.From = cfg.Section("mail").Key("from").String()
	Config.Mail.Username = cfg.Section("mail").Key("username").String()
	Config.Mail.Password = cfg.Section("mail").Key("password").String()
	Config.Mail.Host = cfg.Section("mail").Key("host").String()
	Config.Mail.Address = cfg.Section("mail").Key("address").String()

	// alibaba
	Config.Alibaba.AccessKeyId = cfg.Section("alibaba").Key("accessKeyId").String()
	Config.Alibaba.AccessKeySecret = cfg.Section("alibaba").Key("accessKeySecret").String()
	Config.Alibaba.SignName = cfg.Section("alibaba").Key("signName").String()
	Config.Alibaba.TemplateCode = cfg.Section("alibaba").Key("templateCode").String()

}
