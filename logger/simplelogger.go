package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type SimpleLogger struct {
	logrus.Logger
}

func NewSimpleLogger() *SimpleLogger {
	instance := SimpleLogger{}
	instance.SetFormatter(&logrus.JSONFormatter{})
	instance.SetOutput(os.Stdout)
	instance.SetLogLevel(DefaultLogLevel)
	return &instance
}

func (l *SimpleLogger) Errorf(format string, args ...interface{}) {
	l.Errorf(format, args...)
}

func (l *SimpleLogger) Debugf(format string, args ...interface{}) {
	l.Debugf(format, args...)
}

func (l *SimpleLogger) Infof(format string, args ...interface{}) {
	l.Infof(format, args...)
}

func (l *SimpleLogger) SetLogLevel(lvl string) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		panic(err)
	}
	l.SetLevel(level)
}

