package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/platipy-io/d2s/config"
	"github.com/platipy-io/d2s/internal/http"
	"github.com/platipy-io/d2s/internal/log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	flags *pflag.FlagSet
	dev   bool
	port  int
	host  string
	level config.LogLevel
)

func init() {
	flags = command.PersistentFlags()
	flags.Var(&level, "level", "Specify logger level; allowed: "+config.LogLevelsStr)
	flags.StringVar(&host, "host", "", "Host to listen to")
	flags.IntVar(&port, "port", 8080, "Port to listen to")
	flags.BoolVar(&dev, "dev", false, "Activate dev mode")
}

func main() {
	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	if dev && !flags.Changed("level") {
		level = config.LogLevel{Level: log.TraceLevel}
	}

	logger := log.New(level.Level)
	err := http.ListenAndServe(http.WithLogger(logger), http.WithHost(host), http.WithPort(port))

	if errors.Is(err, http.ErrStopping) {
		logger.Error("failed to stop server")
	} else if errors.Is(err, http.ErrStarting) {
		logger.Fatal("failed to start server")
	}
	return err
}
