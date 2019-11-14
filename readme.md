# profile.go使用方式
profile.go 是用来定时记录一个profile文件
默认的时间间隔是一小时 
默认记录的是内存使用情况 (mem)
pprof文件的默认位置在/tmp/pprof/exeName/day_hour_minute/mem.pprof
修改配置方式

```go
    SetProfile(ProfileMode("cpu"), ProfilePath("./"), ProfilePeriod(time.minute*10))
```
##### mode支持的profile类型
    PModeCpu = "cpu"    
	PModeMem = "mem"   
	PModeMutex = "mutex" 
	PModeThread = "thread"
	PModeTrace = "trace"
	
#### 开始记录
StartProfile()

##### 停止记录
StopProfile()

# ctxlog.go 是利用context做日志的上下文记录