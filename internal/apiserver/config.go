package apiserver

import "os"

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	LogLevel string
}

func NewConfig() Config {
	return Config{
		Server: ServerConfig{
			LogLevel: os.Getenv("LOG_LEVEL"),
		},
	}
}
