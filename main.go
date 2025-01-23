// Custom logger for a gin app
package logger

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	REQUEST_HEADER = "X-Request-ID"
	colorReset     = "\033[0m"
	colorRed       = "\033[31m"
	colorYellow    = "\033[33m"
	colorBlue      = "\033[34m"
	colorGreen     = "\033[32m"
)

// Passing a Gin Context to the logger is for extracting the Request ID from the headers.
// Request ID headers key : X-Request-ID
func Debug(context ...*gin.Context) *zerolog.Event {
	return getLogger("DEBUG", context...)
}

// Passing a Gin Context to the logger is for extracting the Request ID from the headers.
// Request ID headers key : X-Request-ID
func Warn(context ...*gin.Context) *zerolog.Event {
	return getLogger("WARN", context...)
}

// Passing a Gin Context to the logger is for extracting the Request ID from the headers.
// Request ID headers key : X-Request-ID
func Info(context ...*gin.Context) *zerolog.Event {
	return getLogger("INFO", context...)
}

// Passing a Gin Context to the logger is for extracting the Request ID from the headers.
// Request ID headers key : X-Request-ID
func Error(context ...*gin.Context) *zerolog.Event {
	return getLogger("ERROR", context...)
}

// This function should be called in your app so you can have a custom log output.
// Calling this function is a must to format the log as needed.
func ConfigLoader() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatTimestamp: func(i interface{}) string {
			return colorBlue + "[" + i.(string) + "]" + colorReset
		},
		FormatLevel: func(i interface{}) string {
			return "|" + getColoredDebug(i.(string)) + "|"
		},
		FormatMessage: func(i interface{}) string {
			return "- " + i.(string)
		},
		FormatFieldName: func(i interface{}) string {
			return i.(string) + ":"
		},
		FormatFieldValue: func(i interface{}) string {
			return i.(string)
		},
	}
	log.Logger = zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()
}

// Set the global log level for your application.
// Log level should be one of the following : DEBUG / WARN / INFO / ERROR
func SetLoggerLevel(logLevel string) {
	logger_levels := map[string]zerolog.Level{
		"DEBUG": zerolog.DebugLevel,
		"INFO":  zerolog.InfoLevel,
		"ERROR": zerolog.ErrorLevel,
		"WARN":  zerolog.WarnLevel,
	}
	if _, ok := logger_levels[logLevel]; !ok {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	zerolog.SetGlobalLevel(logger_levels[logLevel])
}

func getColoredDebug(level string) string {
	switch strings.ToUpper(level) {
	case "INFO":
		return colorYellow + level + colorReset
	case "WARN", "ERROR":
		return colorRed + level + colorReset
	default:
		return level
	}
}

func getRequestID(c *gin.Context) string {
	return c.Request.Header.Get(REQUEST_HEADER)
}

func getLogger(level string, context ...*gin.Context) *zerolog.Event {
	var logger *zerolog.Event
	loggers := map[string]*zerolog.Event{
		"DEBUG": log.Debug(),
		"WARN":  log.Warn(),
		"INFO":  log.Info(),
		"ERROR": log.Error(),
	}
	if _, ok := loggers[level]; !ok {
		logger = log.Debug()
	} else {
		logger = loggers[level]
	}
	if len(context) == 0 {
		return logger
	}
	c := context[0]
	return logger.Str("request-id", colorGreen+getRequestID(c)+colorReset)
}
