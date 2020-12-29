// 解析message格式消息到json格式
// [2020-12-28 22:10:58] [info] [glog/ctxlog_test.go:67] [/home/index] [10000001] TestContent abc 0 {"app":"xytschol","ip":"10.0.0.1"}
// input  2020-12-28 19:28:42.325
// output 2020-12-29T18:15:23.847+0800
function formatTime(time){
    time = time.replace(/-/g,':').replace('.', ':' ).replace(' ',':')
    time = time.split(':')
    if(time.length != 6){
        var date = new Date()
        var res = date.getFullYear() +"-" + (date.getMonth() + 1)
            + "-" + date.getDate() + "T" + date.getHours() + ":" +
            date.getMinutes() +":" + date.getSeconds()+"." + date.getMilliseconds()+"+0800"
        return res
    }
    var time1  = time[0] +"-" +time[1] +"-"+time[2] +"T"+time[3] +":"+time[4] +":"+time[5]+"."+time[6]+"+0800"
    return time1
}

function process(event) {
    var message = event.Get("message")
    var rule = /\[(.*?)\] \[(.*?)\] \[(.*?)\] \[(.*?)\] \[(.*?)\] (.*)/g
    var res = rule.exec(message)
    if(res){
        event.Put("timestamp", formatTime(res[1]))
        event.Put("time",res[1])
        event.Put("level", res[2])
        event.Put("line", res[3])
        event.Put("pathname", res[4])
        event.Put("trace_id", res[5])
        event.Put("message", res[6])
    }
}