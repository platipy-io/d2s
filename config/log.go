package config

import (
	"github.com/rs/zerolog"
)

func (c Configuration) MarshalZerologObject(e *zerolog.Event) {
	e.Bool("dev", c.Dev)
	e.Str("host", c.Host)
	e.Int("port", c.Port)

	e.Object("logger", c.Logger)
}

func (cl ConfigurationLogger) MarshalZerologObject(e *zerolog.Event) {
	e.Str("level", cl.Level)
}
