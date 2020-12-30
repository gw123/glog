// 解析message格式消息到json格式
// [2020-12-28 22:10:58.112] [info] [glog/ctxlog_test.go:67] [/home/index] [10000001] TestContent abc 0 {"app":"xytschol","ip":"10.0.0.1"}
// input  2020-12-28 19:28:42.325
// output 2020-12-29T18:15:23.847+0800
//        2020-12-30T18:23:42.328+0800
function formatTime(time){
    var times = time.replace(/-/g,':').replace('.', ':' ).replace(' ',':').split(':')
    var resultTime = ""
    if(times.length != 7){
        var date = new Date()
        var hour = date.getHours() <= 9 ? '0' + date.getHours() : date.getHours()
        var month = date.getMonth()+1 <= 9 ? '0' + (date.getMonth() +1) : date.getMonth() + 1
        var day = date.getDate() <=9 ? '0' + date.getDate() : date.getDate()
        var minutes = date.getMinutes() <=9 ? '0' + date.getMinutes() : date.getMinutes()
        var seconds = date.getSeconds() <=9 ? '0' + date.getSeconds() : date.getSeconds()
        var milliSeconds = date.getMilliseconds() < 100 ? ( date.getMilliseconds() < 10 ? '00' + date.getMilliseconds() : '0'+date.getMilliseconds()) : date.getMilliseconds()
        resultTime = date.getFullYear() +"-" + month + "-" + day + "T" +  hour + ":" + minutes +":" + seconds +"." + milliSeconds +"+0800"
    }else {
        resultTime  = times[0] +"-" +times[1] +"-"+times[2] +"T"+times[3] +":"+times[4] +":"+times[5]+"."+times[6]+"+0800"
    }
    return resultTime
}

function process(event) {
    var message = event.Get("message")
    var rule = /\[(.*?)\] \[(.*?)\] \[(.*?)\] \[(.*?)\] \[(.*?)\] (.*)/g
    var res = rule.exec(message)
    if(res){
        event.Put("time",formatTime(res[1]))
        event.Put("level", res[2])
        event.Put("line", res[3])
        event.Put("pathname", res[4])
        event.Put("trace_id", res[5])
        event.Put("message", res[6])
    }else {
        event.Put("time",formatTime(res[1]))
        event.Put("level", 'unknow')
        event.Put("line", 'unknow')
        event.Put("pathname", 'unknow')
        event.Put("trace_id", 'unknow')
    }
}