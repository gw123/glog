package logrus_driver

import (
	"runtime"
	"sync"

	"github.com/gw123/glog/common"
	"github.com/sirupsen/logrus"
)

const ErrorKey = "error"

type Logger struct {
	*logrus.Entry
}

func NewLogger(logger *logrus.Logger) *Logger {
	return &Logger{
		Entry: logrus.NewEntry(logger),
	}
}

var logger *Logger
var loggerOnce sync.Once

func init() {
	_logger := logrus.New()
	_logger.SetReportCaller(true)
	_logger.SetLevel(logrus.DebugLevel)
	_logger.SetFormatter(GTextFormat{})
	logger = NewLogger(_logger)
}

func DefaultLogger() *Logger {
	return logger
}

func DefaultJsonLogger() *Logger {
	jsonLogger := logrus.New()
	jsonLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:  common.TimeFormat,
		DisableTimestamp: false,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, frame.File
		},
	})
	return NewLogger(jsonLogger)
}

func (l *Logger) WithField(key string, value interface{}) common.Logger {
	tmpLoger := &Logger{
		Entry: l.Entry.WithField(key, value),
	}
	return tmpLoger
}

func (l *Logger) WithFields(fields map[string]interface{}) common.Logger {
	tmpLoger := &Logger{
		Entry: l.Entry.WithFields(fields),
	}
	return tmpLoger
}

func (l *Logger) WithError(err error) common.Logger {
	tmpLoger := &Logger{
		Entry: l.Entry.WithField(ErrorKey, err.Error()),
	}
	return tmpLoger
}
