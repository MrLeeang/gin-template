package config

type Server struct {
	ServerPort string `yaml:"serverPort"`
	UploadDir  string `yaml:"uploadDir"`
	MaxRequest int64  `yaml:"maxRequest"`
	Debug      bool   `yaml:"debug"`
	Encrypt    bool   `yaml:"encrypt"`
}
