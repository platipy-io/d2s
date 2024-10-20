package config

import (
	"errors"
	"os"

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
