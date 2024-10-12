package main

import (
	"context"
	"errors"

	"github.com/platipy-io/d2s/internal/http"
	"github.com/platipy-io/d2s/internal/log"
)

func main() {
	logger := log.New(log.LevelTrace)
	err := http.ListenAndServe(logger)

	if errors.Is(err, http.ErrStopping) {
		logger.Error("failed to stop server")
	} else if errors.Is(err, http.ErrStarting) {
		logger.Log(context.TODO(), log.LevelFatal, "failed to start server")
	}
}
