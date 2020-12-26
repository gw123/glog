package logrus_driver

import (
	"runtime"

	"github.com/gw123/glog/common"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Entry
}

func NewLogger(logger *logrus.Logger) *Logger {
	return &Logger{
		Entry: logrus.NewEntry(logger),
	}
}

func DefaultLogger() *Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(GTextFormat{})
	return &Logger{
		Entry: logrus.NewEntry(logger),
	}
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
	l.Entry = l.Entry.WithField(key, value)
	return l
}

func (l *Logger) WithFields(fields map[string]interface{}) common.Logger {
	l.Entry = l.Entry.WithFields(fields)
	return l
}

func (l *Logger) WithError(err error) common.Logger {
	l.Entry = l.Entry.WithError(err)
	return l
}
