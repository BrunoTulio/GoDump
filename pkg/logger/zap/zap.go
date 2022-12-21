package zap

import (
	"context"
	"os"
	"path"
	"reflect"

	"github.com/BrunoTulio/GoDump/pkg/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	KeyContext = "ctxfields"
	KeyError   = "err"
)

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
	core          zapcore.Core
	fields        logger.Fields
}

func (l *zapLogger) Info(v ...interface{}) {
	l.sugaredLogger.Info(v...)
}
func (l *zapLogger) Debug(v ...interface{}) {
	l.sugaredLogger.Debug(v...)
}
func (l *zapLogger) Warn(v ...interface{}) {
	l.sugaredLogger.Warn(v...)
}
func (l *zapLogger) Error(v ...interface{}) {
	l.sugaredLogger.Error(v...)
}
func (l *zapLogger) Fatal(v ...interface{}) {
	l.sugaredLogger.Fatal(v...)
}
func (l *zapLogger) Trace(v ...interface{}) {
	l.sugaredLogger.Debug(v...)
}

func (l *zapLogger) Infof(format string, v ...interface{}) {
	l.sugaredLogger.Infof(format, v...)
}
func (l *zapLogger) Debugf(format string, v ...interface{}) {
	l.sugaredLogger.Debugf(format, v...)
}
func (l *zapLogger) Warnf(format string, v ...interface{}) {
	l.sugaredLogger.Warnf(format, v...)
}
func (l *zapLogger) Errorf(format string, v ...interface{}) {
	l.sugaredLogger.Errorf(format, v...)
}
func (l *zapLogger) Fatalf(format string, v ...interface{}) {
	l.sugaredLogger.Fatalf(format, v...)
}
func (l *zapLogger) Tracef(format string, v ...interface{}) {
	l.sugaredLogger.Debugf(format, v...)
}

func (l *zapLogger) WithFields(fields map[string]interface{}) logger.Logger {
	newFields := logger.Fields{}

	for k, v := range l.fields {
		newFields[k] = v
	}

	for k, v := range fields {
		newFields[k] = v
	}

	f := mapToSlice(newFields)
	newLogger := newSugaredLogger(l.core).With(f...)
	return &zapLogger{newLogger, l.core, newFields}
}

// WithField constructs a new Logger with l.fields and provided key and value field.
func (l *zapLogger) WithField(key string, value interface{}) logger.Logger {
	newFields := logger.Fields{}
	for k, v := range l.fields {
		newFields[k] = v
	}

	newFields[key] = value

	f := mapToSlice(newFields)
	newLogger := newSugaredLogger(l.core).With(f...)
	return &zapLogger{
		sugaredLogger: newLogger,
		core:          l.core,
		fields:        newFields,
	}
}
func (l *zapLogger) WithError(err error) logger.Logger {
	return l.WithField(KeyError, err.Error())
}

func (l *zapLogger) WithTypeOf(obj interface{}) logger.Logger {
	t := reflect.TypeOf(obj)

	return l.WithFields(logger.Fields{
		"reflect.type.name":    t.Name(),
		"reflect.type.package": t.PkgPath(),
	})
}

func (l *zapLogger) ToContext(ctx context.Context) context.Context {
	fields := l.Fields()

	ctxFields := fieldsFromContext(ctx)

	for k, v := range fields {
		ctxFields[k] = v
	}

	return context.WithValue(ctx, KeyContext, ctxFields)
}
func (l *zapLogger) FromContext(ctx context.Context) logger.Logger {
	fields := fieldsFromContext(ctx)
	return l.WithFields(fields)
}

func (l *zapLogger) Fields() logger.Fields {
	return l.fields
}

func newSugaredLogger(core zapcore.Core) *zap.SugaredLogger {
	return zap.New(core,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	).Sugar()
}

func NewZapLogger() logger.Logger {
	var cores []zapcore.Core

	pathLogger, _ := os.Getwd()
	pathLogger = path.Join(pathLogger, "logger.json")

	logFile, _ := os.OpenFile(pathLogger, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig) //zapcore.NewJSONEncoder(encoderConfig)

	level := zapcore.DebugLevel
	writer := zapcore.Lock(os.Stdout)

	coreconsole := zapcore.NewCore(consoleEncoder, writer, level)
	corefile := zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), level)

	cores = append(cores, coreconsole)
	cores = append(cores, corefile)

	combinedCore := zapcore.NewTee(cores...)
	zaplogger := newSugaredLogger(combinedCore)

	newlogger := &zapLogger{
		sugaredLogger: zaplogger,
		core:          combinedCore,
		fields:        logger.Fields{},
	}

	return newlogger

}

func fieldsFromContext(ctx context.Context) logger.Fields {
	fields := make(logger.Fields)

	if ctx == nil {
		return fields
	}

	if f, ok := ctx.Value(KeyContext).(logger.Fields); ok && f != nil {
		for k, v := range f {
			fields[k] = v
		}
	}

	return fields
}

func mapToSlice(m logger.Fields) []interface{} {
	f := make([]interface{}, 2*len(m))
	i := 0
	for k, v := range m {
		f[i] = k
		f[i+1] = v
		i = i + 2
	}

	return f
}
