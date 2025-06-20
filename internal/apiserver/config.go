package apiserver

import "os"

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	LogLevel string
	Port     string
}

func NewConfig() Config {
	return Config{
		Server: ServerConfig{
			LogLevel: os.Getenv("LOG_LEVEL"),
			Port:     os.Getenv("PORT"),
		},
	}
}
