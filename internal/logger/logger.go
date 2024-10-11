package logger

import (
	"context"
	"fmt"
	"io"
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

func Ctx(ctx context.Context) *Logger {
	return zerolog.Ctx(ctx)
}

func New(level Level) Logger {
	writers := io.MultiWriter(NewOtelLogger("d2s"), os.Stdout)
	return zerolog.New(writers).Level(level).With().Timestamp().Logger().Hook(SeverityHook{})
}

func init() {
	zerolog.ErrorStackMarshaler = MarshalStack
}

type SeverityHook struct{}

func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		fmt.Println("foooo")
	}
}
