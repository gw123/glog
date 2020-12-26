package logrus_driver

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gw123/glog/common"

	"github.com/sirupsen/logrus"
)

const DateTimeFormat = "2006-01-02 15:04:05"

type GTextFormat struct {
}

func (t GTextFormat) Format(entry *logrus.Entry) ([]byte, error) {
	buf := strings.Builder{}
	buf.WriteString("[")
	buf.WriteString(entry.Time.Format(DateTimeFormat))
	buf.WriteString("]")

	buf.WriteByte(' ')
	buf.WriteString("[")
	buf.WriteString(entry.Level.String())
	buf.WriteString("]")

	buf.WriteByte(' ')
	if entry.HasCaller() {
		buf.WriteString("[")
		arr := strings.Split(entry.Caller.File, "/")
		if len(arr) > 2 {
			arr = arr[len(arr)-2:]
		}
		buf.WriteString(fmt.Sprintf("%s:%d", strings.Join(arr, "/"), entry.Caller.Line))
		buf.WriteString("]")
	} else {
		buf.WriteString("[")
		buf.WriteString("]")
	}

	tmpData := map[string]interface{}{}
	for key, val := range entry.Data {
		if key == common.KeyPathname || key == common.KeyTraceID {
			continue
		}
		tmpData[key] = val
	}

	buf.WriteByte(' ')
	if val, ok := entry.Data[common.KeyPathname]; ok {
		str, _ := val.(string)
		buf.WriteString("[")
		buf.WriteString(str)
		buf.WriteString("]")
	} else {
		buf.WriteString("[")
		buf.WriteString("]")
	}

	buf.WriteByte(' ')
	if val, ok := entry.Data[common.KeyTraceID]; ok {
		str, _ := val.(string)
		buf.WriteString("[")
		buf.WriteString(str)
		buf.WriteString("]")
	} else {
		buf.WriteString("[")
		buf.WriteString("]")
	}

	buf.WriteByte(' ')
	buf.WriteString(entry.Message)

	if tmpData != nil && len(tmpData) > 0 {
		data, err := json.Marshal(tmpData)
		if err != nil {
			return nil, err
		}
		buf.WriteByte(' ')
		buf.Write(data)
	}
	buf.WriteByte('\n')
	res := []byte(buf.String())
	buf.Reset()
	return res, nil
}
