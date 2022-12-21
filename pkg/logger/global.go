package logger

import (
	"context"
	"sync"
)

var (
	l  Logger = NewNoop()
	mu sync.RWMutex
)

func NewLogger(logger Logger) Logger {
	mu.Lock()
	defer mu.Unlock()
	l = logger

	return l
}

func Info(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	l.Info(v...)
}
func Debug(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	l.Debug(v...)
}
func Warn(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	l.Warn(v...)
}
func Error(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	l.Error(v...)
}
func Fatal(v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	l.Fatal(v...)
}

func Infof(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	l.Infof(format, v...)
}
func Debugf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	l.Debugf(format, v...)
}
func Warnf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	l.Warnf(format, v...)
}
func Errorf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	l.Errorf(format, v...)
}
func Fatalf(format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	l.Fatalf(format, v...)
}

func WithFields(keyValues map[string]interface{}) Logger {
	return l.WithFields(keyValues)
}
func WithField(key string, value interface{}) Logger {
	return l.WithField(key, value)
}
func WithError(err error) Logger {
	return l.WithError(err)
}
func WithTypeOf(obj interface{}) Logger {
	return l.WithTypeOf(obj)
}

func ToContext(ctx context.Context) context.Context {
	return l.ToContext(ctx)
}
func FromContext(ctx context.Context) Logger {
	return l.FromContext(ctx)
}

// GetLogger returns instance of Logger.
func GetLogger() Logger {
	return l
}
