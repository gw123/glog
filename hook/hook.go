package hook

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"sync"
)

type LogHook struct {
	Field  string
	levels []logrus.Level
	mutex  sync.Mutex
}

func (hook *LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *LogHook) Fire(entry *logrus.Entry) error {
	file := findCaller(4)
	hook.mutex.Lock()
	entry.Data[hook.Field] = file
	hook.mutex.Unlock()
	return nil
}

func NewLogHook(field string, levels []logrus.Level) logrus.Hook {
	if len(levels) == 0 {
		levels = []logrus.Level{logrus.ErrorLevel, logrus.WarnLevel}
		//levels = logrus.AllLevels
	}
	hook := LogHook{
		Field:  field,
		levels: levels,
	}
	return &hook
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}
