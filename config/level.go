package config

import (
	"strings"

	"github.com/rs/zerolog"
)

type LogLevel struct {
	zerolog.Level
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
	ll.Level, err = zerolog.ParseLevel(v)
	return
}

var LogLevelsStr = JoinQuoted([]string{
	zerolog.LevelTraceValue,
	zerolog.LevelDebugValue,
	zerolog.LevelInfoValue,
	zerolog.LevelWarnValue,
	zerolog.LevelErrorValue,
	zerolog.LevelFatalValue,
	zerolog.LevelPanicValue,
}, ", ")
