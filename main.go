package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/platipy-io/d2s/config"
	"github.com/platipy-io/d2s/internal/http"
	"github.com/platipy-io/d2s/internal/log"
	"github.com/platipy-io/d2s/internal/telemetry"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// DefaultConfigPath is the default location the application will use to
	// find the configuration.
	DefaultConfigPath = "d2s.toml"
)

var (
	command = &cobra.Command{
		Use:   "d2s [flags]",
		Short: ".",
		RunE: func(cmd *cobra.Command, args []string) error {
			// https://github.com/spf13/cobra/issues/340
			cmd.SilenceUsage = true
			return run()
		},
	}
	flags    *pflag.FlagSet
	logLevel = config.LogLevel{Level: log.InfoLevel}
)

func init() {
	flags = command.PersistentFlags()
	flags.Var(&logLevel, "level", "Specify logger level; allowed: "+config.LogLevelsStr)
	flags.String("host", "", "Host to listen to")
	flags.Int("port", 8080, "Port to listen to")
	flags.Bool("dev", false, "Activate dev mode")
	flags.String("config", DefaultConfigPath, "Path to a configuration file")

	_ = viper.BindPFlag("host", flags.Lookup("host"))
	_ = viper.BindPFlag("port", flags.Lookup("port"))
	_ = viper.BindPFlag("dev", flags.Lookup("dev"))
	_ = viper.BindPFlag("logger.level", flags.Lookup("level"))
}

func main() {
	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	conf, err := config.New(DefaultConfigPath)
	if err != nil {
		return err
	}
	if err := telemetry.InitTrace("d2s", otlptracehttp.WithInsecure()); err != nil {
		return err
	}

	if conf.Dev && !flags.Changed("level") {
		logLevel.Level = log.TraceLevel
		conf.Logger.Level = logLevel.String()
	}
	logger := log.New(logLevel.Level)
	logger.Debug().Object("config", conf).Msg("dumping config")
	opts := []http.ServerOption{http.WithLogger(logger),
		http.WithHost(conf.Host), http.WithPort(conf.Port)}

	err = http.ListenAndServe(opts...)
	if errors.Is(err, http.ErrStopping) {
		logger.Error().Msg("failed to stop server")
	} else if errors.Is(err, http.ErrStarting) {
		logger.Fatal().Msg("failed to start server")
	}
	return err
}
