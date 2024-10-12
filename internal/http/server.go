package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/heptiolabs/healthcheck"
	"github.com/mdobak/go-xerrors"

	"github.com/platipy-io/d2s/app"
	"github.com/platipy-io/d2s/app/lorem"
	"github.com/platipy-io/d2s/internal/log"
)

var timeout = 30 * time.Second
var ErrStarting = xerrors.Message("failed starting")
var ErrStopping = xerrors.Message("failed stopping")

func ListenAndServe(logger *log.Logger) error {
	router := chi.NewRouter()
	server := http.Server{Addr: ":8080", Handler: router}
	errChan := make(chan error)
	health := healthcheck.NewHandler()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		logger.Info("received interrupt, closing server...")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		errChan <- xerrors.New(server.Shutdown(ctx))
		cancel()
		close(errChan)
	}()
	router.HandleFunc("/live", health.LiveEndpoint)
	router.HandleFunc("/ready", health.ReadyEndpoint)

	router.Route("/", func(r chi.Router) {
		r.Use(MiddlewareLogger(logger), MiddlewareRecover, MiddlewareOpenTelemetry)
		r.HandleFunc("/", app.Index)
		r.HandleFunc("/lorem", lorem.Index)
		r.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("I'm about to panic!")) // this will send a response 200 as we write to resp
			panic("some unknown reason")
		})
		r.HandleFunc("/wait", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("starting wait\n"))
			time.Sleep(10 * time.Second)
			w.Write([]byte("ending wait\n"))
		})
	})

	logger.Info("starting server on: " + server.Addr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return xerrors.New(ErrStarting, err)
	}
	if err := xerrors.WithWrapper(ErrStopping, <-errChan); err != nil {
		return err
	}
	logger.Info("server stopped properly")
	return nil
}
