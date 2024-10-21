package config

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/platipy-io/d2s/internal/log"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

type (
	// Configuration hold the current fields to tune the application
	Configuration struct {
		Dev  bool
		Host string
		Port int
		Logger
		Tracer
	}

	Logger struct {
		Level string
	}

	Tracer struct {
		Endpoint string
		Headers  map[string]string
	}
)

func New(path string) (config Configuration, err error) {
	viper.SetConfigFile(path)
	if err = viper.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

func (t Tracer) Opts() (opts []otlptracehttp.Option) {
	if t.Endpoint != "" {
		opts = append(opts, otlptracehttp.WithEndpointURL(t.Endpoint))
	}
	if len(t.Headers) != 0 {
		opts = append(opts, otlptracehttp.WithHeaders(t.Headers))
	}
	return opts
}

func (c Configuration) NewLogger(level zerolog.Level, override bool) zerolog.Logger {
	var output io.Writer = os.Stdout

	if c.Dev {
		output = zerolog.ConsoleWriter{Out: os.Stdout}
		zerolog.ErrorStackMarshaler = log.MarshalStackDev
		if !override {
			level = zerolog.TraceLevel
		}
	} else {
		zerolog.ErrorStackMarshaler = log.MarshalStack
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano

	return zerolog.New(output).Level(level).With().Timestamp().Logger()
}
