package logger

import (
	"fmt"
	"github.com/kintohub/common-go/logger/config"
	"github.com/kintohub/common-go/logger/constants"
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
	instance.SetLogLevel(config.DefaultLogLevel)
	instance.SetLogAlias(constants.DefaultLogAlias)
	instance.Info("Alias logger initialized with LogLevel: %s", config.DefaultLogLevel)

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

func (l *AliasLogger) SetLogAlias(logAlias string) {

	if logAlias == "" {
		logAlias = constants.DefaultLogAlias
	}

	l.LogAlias = fmt.Sprintf("[%s]", logAlias)
}
