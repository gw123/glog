package glog

import (
	"github.com/gw123/glog/common"
	"github.com/gw123/glog/driver/zap_driver"
)

const (
	DriverZap = "zap"
)

// 为了方便创建一个默认的Logger
func DefaultLogger() common.Logger {
	return zap_driver.DefaultLogger()
}

// 为了方便创建一个默认的Logger
func Log() common.Logger {
	return zap_driver.DefaultLogger()
}

func SetDefaultZapLoggerConfig(options common.Options, withFuncList ...common.WithFunc) error {
	return zap_driver.SetDefaultLoggerConfig(options, withFuncList...)
}

func Error(format string) {
	zap_driver.GetInnerLogger().Error(format)
}

func Errorf(format string, other ...interface{}) {
	zap_driver.GetInnerLogger().Errorf(format, other...)
}

func Warn(format string) {
	zap_driver.GetInnerLogger().Warn(format)
}

func Warnf(format string, other ...interface{}) {
	zap_driver.GetInnerLogger().Warnf(format, other...)
}

func Info(format string) {
	zap_driver.GetInnerLogger().Info(format)
}

func Infof(format string, other ...interface{}) {
	zap_driver.GetInnerLogger().Infof(format, other...)
}

func Debug(format string) {
	zap_driver.GetInnerLogger().Debug(format)
}

func Debugf(format string, other ...interface{}) {
	zap_driver.GetInnerLogger().Debugf(format, other...)
}

func WithField(format string, other ...interface{}) common.Logger {
	return zap_driver.GetInnerLogger().WithField(format, other)
}

func WithErr(err error) common.Logger {
	return zap_driver.GetInnerLogger().WithError(err)
}
