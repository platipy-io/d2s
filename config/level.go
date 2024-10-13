package config

import (
	"strings"

	"github.com/platipy-io/d2s/internal/log"
	"go.uber.org/zap/zapcore"
)

type LogLevel struct {
	zapcore.Level
}

func JoinQuoted(in []string, sep string) string {
	var b strings.Builder
	if len(in) == 0 {
		return ""
	}
	b.WriteString(`"`)
	b.WriteString(in[0])
	b.WriteString(`"`)
	for _, str := range in[1:] {
		b.WriteString(`,"`)
		b.WriteString(str)
		b.WriteString(`"`)
	}
	return b.String()
}

func (ll LogLevel) String() string { return ll.Level.String() }
func (ll LogLevel) Type() string   { return "<level>" }
func (ll *LogLevel) Set(v string) (err error) {
	if v == "trace" {
		ll.Level = log.TraceLevel
	}
	ll.Level, err = zapcore.ParseLevel(v)
	return
}

var LogLevelsStr = JoinQuoted([]string{
	"trace",
	log.DebugLevel.String(),
	log.InfoLevel.String(),
	log.WarnLevel.String(),
	log.ErrorLevel.String(),
}, ", ")
