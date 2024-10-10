package main

import (
	"errors"

	"github.com/IxDay/templ-exp/internal/http"
	"github.com/IxDay/templ-exp/internal/logger"
)

func main() {
	logger := logger.New(logger.TraceLevel)
	err := http.ListenAndServe(logger)

	if errors.Is(err, http.ErrStopping) {
		logger.Error().Stack().Err(err).Msg("failed to stop server")
	} else if errors.Is(err, http.ErrStarting) {
		logger.Fatal().Stack().Err(err).Msg("failed to start server")
	}
}
