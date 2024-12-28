package utils

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger() *Logger {
	logger, _ := zap.NewProduction()
	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}