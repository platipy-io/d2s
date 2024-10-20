package log

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
)

func handler(logger Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww, ok := w.(middleware.WrapResponseWriter)
		if !ok {
			ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		}
		span := trace.SpanFromContext(r.Context())
		child := logger.With().
			Str("span_id", span.SpanContext().SpanID().String()).
			Str("trace_id", span.SpanContext().TraceID().String()).
			Logger()
		child.Info().
			Str("method", r.Method).
			Str("url", r.URL.Path).
			Str("user_agent", r.UserAgent()).
			Msg("starting request")

		if r.ContentLength != 0 {
			child.Trace().EmbedObject(Request(r)).Msg("dumping request")
		}
		next.ServeHTTP(ww, r.WithContext(child.WithContext(r.Context())))
		child.Info().
			Int("status", ww.Status()).
			Int("size", ww.BytesWritten()).
			Dur("elapsed_ms", time.Since(start)).
			Msg("ending request")
	})
}

func Middleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handler(logger, next)
	}
}
