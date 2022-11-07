package config

type Server struct {
	ServerPort string
	UploadDir  string
	MaxRequest int64
	Debug      bool
	Encrypt    bool
}
