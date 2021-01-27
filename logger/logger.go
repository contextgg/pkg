package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
}

type logger struct {
	z *zap.SugaredLogger
}

func (l *logger) Debug(msg string, keysAndValues ...interface{}) {
	l.z.Debugw(msg, keysAndValues...)
}
func (l *logger) Info(msg string, keysAndValues ...interface{}) {
	l.z.Infow(msg, keysAndValues...)
}
func (l *logger) Error(msg string, keysAndValues ...interface{}) {
	l.z.Errorw(msg, keysAndValues...)
}
func (l *logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.z.Fatalw(msg, keysAndValues...)
}

func NewLogger(l *zap.Logger) Logger {
	return &logger{
		z: l.Sugar(),
	}
}
