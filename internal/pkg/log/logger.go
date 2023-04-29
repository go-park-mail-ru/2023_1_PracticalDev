package log

import (
	"os"
	"strings"
)

// Logger is a basic logging interface. It supports 4 log levels:
// INFO, WARN, ERROR, DEBUG,
// configurable by LOG_LEVELS env var.
//
// Default levels are Info, Warn, Error.
type Logger interface {
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Debug(v ...interface{})
}

var (
	Warn  = true
	Err   = true
	Debug = true
	Info  = true
)

func init() {
	levels := os.Getenv("LOG_LEVELS")
	if levels == "" {
		return
	}

	if levels == "*" {
		Debug = true
		return
	}

	modes := strings.Split(levels, ",")

	Info, Warn, Err, Debug = false, false, false, false

	for _, mode := range modes {
		mode = strings.ToUpper(mode)
		switch mode {
		case "INFO":
			Info = true
		case "WARN":
			Warn = true
		case "ERROR":
			Err = true
		case "DEBUG":
			Debug = true
		}
	}
}
