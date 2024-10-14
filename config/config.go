package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type (
	// Configuration hold the current fields to tune the application
	Configuration struct {
		Dev    bool
		Host   string
		Port   int
		Logger ConfigurationLogger
	}

	ConfigurationLogger struct {
		Level string
	}
)

func New(path string) (config Configuration, err error) {
	viper.SetConfigFile(path)
	if err = viper.ReadInConfig(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
