package log

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/rs/xid"
)

func mustRead(reader io.Reader) []byte {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return bytes
}

func middleware(logger *Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		child := logger.With(slog.String("req_id", xid.New().String()))
		child.Info("starting request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
			slog.String("user_agent", r.UserAgent()))
		if r.ContentLength != 0 {
			// child.
			// 	Trace().
			// 	Func(func(e *zerolog.Event) {
			// 		e.Dict("headers", RequestHeaders(r.Header))
			// 		body := mustRead(r.Body)
			// 		r.Body = io.NopCloser(bytes.NewBuffer(body))
			// 		if len(body) == 0 {
			// 			return
			// 		}
			// 		if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			// 			e.RawJSON("body", body)
			// 		} else {
			// 			e.Bytes("body", body)
			// 		}
			// 	}).Msg("dumping request")
		}
		child.Info("ending request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
			slog.String("user_agent", r.UserAgent()))
		next.ServeHTTP(w, r.WithContext(WithCtx(r.Context(), logger)))
		// hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		// 	child.
		// 		Info().
		// 		Int("status", status).
		// 		Int("size", size).
		// 		Dur("elapsed_ms", duration).
		// 		Msg("ending request")
		// })(next).ServeHTTP(w, r.WithContext(child.WithContext(r.Context())))
	})
}

func Middleware(logger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware(logger, next)
	}
}
