package glog

import (
	"sync"

	"github.com/gw123/glog/common"
	"github.com/gw123/glog/driver/logrus_driver"
)

var defaultLogger common.Logger
var once = sync.Once{}

// 为了方便创建一个默认的Logger
func DefaultLogger() common.Logger {
	once.Do(func() {
		defaultLogger = logrus_driver.DefaultLogger()
	})

	return defaultLogger
}

func Error(format string) {
	defaultLogger.Error(format)
}

func Errorf(format string, other ...interface{}) {
	defaultLogger.Errorf(format, other)
}

func Warn(format string) {
	defaultLogger.Warn(format)
}

func Warnf(format string, other ...interface{}) {
	defaultLogger.Warnf(format, other)
}

func Info(format string) {
	defaultLogger.Info(format)
}

func Infof(format string, other ...interface{}) {
	defaultLogger.Infof(format, other)
}

func Debug(format string) {
	defaultLogger.Debug(format)
}

func Debugf(format string, other ...interface{}) {
	defaultLogger.Debugf(format, other)
}
