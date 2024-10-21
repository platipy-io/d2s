package log

import (
	"context"

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

func Ctx(ctx context.Context) *Logger {
	return zerolog.Ctx(ctx)
}

func Nop() Logger {
	return zerolog.Nop()
}
