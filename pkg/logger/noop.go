package logger

import "context"

type Noop struct {
}

func NewNoop() Logger {
	return &Noop{}
}

func (n *Noop) Info(v ...interface{})  {}
func (n *Noop) Debug(v ...interface{}) {}
func (n *Noop) Warn(v ...interface{})  {}
func (n *Noop) Error(v ...interface{}) {}
func (n *Noop) Fatal(v ...interface{}) {}
func (n *Noop) Trace(v ...interface{}) {}

func (n *Noop) Infof(format string, v ...interface{})  {}
func (n *Noop) Debugf(format string, v ...interface{}) {}
func (n *Noop) Warnf(format string, v ...interface{})  {}
func (n *Noop) Errorf(format string, v ...interface{}) {}
func (n *Noop) Fatalf(format string, v ...interface{}) {}
func (n *Noop) Tracef(format string, v ...interface{}) {}

func (n *Noop) WithFields(keyValues map[string]interface{}) Logger {
	return n
}
func (n *Noop) WithField(key string, value interface{}) Logger {
	return n
}
func (n *Noop) WithError(err error) Logger {
	return n
}
func (n *Noop) WithTypeOf(obj interface{}) Logger {
	return n
}

func (n *Noop) ToContext(ctx context.Context) context.Context {
	return ctx
}
func (n *Noop) FromContext(ctx context.Context) Logger {
	return n
}

func (n *Noop) Fields() Fields {
	return Fields{}
}
