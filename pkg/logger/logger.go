package logger

import (
	"context"
	"io"
)

const (
	ContextKey = iota
)

type Logger interface {
	Printf(format string, args ...interface{})

	Tracef(format string, args ...interface{})

	Trace(args ...interface{})

	Debugf(format string, args ...interface{})

	Debug(args ...interface{})

	Infof(format string, args ...interface{})

	Info(args ...interface{})

	Warnf(format string, args ...interface{})

	Warn(args ...interface{})

	Errorf(format string, args ...interface{})

	Error(args ...interface{})

	Fatalf(format string, args ...interface{})

	Fatal(args ...interface{})

	Panicf(format string, args ...interface{})

	Panic(args ...interface{})

	Output() io.Writer

	WithField(key string, value interface{}) Logger
	Fields() Fields
	WithFields(keyValues map[string]interface{}) Logger

	ToContext(ctx context.Context) context.Context
	FromContext(ctx context.Context) Logger
}
