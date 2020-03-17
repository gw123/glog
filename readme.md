# ctxlog.go 是利用context做日志的上下文记录


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
