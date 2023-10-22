package log

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"os"
)

type Logger struct {
	log.Logger
}

func New() *Logger {
	return &Logger{
		Logger: log.New(),
	}
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.Error(msg, args...)
	os.Exit(1)
}
