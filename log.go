package glog

import (
	"github.com/gw123/glog/common"
	"github.com/gw123/glog/driver/zap_driver"
)

const (
	DriverLogrus = "logrus"
	DriverZap    = "zap"
)

var defaultDriver = DriverLogrus

func init() {

}

// 为了方便创建一个默认的Logger
func DefaultLogger() common.Logger {
	return zap_driver.DefaultLogger()
}

func SetDefaultLoggerDriver(driver string) {
	if driver == DriverLogrus || driver == DriverZap {
		defaultDriver = driver
	}
}

func SetDefaultZapLoggerConfig(options common.Options, withFuncList ...common.WithFunc) error {
	return zap_driver.SetDefaultLoggerConfig(options, withFuncList...)
}

func Error(format string) {
	DefaultLogger().Error(format)
}

func Errorf(format string, other ...interface{}) {
	DefaultLogger().Errorf(format, other...)
}

func Warn(format string) {
	DefaultLogger().Warn(format)
}

func Warnf(format string, other ...interface{}) {
	DefaultLogger().Warnf(format, other...)
}

func Info(format string) {
	DefaultLogger().Info(format)
}

func Infof(format string, other ...interface{}) {
	DefaultLogger().Infof(format, other...)
}

func Debug(format string) {
	DefaultLogger().Debug(format)
}

func Debugf(format string, other ...interface{}) {
	DefaultLogger().Debugf(format, other...)
}

func WithField(format string, other ...interface{}) common.Logger {
	return DefaultLogger().WithField(format, other)
}

func WithErr(err error) common.Logger {
	return DefaultLogger().WithError(err)
}
