package logger

import "go.uber.org/zap"

type Logger struct {
	logger *zap.Logger
}

func NewLogger(l *zap.Logger) *Logger {
	return &Logger{logger: l}
}

func (l *Logger) Base() *zap.Logger {
	return l.logger
}

func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.logger.Sugar()
}
