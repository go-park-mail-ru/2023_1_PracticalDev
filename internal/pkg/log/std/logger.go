// Simple logging package, based on golang built-in log package.
package stdlogger

import (
	"log"
	"os"

	pkgLog "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
)

type StdLogger struct {
	log *log.Logger
}

func New() pkgLog.Logger {
	return &StdLogger{
		log: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}
}

func (log *StdLogger) Info(v ...interface{}) {
	if pkgLog.Info {
		log.log.Println(colorBlue, "INF:", colorReset, v)
	}
}

func (log *StdLogger) Warn(v ...interface{}) {
	if pkgLog.Warn {
		log.log.Println(colorYellow, "WRN:", colorReset, v)
	}
}

func (log *StdLogger) Error(v ...interface{}) {
	if pkgLog.Err {
		log.log.Println(colorRed, "ERR:", colorReset, v)
	}
}

func (log *StdLogger) Debug(v ...interface{}) {
	if pkgLog.Debug {
		log.log.Println(colorBlack, "DBG:", colorReset, v)
	}
}
