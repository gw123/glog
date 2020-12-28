// 解析message格式消息到json格式
// [2020-12-28 22:10:58] [info] [glog/ctxlog_test.go:67] [/home/index] [10000001] TestContent abc 0 {"app":"xytschol","ip":"10.0.0.1"}
function process(event) {
    var message = event.Get("message")
    var rule = /\[(.*?)\] \[(.*?)\] \[(.*?)\] \[(.*?)\] \[(.*?)\] (.*)/g
    var res = rule.exec(message)
    if(res){
        event.Put("time",res[1])
        event.Put("level", res[2])
        event.Put("line", res[3])
        event.Put("pathname", res[4])
        event.Put("trace_id", res[5])
        event.Put("message", res[6])
    }
}