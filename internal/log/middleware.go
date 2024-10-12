package log

import (
	"net/http"
	"time"

	"github.com/platipy-io/d2s/internal/http/mutil"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func middleware(logger *Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := mutil.WrapWriter(w)
		span := trace.SpanFromContext(r.Context())
		child := logger.With(
			zap.String("span_id", span.SpanContext().SpanID().String()),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)
		child.Info("starting request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.String("user_agent", r.UserAgent()))
		defer func() {
			child.Info("ending request",
				zap.Int("status", lw.Status()),
				zap.Int("size", lw.BytesWritten()),
				zap.Duration("elapsed_ms", time.Since(start)))
		}()
		if r.ContentLength != 0 {
			child.Log(TraceLevel, "dumping request", zap.Inline(Request(r)))
		}
		next.ServeHTTP(lw, r.WithContext(WithCtx(r.Context(), child)))
	})
}

func Middleware(logger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware(logger, next)
	}
}
