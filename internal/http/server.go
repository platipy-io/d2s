package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/heptiolabs/healthcheck"
	"github.com/mdobak/go-xerrors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/platipy-io/d2s/app"
	"github.com/platipy-io/d2s/app/lorem"
	"github.com/platipy-io/d2s/internal/log"
)

var timeout = 30 * time.Second
var ErrStarting = xerrors.Message("failed starting")
var ErrStopping = xerrors.Message("failed stopping")

type serverConfig struct {
	host   string
	port   int
	logger log.Logger
}

func (sc serverConfig) addr() string {
	return sc.host + ":" + strconv.Itoa(sc.port)
}

// ServerOption applies a configuration option value to a Server.
type ServerOption interface {
	apply(serverConfig) serverConfig
}

type ServerOptionFunc func(serverConfig) serverConfig

func (fn ServerOptionFunc) apply(c serverConfig) serverConfig {
	return fn(c)
}

func newServerConfig(opts []ServerOption) serverConfig {
	sc := serverConfig{port: 8080, logger: log.Nop()}
	for _, opt := range opts {
		sc = opt.apply(sc)
	}
	return sc
}

func WithHost(host string) ServerOption {
	return ServerOptionFunc(func(sc serverConfig) serverConfig {
		sc.host = host
		return sc
	})
}

func WithLogger(logger log.Logger) ServerOption {
	return ServerOptionFunc(func(sc serverConfig) serverConfig {
		sc.logger = logger
		return sc
	})
}

func WithPort(port int) ServerOption {
	return ServerOptionFunc(func(sc serverConfig) serverConfig {
		sc.port = port
		return sc
	})
}

func ListenAndServe(opts ...ServerOption) error {
	router := chi.NewRouter()
	errChan := make(chan error)
	health := healthcheck.NewHandler()
	config := newServerConfig(opts)
	logger := config.logger
	server := http.Server{Addr: config.addr(), Handler: router}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		logger.Info().Msg("received interrupt, closing server...")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		errChan <- xerrors.New(server.Shutdown(ctx))
		cancel()
		close(errChan)
	}()
	router.HandleFunc("/live", health.LiveEndpoint)
	router.HandleFunc("/ready", health.ReadyEndpoint)
	router.Handle("/metrics", promhttp.Handler())

	router.Route("/", func(r chi.Router) {
		r.Use(MiddlewareOpenTelemetry, MiddlewareMetrics, MiddlewareLogger(logger), MiddlewareRecover)
		r.HandleFunc("/", app.Index)
		r.HandleFunc("/lorem", lorem.Index)
		r.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
			// w.Write([]byte("I'm about to panic!")) // this will send a response 200 as we write to resp
			panic("some unknown reason")
		})
		r.HandleFunc("/wait", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("starting wait\n"))
			time.Sleep(10 * time.Second)
			w.Write([]byte("ending wait\n"))
		})
	})

	logger.Info().Msg("starting server on: " + server.Addr)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return xerrors.New(ErrStarting, err)
	}
	if err := xerrors.WithWrapper(ErrStopping, <-errChan); err != nil {
		return err
	}
	logger.Info().Msg("server stopped properly")
	return nil
}
