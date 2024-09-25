package config

type Mail struct {
	From     string `yaml:"from"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Address  string `yaml:"address"`
}
