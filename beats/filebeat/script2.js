// 解析message格式消息到json格式
// [2020-12-28 19:28:42.325] [info] [] cluster/service_manager.go 190 [] - ServiceManager[businesspoint].
function process(event) {
    var message = event.Get("message")
    var rule = /\[(.*?)\] \[(.*?)\] \[(.*?)\] (.+) (\d{1,4}) \[(.*?)\] - (.*)/g
    var res = rule.exec(message)
    if(res){
        event.Put("time",res[1])
        event.Put("level", res[2])
        event.Put("line", res[4] + " "+res[5])
        event.Put("message", res[7])
    }
}