package log

import (
	"context"
	"log/slog"
	"os"
)

type (
	Level  = slog.Level
	Logger = slog.Logger
	key    = struct{}
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelFatal = slog.Level(5)
	LevelPanic = slog.Level(6)
	LevelTrace = slog.Level(-5)
)

func Ctx(ctx context.Context) *Logger {
	return ctx.Value(key{}).(*Logger)
}

func WithCtx(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, key{}, logger)
}

func New(level Level) *Logger {
	// writers := io.MultiWriter(NewOtelLogger("d2s"), os.Stdout)
	// return zerolog.New(writers).Level(level).With().Timestamp().Logger().Hook(SeverityHook{})
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
}
