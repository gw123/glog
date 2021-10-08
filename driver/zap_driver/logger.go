package zap_driver

import (
	"sync"
	"time"

	"github.com/gw123/glog/common"

	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05.000"
)

type Logger struct {
	*zap.SugaredLogger
}

func (l Logger) WithField(key string, value interface{}) common.Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(zap.Any(key, value)),
	}
}

func (l Logger) WithFields(fields map[string]interface{}) common.Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(zap.Any("fields", fields)),
	}
}

func (l Logger) WithError(err error) common.Logger {
	return &Logger{
		l.SugaredLogger.With(zap.Error(err)),
	}
}

func (l Logger) Warningf(format string, args ...interface{}) {
	l.SugaredLogger.Warnf(format, args...)
}

func (l Logger) Warning(args ...interface{}) {
	l.SugaredLogger.Warn(args...)
}

var defaultLogger *Logger
var loggerOnce sync.Once

func SetDefaultLoggerConfig(options common.Options, withFuncList ...common.WithFunc) error {
	var err error
	newLogger, err := NewLogger(options, withFuncList...)
	if err != nil {
		return err
	}
	defaultLogger = newLogger
	return err
}

func init() {
	var err error
	option := common.Options{}
	defaultLogger, err = NewLogger(option, common.WithConsoleEncoding(), common.WithLevel(common.DebugLevel), common.WithStdoutOutputPath(), common.WithStderrErrorOutputPath())
	if err != nil {
		panic(err)
	}
}

func DefaultLogger() *Logger {
	return defaultLogger
}

func NewLogger(options common.Options, withFuncs ...common.WithFunc) (*Logger, error) {
	for _, withFunc := range withFuncs {
		withFunc(&options)
	}

	if len(options.OutputPaths) == 0 {
		options.OutputPaths = append(options.OutputPaths, common.PathStdout)
	}

	if options.Encoding == "" {
		options.Encoding = common.EncodeConsole
	}

	encodeCfg := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "defaultLogger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel: func(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + lv.String() + "]")
		},
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[v1] [" + t.Format(DateTimeFormat) + "]")
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: func(call zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			trPath := call.TrimmedPath()
			enc.AppendString("[] " + trPath + " [] ")
		},

		ConsoleSeparator: " ",
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.Level(options.Level)),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		DisableStacktrace: true,
		Encoding:          options.Encoding,
		EncoderConfig:     encodeCfg,
		OutputPaths:       options.OutputPaths,
		ErrorOutputPaths:  options.ErrorOutputPaths,
	}

	logger, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(0))
	if err != nil {
		return nil, err
	}

	su := logger.Sugar()
	return &Logger{SugaredLogger: su}, nil
}
