package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type SimpleLogger struct {
	logger logrus.Logger
}

func NewSimpleLogger() *SimpleLogger {
	instance := SimpleLogger{}
	instance.logger.SetFormatter(&logrus.JSONFormatter{})
	instance.logger.SetOutput(os.Stdout)
	instance.SetLogLevel(DefaultLogLevel)
	return &instance
}

func (l *SimpleLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *SimpleLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *SimpleLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *SimpleLogger) SetLogLevel(lvl string) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		panic(err)
	}
	l.logger.SetLevel(level)
}

