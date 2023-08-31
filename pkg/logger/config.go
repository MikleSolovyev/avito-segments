package logger

type Config struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"debug"`
}
