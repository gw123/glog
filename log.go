package glog

import (
	"runtime"

	"github.com/gw123/glog/hook"
	"github.com/sirupsen/logrus"
)

var defaultLogger *logrus.Logger
var isDebug = false

func SetDebug(flag bool) {
	isDebug = flag
}

// 为了方便创建一个默认的Logger
func DefaultLogger() *logrus.Logger {
	if defaultLogger == nil {
		defaultLogger = TextLogger()
	}
	return defaultLogger
}

func JsonLogger() *logrus.Logger {
	jsonLogger := logrus.New()
	jsonLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:  TimeFormat,
		DisableTimestamp: false,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, frame.File
		},
		PrettyPrint: isDebug,
	})
	jsonLogger.AddHook(&hook.LogHook{Field: "caller"})

	return jsonLogger
}

func GetDefaultJsonLoggerFormatter() logrus.Formatter {
	return &logrus.JSONFormatter{
		TimestampFormat:  TimeFormat,
		DisableTimestamp: false,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, frame.File
		},
		PrettyPrint: isDebug,
	}
}

func JsonEntry() *logrus.Entry {
	return logrus.NewEntry(JsonLogger())
}

// 为了方便创建一个默认的Logger
func TextLogger() *logrus.Logger {
	textLogger := logrus.New()
	textLogger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  TimeFormat,
		DisableTimestamp: false,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, frame.File
		},
	})
	return textLogger
}

func TextEntry() *logrus.Entry {
	return logrus.NewEntry(TextLogger())
}

func SetDefaultLoggerFormatter(formatter logrus.Formatter) {
	newLogger := logrus.New()
	newLogger.SetFormatter(formatter)
	newLogger.AddHook(&hook.LogHook{Field: "caller"})
	defaultLogger = newLogger
}

func Error(format string) {
	defaultLogger.Error(format)
}

func Errorf(format string, other ...interface{}) {
	defaultLogger.Errorf(format, other...)
}

func Warn(format string) {
	defaultLogger.Warn(format)
}

func Warnf(format string, other ...interface{}) {
	defaultLogger.Warnf(format, other...)
}

func Info(format string) {
	defaultLogger.Info(format)
}

func Infof(format string, other ...interface{}) {
	defaultLogger.Infof(format, other...)
}

func Debug(format string) {
	defaultLogger.Debug(format)
}

func Debugf(format string, other ...interface{}) {
	defaultLogger.Debugf(format, other...)
}
