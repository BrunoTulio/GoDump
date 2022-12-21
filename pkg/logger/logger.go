package logger

import "context"

type Logger interface {
	Info(v ...interface{})
	Debug(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Trace(v ...interface{})

	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	Tracef(format string, v ...interface{})

	WithFields(keyValues map[string]interface{}) Logger
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	WithTypeOf(obj interface{}) Logger

	ToContext(ctx context.Context) context.Context
	FromContext(ctx context.Context) Logger

	Fields() Fields
}
