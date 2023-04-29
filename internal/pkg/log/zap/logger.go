package zaplogger

import (
	"go.uber.org/zap"

	pkgLog "github.com/go-park-mail-ru/2023_1_PracticalDev/internal/pkg/log"
)

type ZapLogger struct {
	log *zap.SugaredLogger
}

func New() (*ZapLogger, error) {
	log, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{log: log.Sugar()}, nil
}

func (log *ZapLogger) Sync() error {
	return log.log.Sync()
}

func (log *ZapLogger) Info(args ...interface{}) {
	if pkgLog.Info {
		log.log.Info(args)
	}
}

func (log *ZapLogger) Warn(args ...interface{}) {
	if pkgLog.Warn {
		log.log.Warn(args)
	}
}

func (log *ZapLogger) Error(args ...interface{}) {
	if pkgLog.Err {
		log.log.Error(args)
	}
}

func (log *ZapLogger) Debug(args ...interface{}) {
	if pkgLog.Debug {
		log.log.Debug(args)
	}
}
