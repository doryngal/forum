package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
}

type stdLogger struct {
	base *slog.Logger
}

func New() Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return &stdLogger{
		base: slog.New(handler),
	}
}

func (l *stdLogger) Info(msg string, keysAndValues ...interface{}) {
	l.base.Info(msg, keysAndValues...)
}

func (l *stdLogger) Error(msg string, keysAndValues ...interface{}) {
	l.base.Error(msg, keysAndValues...)
}

func (l *stdLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.base.Debug(msg, keysAndValues...)
}
