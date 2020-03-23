package logger

import (
	"github.com/stretchr/testify/mock"
	"testing"
)

const (
	TestLogFormat = "Testing 123456 %s %i"
	TestLogParamS = "string"
	TestLogParamI = 123
)

type MockLogger struct {
	mock.Mock
}

func (l *MockLogger) Debugf(format string, args ...interface{}) {
	args = append([]interface{}{format}, args...)
	l.Called(args...)
}

func (l *MockLogger) Infof(format string, args ...interface{}) {
	args = append([]interface{}{format}, args...)
	l.Called(args...)
}

func (l *MockLogger) Errorf(format string, args ...interface{}) {
	args = append([]interface{}{format}, args...)
	l.Called(args...)
}

func (l *MockLogger) SetLogLevel(lvl string) {
}

func CreateMockLogger() *MockLogger {
	return &MockLogger{}
}

func TestLogger_SetLogAlias(t *testing.T) {

}

func TestDebug(t *testing.T) {
	logger := CreateMockLogger()
	_instance = logger

	logger.On("Debugf", TestLogFormat, TestLogParamS, TestLogParamI).Return(nil)

	Debugf(TestLogFormat, TestLogParamS, TestLogParamI)

	logger.AssertCalled(t, "Debugf", TestLogFormat, TestLogParamS, TestLogParamI)
}

func TestInfo(t *testing.T) {
	logger := CreateMockLogger()
	_instance = logger

	logger.On("Infof", TestLogFormat, TestLogParamS, TestLogParamI).Return(nil)

	Infof(TestLogFormat, TestLogParamS, TestLogParamI)

	logger.AssertCalled(t, "Infof", TestLogFormat, TestLogParamS, TestLogParamI)
}

func TestError(t *testing.T) {
	logger := CreateMockLogger()
	_instance = logger

	logger.On("Errorf", TestLogFormat, TestLogParamS, TestLogParamI).Return(nil)

	Errorf(TestLogFormat, TestLogParamS, TestLogParamI)

	logger.AssertCalled(t, "Errorf", TestLogFormat, TestLogParamS, TestLogParamI)
}
