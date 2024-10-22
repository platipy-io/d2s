package log

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

type MarshalerFunc func(*zerolog.Event)

func (mf MarshalerFunc) MarshalZerologObject(e *zerolog.Event) { mf(e) }

func mustRead(reader io.Reader) []byte {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return bytes
}

func RequestHeaders(h http.Header) *zerolog.Event {
	dict := zerolog.Dict()
	for k, v := range h {
		if len(v) == 1 {
			dict.Str(k, v[0])
		} else {
			dict.Strs(k, v)
		}
	}
	return dict
}

func Request(r *http.Request) zerolog.LogObjectMarshaler {
	return MarshalerFunc(func(e *zerolog.Event) {
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
	})

}
