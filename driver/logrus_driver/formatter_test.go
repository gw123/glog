package logrus_driver

import (
	"testing"

	"github.com/gw123/glog/common"

	"github.com/sirupsen/logrus"
)

func TestTextFormat_Format(t1 *testing.T) {
	logger := logrus.New()
	format := GTextFormat{}
	logger.SetFormatter(format)

	logger.SetReportCaller(true)

	logger.Infof("hello")
	logger.Infof("hello 123")

	entry := logger.WithField(common.KeyPathname, "/index/index").WithField(common.KeyTraceID, "x109123229883")
	entry.Infof("abc")

}
