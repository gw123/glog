package glog

import (
	"errors"
	"testing"

	"github.com/gw123/glog/common"
)

func TestError(t *testing.T) {
	type args struct {
		format string
		params []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "abc",
			args: args{
				format: "this is a log %s",
				params: []string{"test"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args.format)
			Infof(tt.args.format, tt.args.params[0])
			Warn(tt.args.format)
			Warnf(tt.args.format, tt.args.params[0])
			Error(tt.args.format)
			Errorf(tt.args.format, tt.args.params[0])
			Debug(tt.args.format)
			Debugf(tt.args.format, tt.args.params[0])

			Info(tt.args.format)
			Infof(tt.args.format, tt.args.params[0])
			Warn(tt.args.format)
			Warnf(tt.args.format, tt.args.params[0])
			Error(tt.args.format)
			Errorf(tt.args.format, tt.args.params[0])
			Debug(tt.args.format)
			Debugf(tt.args.format, tt.args.params[0])
		})
	}
}

func TestWithError(t *testing.T) {
	var err = errors.New("xxxxx")
	DefaultLogger().WithError(err).Error("xxxxxxxx")
}

func TestSetDefaultLoggerDriver(t *testing.T) {
	DefaultLogger().WithField("abc", "hello").Info("show log")
}

func TestSetDefaultZapLoggerConfig(t *testing.T) {
	SetDefaultZapLoggerConfig(common.Options{
		OutputPaths:      []string{common.PathStdout, "test.log"},
		ErrorOutputPaths: []string{common.PathStderr},
		Encoding:         common.EncodeConsole,
		Level:            common.DebugLevel,
	})
	DefaultLogger().WithField("abc", "hello").Debug("show log")
	DefaultLogger().WithField("abc", "hello").Info("show log")
	DefaultLogger().Named("glog").WithField("abc", "hello").Info("show log")
}

func TestSetDefaultZapLoggerConfig2(t *testing.T) {
	err := SetDefaultZapLoggerConfig(common.Options{
		OutputPaths:      []string{common.PathStdout, "./test.log"},
		ErrorOutputPaths: []string{common.PathStderr},
		Encoding:         common.EncodeConsole,
		Level:            common.DebugLevel,
	})
	if err != nil {
		t.Error(err)
	}
	DefaultLogger().WithField("abc", "hello").Debug("show log")
	DefaultLogger().WithField("abc", "hello").Info("show log")
}

func BenchmarkTestInfof(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DefaultLogger().Infof("hello %d", i)
	}
}

func BenchmarkTestWithPath(b *testing.B) {
	SetDefaultZapLoggerConfig(common.Options{
		OutputPaths:      []string{common.PathStdout, "test.log"},
		ErrorOutputPaths: []string{common.PathStderr},
		Encoding:         common.EncodeConsole,
		Level:            common.DebugLevel,
	})

	for i := 0; i < b.N; i++ {
		DefaultLogger().WithField("abc", "hello").Info("show log")
	}
}
