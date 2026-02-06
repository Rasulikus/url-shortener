package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Level string
}

func Init(cfg Config) error {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	zerolog.SetGlobalLevel(level)

	zerolog.TimeFieldFormat = time.RFC3339

	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	return nil
}
