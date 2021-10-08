package zap_driver

import (
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLog(t *testing.T) {
	url := "http://xxxx.html"
	logger, _ := zap.NewProduction()
	//defer defaultLogger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "url",
		"backoff", time.Second,
		"attempt", 3,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	sugar.With(zap.Any("response", 12)).Info("test any")
}

func TestTextLog(t *testing.T) {

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
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: func(call zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			trPath := strings.Replace(call.TrimmedPath(), ":", " ", -1)
			enc.AppendString("[] " + trPath + " [] -")
		},
	}
	encodeCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + t.Format(DateTimeFormat) + "]")
	}

	encodeCfg.EncodeLevel = func(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("[" + lv.String() + "]")
	}

	encodeCfg.EncodeCaller = func(call zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		trPath := strings.Replace(call.TrimmedPath(), ":", " ", -1)
		enc.AppendString("[] " + trPath + " [] -")
	}
	encodeCfg.ConsoleSeparator = " "

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    encodeCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(0))
	if err != nil {
		t.Error(err)
	}

	su := logger.Sugar()

	su.With(zap.Any("xxx", 123)).Info("hello world")
}
