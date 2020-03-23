package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

type AliasLogger struct {
	logrus.Logger
	LogAlias string
}

func NewAliasLogger() *AliasLogger {
	instance := AliasLogger{}
	instance.SetFormatter(&logrus.JSONFormatter{})
	instance.SetOutput(os.Stdout)
	instance.SetLogLevel(DefaultLogLevel)
	instance.Info("Alias logger initialized with LogLevel: %s", DefaultLogLevel)

	return &instance
}

func (l *AliasLogger) Error(format string, args ...interface{}) {
	l.Errorf(l.LogAlias+format, args...)
}

func (l *AliasLogger) Debug(format string, args ...interface{}) {
	l.Debugf(l.LogAlias+format, args...)
}

func (l *AliasLogger) Info(format string, args ...interface{}) {
	l.Infof(l.LogAlias+format, args...)
}

func (l *AliasLogger) SetLogLevel(lvl string) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		panic(err)
	}
	l.SetLevel(level)
}

