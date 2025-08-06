package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Errorf(format string, args ...interface{})
	With(fields ...Field) Logger
}

type Field struct {
	Key   string
	Value interface{}
}

func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

type stdLogger struct {
	out    *log.Logger
	fields []Field
}

func New() Logger {
	return &stdLogger{
		out: log.New(os.Stdout, "", 0),
	}
}

func (l *stdLogger) With(fields ...Field) Logger {
	return &stdLogger{
		out:    l.out,
		fields: append(l.fields, fields...),
	}
}

func (l *stdLogger) Debug(msg string, fields ...Field) {
	l.log(LevelDebug, msg, fields...)
}

func (l *stdLogger) Info(msg string, fields ...Field) {
	l.log(LevelInfo, msg, fields...)
}

func (l *stdLogger) Warn(msg string, fields ...Field) {
	l.log(LevelWarn, msg, fields...)
}

func (l *stdLogger) Error(msg string, fields ...Field) {
	l.log(LevelError, msg, fields...)
}

func (l *stdLogger) Errorf(format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	l.log(LevelError, formatted)
}

func (l *stdLogger) log(level Level, msg string, fields ...Field) {
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	allFields := append(l.fields, fields...)
	var fieldParts []string
	for _, f := range allFields {
		fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", f.Key, f.Value))
	}
	fieldsStr := ""
	if len(fieldParts) > 0 {
		fieldsStr = " " + strings.Join(fieldParts, " ")
	}

	color := levelColor(level)
	reset := "\033[0m"

	line := fmt.Sprintf("%s%-7s%s %s %s%s", color, level, reset, timestamp, msg, fieldsStr)
	l.out.Println(line)
}

func levelColor(level Level) string {
	switch level {
	case LevelDebug:
		return "\033[36m" // Cyan
	case LevelInfo:
		return "\033[32m" // Green
	case LevelWarn:
		return "\033[33m" // Yellow
	case LevelError:
		return "\033[31m" // Red
	default:
		return ""
	}
}
