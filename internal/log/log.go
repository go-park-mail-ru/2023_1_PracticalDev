// Simple logging package, based on golang built-in log package.
package log

import (
	"log"
	"os"
	"strings"
)

var (
	info, warn, err, debug bool = true, true, true, true
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

type logger struct {
	log *log.Logger
}

func init() {
	levels := os.Getenv("LOG_LEVELS")
	if levels == "" {
		return
	}

	if levels == "*" {
		debug = true
		return
	}

	modes := strings.Split(levels, ",")

	info, warn, err, debug = false, false, false, false

	for _, mode := range modes {
		mode = strings.ToUpper(mode)
		switch mode {
		case "INFO":
			info = true
		case "WARN":
			warn = true
		case "ERROR":
			err = true
		case "DEBUG":
			debug = true
		}
	}
}

func New() Logger {
	a := logger{
		log: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	return Logger(&a)
}

func (log *logger) Info(v ...interface{}) {
	if info {
		log.log.Println((colorBlue), "INF:", (colorReset), v)
	}
}

func (log *logger) Warn(v ...interface{}) {
	if warn {
		log.log.Println((colorYellow), "WRN:", (colorReset), v)
	}
}

func (log *logger) Error(v ...interface{}) {
	if err {
		log.log.Println((colorRed), "ERR:", (colorReset), v)
	}
}

func (log *logger) Debug(v ...interface{}) {
	if debug {
		log.log.Println((colorBlack), "DBG:", (colorReset), v)
	}
}
