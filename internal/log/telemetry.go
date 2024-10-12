package log

import (
	"context"
	"io"

	"github.com/valyala/fastjson"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
)

/*
https://github.com/rs/zerolog/pull/682#issuecomment-2395677680
*/

type config struct {
	provider  log.LoggerProvider
	version   string
	schemaURL string
}

func newConfig(options []Option) config {
	var c config
	for _, opt := range options {
		c = opt.apply(c)
	}

	if c.provider == nil {
		c.provider = global.GetLoggerProvider()
	}
	return c
}

func (c config) logger(name string) log.Logger {
	var opts []log.LoggerOption
	if c.version != "" {
		opts = append(opts, log.WithInstrumentationVersion(c.version))
	}
	if c.schemaURL != "" {
		opts = append(opts, log.WithSchemaURL(c.schemaURL))
	}
	return c.provider.Logger(name, opts...)
}

// Option configures a Hook.
type Option interface {
	apply(config) config
}
type optFunc func(config) config

func (f optFunc) apply(c config) config { return f(c) }

// WithVersion returns an [Option] that configures the version of the
// [log.Logger] used by a [OtelLogger]. The version should be the version of the
// package that is being logged.
func WithVersion(version string) Option {
	return optFunc(func(c config) config {
		c.version = version
		return c
	})
}

// WithSchemaURL returns an [Option] that configures the semantic convention
// schema URL of the [log.Logger] used by a [OtelLogger]. The schemaURL should be
// the schema URL for the semantic conventions used in log records.
func WithSchemaURL(schemaURL string) Option {
	return optFunc(func(c config) config {
		c.schemaURL = schemaURL
		return c
	})
}

// WithLoggerProvider returns an [Option] that configures [log.LoggerProvider]
//
// By default if this Option is not provided, the Logger will use the global LoggerProvider.
func WithLoggerProvider(provider log.LoggerProvider) Option {
	return optFunc(func(c config) config {
		c.provider = provider
		return c
	})
}

var _ io.WriteCloser = (*OtelLogger)(nil)

type OtelLogger struct {
	logger log.Logger
}

func NewOtelLogger(name string, options ...Option) *OtelLogger {
	cfg := newConfig(options)
	return &OtelLogger{
		logger: cfg.logger(name),
	}
}

func (l *OtelLogger) Write(p []byte) (n int, err error) {
	r := log.Record{}
	r.SetSeverity(log.SeverityInfo)
	r.SetBody(log.BytesValue(p))
	r.SetSeverityText(fastjson.GetString(p, "level"))
	l.logger.Emit(context.Background(), r)
	return len(p), nil
}

// Close implements io.Closer, and closes the current logfile.
func (l *OtelLogger) Close() error {
	return nil
}
