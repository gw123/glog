# 基础使用基于zap做的封装,默认记录时间和代码位置
```go
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
```
### 输出
```
[2022-04-01 14:12:15.891] [info] [] demo/demo.go:8 []  show log level other a
[2022-04-01 14:12:15.892] [error] [] demo/demo.go:9 []  show log error other b
[2022-04-01 14:12:15.892] [warn] [] demo/demo.go:10 []  show log warn other c
[2022-04-01 14:12:15.892] [debug] [] demo/demo.go:11 []  show log debug other d
[2022-04-01 14:12:15.892] [info] [] demo/demo.go:13 []  show log
[2022-04-01 14:12:15.892] [debug] [] demo/demo.go:14 []  show log debug other
[2022-04-01 14:12:15.892] [info] [] demo/demo.go:15 []  show log info other
[2022-04-01 14:12:15.892] [warn] [] demo/demo.go:16 []  show log warn other
[2022-04-01 14:12:15.892] [debug] [] demo/demo.go:17 []  show log debug other
[2022-04-01 14:12:15.892] [info] [] demo/demo.go:19 []  show log info {"key": "val"}
```

# ctxlog.go 是利用context做日志的上下文记录
## 使用步骤
``` 
//1. 一般是在应用的入口创建一个根context
ctx := context.Background()
//2. 新建立entry 使用ToContext将entry传入context
entry := NewDefaultEntry()
ctx = ToContext(ctx, entry)
//3. 加入一些应用全局的描述
AddField(ctx, "app_name", "web")
//4. 然后将这个根context传到应用框架中

//5. 在中间件里调用AddRequestID(ctx, requestID) 记录一次请求的requestID, 
//后面在调用ExtractEntry()输出日志时候会把ctx里面记录的requestID输出
ctx = AddRequestID(ctx, "10000001")

//6. 在action 或者service 等地方记录日志记录日志
ExtractEntry(ctx).WithField("ip", "10.0.0.1").Info("TestContent")

//输出结果
//{"RequestID":"10000001","app_name":"web","ip":"10.0.0.1","level":"info","msg":"TestContent","time":"2020-03-17 20:34:14"}
```