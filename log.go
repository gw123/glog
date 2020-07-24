package glog

import (
	"github.com/gw123/glog/hook"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	formatter := logrus.TextFormatter{
		ForceColors:               false,
		DisableColors:             false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	}
	SetDefaultLoggerFormatter(&formatter)
}

func GetLogger() *logrus.Logger {
	return logger
}

func SetDefaultLogger(l *logrus.Logger) {
	logger = l
}

func SetDefaultJsonLogger() {
	formatter := logrus.JSONFormatter{
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02 15:04",
		FieldMap:         nil,
		CallerPrettyfier: nil,
	}
	SetDefaultLoggerFormatter(&formatter)
}

func SetDefaultLoggerFormatter(formatter logrus.Formatter) {
	newLogger := logrus.New()
	newLogger.SetFormatter(formatter)
	newLogger.AddHook(&hook.LogHook{Field: "caller"})
	logger = newLogger
}

func Error(format string) {
	logger.Error(format)
}

func Errorf(format string, other ...interface{}) {
	logger.Errorf(format, other...)
}

func Warn(format string) {
	logger.Warn(format)
}

func Warnf(format string, other ...interface{}) {
	logger.Warnf(format, other...)
}

func Info(format string) {
	logger.Info(format)
}

func Infof(format string, other ...interface{}) {
	logger.Infof(format, other...)
}

func Debug(format string) {
	logger.Debug(format)
}

func Debugf(format string, other ...interface{}) {
	logger.Debugf(format, other...)
}
