package log

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/platipy-io/d2s/internal/http/mutil"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

func mustRead(reader io.Reader) []byte {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return bytes
}

func middleware(logger zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := mutil.WrapWriter(w)
		child := logger.With().Str("req_id", xid.New().String()).Logger()
		child.
			Info().
			Str("method", r.Method).
			Str("url", r.URL.Path).
			Str("user_agent", r.UserAgent()).
			Msg("starting request")
		defer func() {
			child.
				Info().
				Int("status", lw.Status()).
				Int("size", lw.BytesWritten()).
				Dur("elapsed_ms", time.Since(start)).
				Msg("ending request")
		}()
		if r.ContentLength != 0 {
			child.
				Trace().
				Func(func(e *zerolog.Event) {
					e.Dict("headers", RequestHeaders(r.Header))
					body := mustRead(r.Body)
					r.Body = io.NopCloser(bytes.NewBuffer(body))
					if len(body) == 0 {
						return
					}
					if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
						e.RawJSON("body", body)
					} else {
						e.Bytes("body", body)
					}
				}).Msg("dumping request")
		}

		next.ServeHTTP(lw, r.WithContext(child.WithContext(r.Context())))
	})
}

func Middleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware(logger, next)
	}
}
