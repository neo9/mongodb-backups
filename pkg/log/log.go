package log

import (
	"log/slog"
	"os"
	"sync"
)

var loggerInstance *slog.Logger

var once sync.Once

func GetLogger() *slog.Logger {
	once.Do(func() {
		loggerInstance = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	})
	return loggerInstance
}

func Error(log string, args ...interface{}) {
	GetLogger().Error(log, args...)
}

func Info(log string, args ...interface{}) {
	GetLogger().Info(log, args...)
}

func Debug(log string, args ...interface{}) {
	GetLogger().Debug(log, args...)
}
