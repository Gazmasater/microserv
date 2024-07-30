package logger

import (
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

// Init initializes the logger
func Init() {
	logger, _ := zap.NewProduction()
	sugar = logger.Sugar()
}

// GetLogger returns the Sugared Logger
func GetLogger() *zap.SugaredLogger {
	return sugar
}
