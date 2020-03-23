package logger

import (
	"github.com/kintohub/common-go/logger/config"
)

const (
	LogLevelPanic   = "panic"
	LogLevelError   = "error"
	LogLevelInfo    = "info"
	LogLevelWarning = "warning"
	LogLevelDebug   = "debug"
	LogLevelTrace   = "trace"
)

const DefaultLogLevel = LogLevelDebug

type LogLevel int32

var (
	_instance ILogger
)

type ILogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	SetLogLevel(lvl string)
}

/**
  Instance method for logger
*/

func SetLogger(logger ILogger) {
	_instance = logger
}

//gets singleton of logger
func GetLogger() ILogger {
	if _instance == nil {
		if config.DefaultLogLevel == LogLevelDebug {
			_instance = NewReflectLogger()
		} else {
			_instance = NewSimpleLogger()

		}
	}
	return _instance
}

/**
  Static methods for logging with default logger
*/

func SetLogLevel(lvl string) {
	GetLogger().SetLogLevel(lvl)
}

func Error(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

func Debug(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}
