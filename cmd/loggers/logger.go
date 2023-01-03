package loggers

import (
	"os"

	"github.com/rs/zerolog"
)

// Logger структура логгера.
type Logger struct {
	logger zerolog.Logger
}

// NewLogger создание нового логгера.
func NewLogger() *Logger {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &Logger{
		log,
	}
}

// LogErr обработка уровень ошибок.
func (l *Logger) LogErr(err error, msg string) {
	l.logger.Error().Err(err).Msg(msg)
}

// LogInfo обработка уровень инфо.
func (l *Logger) LogInfo(key, value, msg string) {
	l.logger.Info().Str(key, value).Msg(msg)
}

// LogDebug обработка уровень дебаг.
func (l *Logger) LogDebug(key, value, msg string) {
	l.logger.Debug().Str(key, value).Msg(msg)
}
