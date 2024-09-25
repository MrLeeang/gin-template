package config

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}
