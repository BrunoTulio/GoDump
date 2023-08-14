package zap

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/BrunoTulio/GoDump/pkg/logger"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	fields        logger.Fields
	sugaredLogger *zap.SugaredLogger
	writers       []io.Writer
	core          zapcore.Core
}

// Debug implements logger.Logger.
func (l *zapLogger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

// Debugf implements logger.Logger.
func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

// Error implements logger.Logger.
func (l *zapLogger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

// Errorf implements logger.Logger.
func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

// Fatal implements logger.Logger.
func (l *zapLogger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

// Fatalf implements logger.Logger.
func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}

// Info implements logger.Logger.
func (l *zapLogger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

// Infof implements logger.Logger.
func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

// Panic implements logger.Logger.
func (l *zapLogger) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(args...)
}

// Panicf implements logger.Logger.
func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Panicf(format, args...)
}

// Printf implements logger.Logger.
func (l *zapLogger) Printf(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

// Trace implements logger.Logger.
func (l *zapLogger) Trace(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

// Tracef implements logger.Logger.
func (l *zapLogger) Tracef(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

// Warn implements logger.Logger.
func (l *zapLogger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

// Warnf implements logger.Logger.
func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

// Output returns a Writer that represents the zap writers.
func (l *zapLogger) Output() io.Writer {
	return io.MultiWriter(l.writers...)
}

// WithField constructs a new Logger with l.fields and provided key and value field.
func (l *zapLogger) WithField(key string, value interface{}) logger.Logger {
	newFields := logger.Fields{}
	for k, v := range l.fields {
		newFields[k] = v
	}

	newFields[key] = value

	values := mapToSlice(newFields)

	newLogger := newSugaredLogger(l.core).With(values...)
	return &zapLogger{newFields, newLogger, l.writers, l.core}
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
	return &zapLogger{newFields, newLogger, l.writers, l.core}
}

// Fields implements logger.Logger
func (l *zapLogger) Fields() logger.Fields {
	return l.fields
}

// FromContext implements logger.Logger
func (l *zapLogger) FromContext(ctx context.Context) logger.Logger {
	fields := fieldsFromContext(ctx)
	return l.WithFields(fields)
}

func (l *zapLogger) ToContext(ctx context.Context) context.Context {
	fields := l.Fields()

	ctxFields := fieldsFromContext(ctx)

	for k, v := range fields {
		ctxFields[k] = v
	}

	return context.WithValue(ctx, logger.ContextKey, ctxFields)
}

func fieldsFromContext(ctx context.Context) logger.Fields {
	fields := make(logger.Fields)

	if ctx == nil {
		return fields
	}

	if f, ok := ctx.Value(logger.ContextKey).(logger.Fields); ok && f != nil {
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

func NewWithOptions(options *Options) logger.Logger {

	cores := []zapcore.Core{}
	var writers []io.Writer

	if options.Console.Enabled {
		level := logLevel(options.Console.Level)
		writer := zapcore.Lock(os.Stdout)
		coreconsole := zapcore.NewCore(getEncoder(options.Console.Formatter), writer, level)
		cores = append(cores, coreconsole)
		writers = append(writers, writer)
	}

	if options.File.Enabled {
		s := []string{options.File.Path, "/", options.File.Name}
		fileLocation := strings.Join(s, "")

		lumber := &lumberjack.Logger{
			Filename: fileLocation,
			MaxSize:  options.File.MaxSize,
			Compress: options.File.Compress,
			MaxAge:   options.File.MaxAge,
		}

		level := logLevel(options.File.Level)
		writer := zapcore.AddSync(lumber)
		corefile := zapcore.NewCore(getEncoder(options.File.Formatter), writer, level)
		cores = append(cores, corefile)
		writers = append(writers, lumber)
	}

	combinedCore := zapcore.NewTee(cores...)

	// AddCallerSkip skips 2 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	zaplogger := newSugaredLogger(combinedCore)

	newlogger := &zapLogger{
		sugaredLogger: zaplogger,
		writers:       writers,
		core:          combinedCore,
	}

	logger.SetGlobalLogger(newlogger)

	return newlogger
}

func newSugaredLogger(core zapcore.Core) *zap.SugaredLogger {
	return zap.New(core,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	).Sugar()
}

func getEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	switch format {
	case "JSON":
		return zapcore.NewJSONEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func logLevel(level string) zapcore.Level {
	switch level {
	case "TRACE":
		return zapcore.DebugLevel
	case "WARN":
		return zapcore.WarnLevel
	case "DEBUG":
		return zapcore.DebugLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "PANIC":
		return zapcore.PanicLevel
	case "FATAL":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
