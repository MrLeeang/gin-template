package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Mysql   Mysql   `yaml:"mysql"`
	Server  Server  `yaml:"server"`
	Consul  Consul  `yaml:"consul"`
	Mail    Mail    `yaml:"mail"`
	Service Service `yaml:"service"`
	Alibaba Alibaba `yaml:"alibaba"`
}

var Global = &Config{}

func InitializeConfig() {

	// 设置配置文件名和类型
	viper.SetConfigName("config") // 不需要文件扩展名
	viper.SetConfigType("yaml")   // 指定文件类型
	viper.AddConfigPath(".")      // 设置查找配置文件的路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// 将配置解码到结构体
	if err := viper.Unmarshal(&Global); err != nil {
		panic(err)
	}

}
