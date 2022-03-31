package logging

import (
	"github.com/rs/zerolog"
	"os"
)

const timeFormat = "2006-01-02 15:04:05"

var logger zerolog.Logger

func init() {
	zerolog.CallerSkipFrameCount = 3
	zerolog.TimeFieldFormat = timeFormat

	logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func Debug() *zerolog.Event {
	return logger.Debug().Caller()
}

func Info() *zerolog.Event {
	return logger.Info().Caller()
}

func Warn() *zerolog.Event {
	return logger.Warn().Caller()
}

func Error(err error) *zerolog.Event {
	return logger.Err(err).Caller()
}

func Fatal() *zerolog.Event {
	return logger.Fatal().Caller()
}
