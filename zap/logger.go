package zap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(args...),
	}
}

func (l Logger) WithError(err error) common.Logger {
	if err == nil {
		return l
	}

	return &Logger{
		l.SugaredLogger.With(zap.String("error", err.Error())),
	}
}

func (l Logger) Warningf(format string, args ...interface{}) {
	l.SugaredLogger.Warnf(format, args...)
}

func (l Logger) Warning(args ...interface{}) {
	l.SugaredLogger.Warn(args...)
}

func (l Logger) Named(name string) common.Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.Named(name),
	}
}

var (
	defaultLogger *Logger
	innerLogger   *Logger
	loggerMutex   sync.RWMutex
)

func SetDefaultLoggerConfig(options common.Options, withFuncList ...common.WithFunc) error {
	newLogger, err := NewLogger(options, withFuncList...)
	if err != nil {
		return err
	}

	innerOptions := options
	innerWithFuncs := make([]common.WithFunc, len(withFuncList)+1)
	copy(innerWithFuncs, withFuncList)
	innerWithFuncs[len(withFuncList)] = common.WithCallerSkip(1)

	newInnerLogger, err := NewLogger(innerOptions, innerWithFuncs...)
	if err != nil {
		return err
	}

	loggerMutex.Lock()
	defaultLogger = newLogger
	innerLogger = newInnerLogger
	loggerMutex.Unlock()
	return nil
}

func init() {
	// Register custom console encoder globally
	err := zap.RegisterEncoder("custom-console", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return newCustomConsoleEncoder(cfg), nil
	})
	if err != nil {
		fmt.Printf("WARNING: Failed to register custom console encoder: %v\n", err)
		panic(err)
	}

	option := common.Options{}
	defaultLogger, err = NewLogger(option, common.WithConsoleEncoding(), common.WithLevel(common.InfoLevel), common.WithStdoutOutputPath(), common.WithStderrErrorOutputPath())
	if err != nil {
		panic(err)
	}

	innerLogger, err = NewLogger(option, common.WithCallerSkip(1), common.WithConsoleEncoding(), common.WithLevel(common.InfoLevel), common.WithStdoutOutputPath(), common.WithStderrErrorOutputPath())
	if err != nil {
		panic(err)
	}
}

func DefaultLogger() *Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	return defaultLogger
}

func GetInnerLogger() *Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	return innerLogger
}

func NewLogger(options common.Options, withFuncs ...common.WithFunc) (*Logger, error) {
	for _, withFunc := range withFuncs {
		withFunc(&options)
	}

	if len(options.OutputPaths) == 0 {
		options.OutputPaths = append(options.OutputPaths, common.PathStdout)
	}

	if options.Encoding == "" {
		options.Encoding = "custom-console"
	} else if options.Encoding == common.EncodeConsole {
		options.Encoding = "custom-console"
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
			enc.AppendString("[" + t.Format(DateTimeFormat) + "]")
		},
		//EncodeDuration: func(duration time.Duration, encoder zapcore.PrimitiveArrayEncoder) {
		//	encoder.AppendString(duration.String())
		//},
		EncodeCaller: func(call zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			trPath := call.TrimmedPath()
			enc.AppendString(trPath)
		},
		EncodeName: func(name string, enc zapcore.PrimitiveArrayEncoder) {
			names := strings.Split(name, ".")
			if len(names) == 1 {
				enc.AppendString("[]")
				return
			}
			enc.AppendString("[" + strings.Join(names[1:], ".") + "]")
		},
		ConsoleSeparator: " ",
	}

	allLogPath := append(options.OutputPaths, options.ErrorOutputPaths...)
	for _, path := range allLogPath {
		if path == common.PathStderr || path == common.PathStdout {
			continue
		}

		dir := filepath.Dir(path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0760); err != nil {
				return nil, err
			}
		}
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

	logger, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(options.CallerSkip))
	if err != nil {
		return nil, err
	}

	su := logger.Sugar().Named("-")
	return &Logger{SugaredLogger: su}, nil
}
