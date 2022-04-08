package main

import (
	"context"
	"github.com/gw123/glog"
)

func main() {
	glog.Infof("show log level %s", "other a")
	glog.Errorf("show log error %s", "other b")
	glog.Warnf("show log warn %s", "other c")
	glog.Debugf("show log debug %s", "other d")

	glog.Log().Info("show log")
	glog.Log().Debugf("show log debug %s", "other")
	glog.Log().Infof("show log info %s", "other")
	glog.Log().Warnf("show log warn %s", "other")
	glog.Log().Debugf("show log debug %s", "other")

	glog.Log().WithField("key", "val").Infof("show log info")

	glog.WithOTEL(context.Background()).Info("with otel id")
}
