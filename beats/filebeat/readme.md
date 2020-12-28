## 启动filebeat
filebeat -e  --path.config .


## 日志采集结果
``` json
{"@timestamp":"2020-12-28T14:58:09.177Z","@metadata":{"beat":"filebeat","type":"_doc","version":"7.10.1"},"time":"2020-12-28 22:10:58","message":"TestContent abc 20 {\"app\":\"xytschol\",\"ip\":\"10.0.0.1\"}","level":"info","line":"glog/ctxlog_test.go:67","pathname":"/home/index","trace_id":"10000001"}
{"@timestamp":"2020-12-28T14:58:09.177Z","@metadata":{"beat":"filebeat","type":"_doc","version":"7.10.1"},"time":"2020-12-28 22:10:58","message":"TestContent abc 20 {\"app\":\"xytschol\",\"ip\":\"10.0.0.1\"}","level":"info","line":"glog/ctxlog_test.go:67","pathname":"/home/index","trace_id":"10000001"}
```