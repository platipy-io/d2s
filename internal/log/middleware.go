package log

import (
	"net/http"
	"time"

	"github.com/platipy-io/d2s/internal/http/mutil"
	"go.opentelemetry.io/otel/trace"
)

func middleware(logger Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := mutil.WrapWriter(w)
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
		defer func() {
			child.Info().
				Int("status", lw.Status()).
				Int("size", lw.BytesWritten()).
				Dur("elapsed_ms", time.Since(start)).
				Msg("ending request")
		}()
		if r.ContentLength != 0 {
			child.Trace().EmbedObject(Request(r)).Msg("dumping request")
		}
		next.ServeHTTP(lw, r.WithContext(child.WithContext(r.Context())))
	})
}

func Middleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware(logger, next)
	}
}
