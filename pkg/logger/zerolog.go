package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func NewZerolog(cfg *Config) (*zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("incorrect logger level: %s", cfg.Level)
	}

	logger := zerolog.Logger{}
	if level == 0 {
		logger = zerolog.
			New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
			With().
			Timestamp().
			Logger()
	} else {
		logger = zerolog.
			New(os.Stderr).
			Level(level).
			With().
			Timestamp().
			Logger()
	}

	return &logger, nil
}
