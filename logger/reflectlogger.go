package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type ReflectLogger struct {
	logger logrus.Logger
	ShouldShowFileName bool
}

func NewReflectLogger() *ReflectLogger {
	instance := ReflectLogger{}
	instance.logger.SetFormatter(&logrus.JSONFormatter{})
	instance.logger.SetOutput(os.Stdout)
	instance.SetLogLevel(DefaultLogLevel)
	instance.ShouldShowFileName = false
	return &instance
}

func (l *ReflectLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(l.getCallerName()+format, args...)
}

func (l *ReflectLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(l.getCallerName()+format, args...)
}

func (l *ReflectLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(l.getCallerName()+format, args...)
}

func (l *ReflectLogger) SetLogLevel(lvl string) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		panic(err)
	}
	l.logger.SetLevel(level)
}


func (l *ReflectLogger) getCallerName() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(4, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	funcName := getFuncName(frame.Function)
	prefix := ""
	if l.ShouldShowFileName {
		return fmt.Sprintf("[%s][%s:%d]", funcName, frame.File, frame.Line)
	} else {
		return fmt.Sprintf("[%s]", funcName)
	}

	// TODO: go vet prompts this line of code is unreachable
	return prefix
}

func getFuncName(fullName string) string {
	paths := strings.Split(fullName, "/")
	fullFuncName := paths[len(paths)-1]
	components := strings.Split(fullFuncName, ".")
	return components[0] + "." + components[len(components)-1]
}
