package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type (
	Level  = zerolog.Level
	Logger = zerolog.Logger
)

const (
	DebugLevel = zerolog.DebugLevel
	InfoLevel  = zerolog.InfoLevel
	WarnLevel  = zerolog.WarnLevel
	ErrorLevel = zerolog.ErrorLevel
	FatalLevel = zerolog.FatalLevel
	PanicLevel = zerolog.PanicLevel
	NoLevel    = zerolog.NoLevel
	Disabled   = zerolog.Disabled
	TraceLevel = zerolog.TraceLevel
)

func New(level Level) Logger {
	return zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()
}
