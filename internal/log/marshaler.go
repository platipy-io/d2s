package log

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap/zapcore"
)

func mustRead(reader io.Reader) []byte {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return bytes
}

func StringArray(slice []string) zapcore.ArrayMarshaler {
	return zapcore.ArrayMarshalerFunc(func(ae zapcore.ArrayEncoder) error {
		for _, s := range slice {
			ae.AppendString(s)
		}
		return nil
	})
}

func RequestHeaders(h http.Header) zapcore.ObjectMarshaler {
	return zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
		for k, v := range h {
			if len(v) == 1 {
				oe.AddString(k, v[0])
			} else {
				oe.AddArray(k, StringArray(v))
			}
		}
		return nil
	})
}

func Request(r *http.Request) zapcore.ObjectMarshaler {
	return zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
		oe.AddObject("headers", RequestHeaders(r.Header))
		body := mustRead(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		if len(body) == 0 {
			return nil
		}
		if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			foo := json.RawMessage(body)
			oe.AddReflected("body", &foo)
		} else {
			oe.AddByteString("body", body)
		}
		return nil
	})
}
