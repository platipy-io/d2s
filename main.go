package main

import (
	"errors"

	"github.com/platipy-io/d2s/internal/http"
	"github.com/platipy-io/d2s/internal/log"
)

func main() {
	logger := log.New(log.TraceLevel)
	err := http.ListenAndServe(logger)

	if errors.Is(err, http.ErrStopping) {
		logger.Error("failed to stop server")
	} else if errors.Is(err, http.ErrStarting) {
		logger.Fatal("failed to start server")
	}
}
