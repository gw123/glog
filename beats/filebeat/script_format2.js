// 解析message格式消息到json格式
// [2020-12-28 19:28:42.325] [info] [] cluster/service_manager.go 190 [] - ServiceManager[point].
// 2020-12-29T18:15:23.847+0800

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
    var rule = /\[(.*?)\] \[(.*?)\] \[(.*?)\] (.+) (\d{1,4}) \[(.*?)\] - (.*)/g
    var res = rule.exec(message)

    if(res){
        event.Put("level", res[2])
        event.Put("line", res[4] + " "+res[5])
        event.Put("message", res[7])
        event.Put("timestamp", formatTime(res[1]))
    }
}