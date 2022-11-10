package loggers

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger() *Logger {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &Logger{
		log,
	}
}

func (l *Logger) LogErr(err error, msg string) {
	l.logger.Error().Err(err).Msg(msg)
}

func (l *Logger) LogInfo(key, value, msg string) {
	l.logger.Info().Str(key, value).Msg(msg)
}

func (l *Logger) LogDebug(key, value, msg string) {
	l.logger.Debug().Str(key, value).Msg(msg)
}
