package glog

import (
	"context"
	"testing"
)

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
	entry = ExtractEntry(ctx).WithField("ip2", "10.0.0.1")
	for i := 0; i < 20; i++ {
		entry.WithField("key", i).Infof("TestContent abc %d", i)
	}
	//输出结果
	//{"RequestID":"10000001","app_name":"web","ip":"10.0.0.1","level":"info","msg":"TestContent","time":"2020-03-17 20:34:14"}
}
