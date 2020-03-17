package driver

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

var dataBuffer *bytes.Buffer

type LogOutPut struct {
	buffer   *bytes.Buffer
	interval int
	writer   io.Writer
	stop     bool
}

func NewLogOutPut(io io.Writer) *LogOutPut {
	byteArr := make([]byte, 1>>20)
	dataBuffer = bytes.NewBuffer(byteArr)
	out := &LogOutPut{
		buffer:   dataBuffer,
		interval: 0,
		writer:   io,
	}
	out.Start()
	return out
}

func (l *LogOutPut) Write(data []byte) {
	l.buffer.Write([]byte(data))
}

func (l *LogOutPut) Stop()  {
	l.stop = true
}

func (l *LogOutPut) Start() {
	go func() {
		var buffer = make([]byte, 4096)
		var length = 1
		var err error
		for ; ; {
			if l.buffer.Len() == 0 {
				for l.buffer.Len() == 0{
					time.Sleep(time.Second)
				}
			}
			length, err = l.buffer.Read(buffer)
			if err != nil {
				fmt.Println("LogManager start", err)
				break
			}
			for ; length >0; {
				_, err = l.writer.Write(buffer[0:length])
				if err != nil {
					fmt.Println("LogManager Write", err)
					if err.Error() == "EOF" {
						break
					}
				}
			}
		}
	}()
}

//func SetTraceId(id string) {
//	traceIdName = id
//}

////向context中注入 traceId
//func WithTraceID(ctx context.Context) (context.Context, int64) {
//	nano := time.Now().UnixNano()
//	var id int64
//	id = rand.Int63()
//	id = nano & (id & 0xffffffff000000000)
//	return context.WithValue(ctx, traceIdName, id), id
//}
//
////记录调用位置 ,调用时间点 . ctx 必须存在否则触发panic
//func RecordPoint(ctx context.Context, point string) (context.Context, error) {
//	traceId, ok := ctx.Value(traceIdName).(int64)
//	if !ok {
//		ctx, traceId = WithTraceID(ctx)
//	}
//	record, ok := ctx.Value(lastLogRecord).(*Record)
//	if !ok {
//		record = &Record{}
//	}
//	record.Point = point
//	record.TraceId = traceId
//	record.CreatedAt = time.Now().UnixNano()
//	return
//}
//
////记录调用位置 ,调用时间点 . ctx 必须存在否则触发panic
//func RecordPointWithCall(ctx context.Context, point string) (context.Context, error) {
//	traceId, ok := ctx.Value(traceIdName).(int64)
//	if !ok {
//		ctx, traceId = WithTraceID(ctx)
//	}
//	record, ok := ctx.Value(lastLogRecord).(*Record)
//	if !ok {
//		record = &Record{}
//	}
//	record.Place = runtime.Caller(1)
//	item, err := json.Marshal(record)
//	if err != nil {
//		return
//	}
//
//	return nil
//}
