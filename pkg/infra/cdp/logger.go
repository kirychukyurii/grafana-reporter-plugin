package cdp

import "github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"

type Logger struct {
	logger *log.Logger
}

func NewLogger(logger *log.Logger) *Logger {
	return &Logger{logger: logger}
}

func (l *Logger) Println(args ...interface{}) {
	l.logger.Debug("headless browser", args...)
}
