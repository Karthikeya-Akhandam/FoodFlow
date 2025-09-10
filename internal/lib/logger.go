package lib

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	// Configure console writer for development
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		if level, ok := i.(zerolog.Level); ok {
			return zerolog.LevelFieldMarshalFunc(level)
		}
		if str, ok := i.(string); ok {
			return str
		}
		return ""
	}
	output.FormatMessage = func(i interface{}) string {
		return i.(string)
	}
	output.FormatFieldName = func(i interface{}) string {
		return i.(string) + "="
	}
	output.FormatFieldValue = func(i interface{}) string {
		return i.(string)
	}

	log.Logger = zerolog.New(output).With().Timestamp().Logger()

	// Set log level based on environment
	if os.Getenv("LOG_LEVEL") == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func GetLogger() zerolog.Logger {
	return log.Logger
}
