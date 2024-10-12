package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/platipy-io/d2s/internal/http"
	"github.com/platipy-io/d2s/internal/log"

	"github.com/spf13/cobra"
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
	dev  bool
	port int
	host string
)

func init() {
	flags := command.PersistentFlags()

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
	level := log.WarnLevel
	if dev {
		level = log.TraceLevel
	}

	logger := log.New(level)
	err := http.ListenAndServe(http.WithLogger(logger), http.WithHost(host), http.WithPort(port))

	if errors.Is(err, http.ErrStopping) {
		logger.Error("failed to stop server")
	} else if errors.Is(err, http.ErrStarting) {
		logger.Fatal("failed to start server")
	}
	return err
}
