package logger

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

func RequestHeaders(h http.Header) *zerolog.Event {
	dict := zerolog.Dict()
	for k, v := range h {
		dict.Str(k, strings.Join(v, ","))
	}
	return dict
}
