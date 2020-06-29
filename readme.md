# 基础使用基于 logrus 做了简单的封装,默认记录时间和代码位置
```go
    Info(tt.args.format)
    Infof(tt.args.format, tt.args.params[0])
	Warn(tt.args.format)
	Warnf(tt.args.format, tt.args.params[0])
	Error(tt.args.format)
	Errorf(tt.args.format, tt.args.params[0])
	Debug(tt.args.format)
	Debugf(tt.args.format, tt.args.params[0])

    //切换默认的日志格式为json
    SetDefaultJsonLogger()
	Info(tt.args.format)
	Infof(tt.args.format, tt.args.params[0])
	Warn(tt.args.format)
	Warnf(tt.args.format, tt.args.params[0])
	Error(tt.args.format)
	Errorf(tt.args.format, tt.args.params[0])
	Debug(tt.args.format)
	Debugf(tt.args.format, tt.args.params[0])
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

//5. 在中间件里调用AddRequestId(ctx, requestId) 记录一次请求的requestId, 
//后面在调用ExtractEntry()输出日志时候会把ctx里面记录的requestId输出
ctx = AddRequestId(ctx, "10000001")

//6. 在action 或者service 等地方记录日志记录日志
ExtractEntry(ctx).WithField("ip", "10.0.0.1").Info("TestContent")

//输出结果
//{"RequestId":"10000001","app_name":"web","ip":"10.0.0.1","level":"info","msg":"TestContent","time":"2020-03-17 20:34:14"}
```

# 记录日志
# profile.go使用方式
profile.go 是用来定时记录一个profile文件
默认的时间间隔是一小时 
默认记录的是内存使用情况 (mem)
pprof文件的默认位置在/tmp/pprof/exeName/day_hour_minute/mem.pprof
修改配置方式

tip:go tool pprof 和 pprof 工具使用请查看:
http://www.xytschool.com/chapter/454.html

### 修改profile类型
    SetProfile(ProfileMode(mode))
    mode支持的profile类型
    PModeCpu = "cpu"    
	PModeMem = "mem"   
	PModeMutex = "mutex" 
	PModeThread = "thread"
	PModeTrace = "trace"
	
## 修改pprof文件生成路径
	SetProfile(ProfilePath("./"))
## 修改单个pprof采集时长
	SetProfile(ProfilePeriod(time.minute*10))
## 多个配置一起修改
    SetProfile(ProfileMode("cpu"), ProfilePath("./"), ProfilePeriod(time.minute*10))
#### 开始记录
    StartProfile()
##### 停止记录
    StopProfile()
