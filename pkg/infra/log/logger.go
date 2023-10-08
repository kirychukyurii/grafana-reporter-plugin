package log

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

type Logger struct {
	log.Logger
}

func New() *Logger {
	return &Logger{
		Logger: log.New(),
	}
}

func (l *Logger) Println(args ...interface{}) {
	l.Debug("headless browser", args...)
}
