package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"helpdesk/backend/internal/config"
)

func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return zerolog.DebugLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

func Init(cfg config.Config) error {
	zerolog.SetGlobalLevel(parseLevel(cfg.LogLevel))
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if strings.EqualFold(cfg.GoEnv, "production") {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	} else {
		console := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05",
		}
		log.Logger = zerolog.New(console).With().Timestamp().Logger()
	}

	log.Info().Msg("logger initialized")
	return nil
}

func L() *zerolog.Logger {
	return &log.Logger
}

func Sync() error {
	return nil
}
