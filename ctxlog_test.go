package glog

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestDefaultLogger(t *testing.T) {
	log := DefaultLogger()
	log.Info("hello world")
	log.WithError(errors.New("test err")).Info("hello world")
}

func TestContext(t *testing.T) {
	entry := DefaultLogger()
	ctx := context.Background()
	ctx = ToContext(ctx, entry)
	AddFields(ctx, logrus.Fields{
		"app_name": "web",
	})
	ExtractEntry(ctx).WithFields(logrus.Fields{"field1": "hello world"}).Info("TestContent")
}

func TestExtractEntryWithID(t *testing.T) {
	ctx := context.Background()

	AddTraceID(ctx, "10000001")
	entry := DefaultLogger()
	ctx = ToContext(ctx, entry)
	AddFields(ctx, logrus.Fields{
		"app_name": "web",
	})
	ExtractEntry(ctx).WithFields(logrus.Fields{"field1": "hello world"}).Info("TestContent")
	ExtractEntry(ctx).WithFields(logrus.Fields{"field1": "cpu"}).Info("very nice")
}

func TestExtractEntryWithID2(t *testing.T) {
	ctx := context.Background()
	entry := DefaultLogger()
	ctx = ToContext(ctx, entry)
	AddTraceID(ctx, "10000001")
	AddField(ctx, "app_name", "web")
	ExtractEntry(ctx).WithFields(logrus.Fields{"field1": "hello world"}).Info("TestContent")
	ExtractEntry(ctx).WithFields(logrus.Fields{"field1": "cpu"}).Info("very nice")
}

func TestALl(t *testing.T) {
	//1. 一般是在应用的入口创建一个根context
	ctx := context.Background()
	//2. 新建立entry 使用ToContext将entry传入context
	entry := DefaultLogger()
	//3. 加入一些应用全局的描述
	entry.WithField("app", "xytschol")
	//4. 然后将这个根context传到应用框架中

	//5. 在中间件里调用AddRequestID(ctx) 记录一次请求的同一个 requestID
	ctx = ToContext(ctx, entry)
	AddTraceID(ctx, "10000001")
	AddPathname(ctx, "/home/index")

	//6. 在action 或者service 等地方记录日志记录日志
	entry = ExtractEntry(ctx).WithField("ip", "10.0.0.1")
	for i := 0; i < 20; i++ {
		entry.Infof("TestContent abc %d", i)
	}
	//输出结果
	//{"RequestID":"10000001","app_name":"web","ip":"10.0.0.1","level":"info","msg":"TestContent","time":"2020-03-17 20:34:14"}
}
