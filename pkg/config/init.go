package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
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

func init() {

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

	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件发生变化:", e.Name)
		if err := viper.Unmarshal(&Global); err != nil {
			panic(err)
		}
	})
}
