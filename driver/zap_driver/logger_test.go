package zap_driver

import (
	"testing"

	"github.com/gw123/glog/common"
)

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger(common.Options{})

	if err != nil {
		t.Errorf(err.Error())
	}

	for i := 0; i < 100000; i++ {
		logger.Infof("hello")
	}
}

func BenchmarkTestInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DefaultLogger().Info("hello ")
	}
}

func BenchmarkTestInfof(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DefaultLogger().Infof("hello %d", i)
	}
}
