package glog

import (
	"github.com/gw123/glog/common"
	"github.com/gw123/glog/zap"
)

func DefaultLogger() common.Logger {
	return zap.DefaultLogger()
}

func Log() common.Logger {
	return zap.DefaultLogger()
}

func SetDefaultLoggerConfig(options common.Options, withFuncList ...common.WithFunc) error {
	return zap.SetDefaultLoggerConfig(options, withFuncList...)
}

func Error(format string) {
	zap.GetInnerLogger().Error(format)
}

func Errorf(format string, other ...interface{}) {
	zap.GetInnerLogger().Errorf(format, other...)
}

func Warn(format string) {
	zap.GetInnerLogger().Warn(format)
}

func Warnf(format string, other ...interface{}) {
	zap.GetInnerLogger().Warnf(format, other...)
}

func Info(format string) {
	zap.GetInnerLogger().Info(format)
}

func Infof(format string, other ...interface{}) {
	zap.GetInnerLogger().Infof(format, other...)
}

func Debug(format string) {
	zap.GetInnerLogger().Debug(format)
}

func Debugf(format string, other ...interface{}) {
	zap.GetInnerLogger().Debugf(format, other...)
}

func WithField(format string, other ...interface{}) common.Logger {
	return zap.GetInnerLogger().WithField(format, other)
}

func WithErr(err error) common.Logger {
	return zap.GetInnerLogger().WithError(err)
}

func WithError(err error) common.Logger {
	return zap.GetInnerLogger().WithError(err)
}
