package http

import (
	"net/http"

	"github.com/platipy-io/d2s/internal/log"
	"github.com/platipy-io/d2s/internal/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/mdobak/go-xerrors"
)

func MiddlewareRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer xerrors.Recover(func(err error) {
			w.WriteHeader(http.StatusInternalServerError)
			log.Ctx(r.Context()).Error("recovering from panic!")
		})
		next.ServeHTTP(w, r)
	})
}

var MiddlewareLogger = log.Middleware

var MiddlewareOpenTelemetry = otelhttp.NewMiddleware("fooo",
	otelhttp.WithTracerProvider(telemetry.Provider),
	otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
		return "server"
	}),
)
