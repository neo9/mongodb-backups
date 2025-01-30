package log

import (
	"fmt"
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
	GetLogger().Error(fmt.Sprintf(log, args...))
}

func Info(log string, args ...interface{}) {
	GetLogger().Info(fmt.Sprintf(log, args...))
}

func Debug(log string, args ...interface{}) {
	GetLogger().Debug(fmt.Sprintf(log, args...))
}

func Warn(log string, args ...interface{}) {
	GetLogger().Warn(fmt.Sprintf(log, args...))
}
