package httpserver

import "time"

type Config struct {
	Address         string        `yaml:"address" env:"HTTP_ADDR" env-default:"localhost:8080"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"5s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"3s"`
}
