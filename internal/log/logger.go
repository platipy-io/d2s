package log

import (
	"context"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/log/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Level  = zapcore.Level
	Logger = zap.Logger
	key    struct{}
)

const (
	TraceLevel = zapcore.Level(-2)
	DebugLevel = zap.DebugLevel
	InfoLevel  = zap.InfoLevel
	WarnLevel  = zap.WarnLevel
	ErrorLevel = zap.ErrorLevel
)

func Ctx(ctx context.Context) *Logger {
	return ctx.Value(key{}).(*Logger)
}

func WithCtx(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, key{}, logger)
}

func New(level Level) *Logger {
	provider := global.GetLoggerProvider()
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(os.Stdout), level),
		otelzap.NewCore("d2s", otelzap.WithLoggerProvider(provider)),
	)
	return zap.New(core)
}
