package utils

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

var loggerInstance *slog.Logger // Declare as a pointer to *slog.Logger

var once sync.Once

func GetLogger() *slog.Logger {
	once.Do(func() {
		loggerInstance = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	})
	return loggerInstance
}

func Error(log string, args ...interface{}) {
	// Log the error
	GetLogger().Error(fmt.Sprintf(log, args...))
}

func Info(log string, args ...interface{}) {
	// Log the info
	GetLogger().Info(fmt.Sprintf(log, args...))
}
