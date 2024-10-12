package log

import (
	"log/slog"
	"net/http"
	"strings"
)

type LogValuer func() slog.Value

func (lg LogValuer) LogValue() slog.Value { return lg() }

func RequestHeaders(h http.Header) slog.LogValuer {
	f := []slog.Attr{}
	for k, v := range h {
		f = append(f, slog.String(k, strings.Join(v, ",")))
	}
	return LogValuer(func() slog.Value { return slog.GroupValue(f...) })
}
