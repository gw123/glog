package glog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	ColorBlack       = 0
	ColorRed         = 1
	ColorGreen       = 2
	ColorYello       = 3
	ColorBule        = 4
	ColorAmaranth    = 5
	ColorUltramarine = 6
	ColorWhite       = 7

	Style0 = 0
	Style1 = 1
	Style4 = 4
	Style5 = 5
	Style7 = 7
	Style8 = 8
)

func Dump(input interface{}, label ...string) {
	out1, err := json.MarshalIndent(input, "", " ")
	if err != nil {
		fmt.Println("Dump", err.Error())
		return
	}
	if len(label) >= 1 {
		logColorful("label:"+label[0]+string(out1), ColorUltramarine, ColorUltramarine, 1)
	} else {
		logColorful(string(out1), ColorUltramarine, ColorUltramarine, 1)
	}
}

func Error(format string, other ...interface{}) {
	str := fmt.Sprintf(format, other...)
	logColorful("[E] "+str, ColorBlack, ColorRed, Style1)
}

func Warn(format string, other ...interface{}) {
	str := fmt.Sprintf(format, other...)
	logColorful("[W] "+str, ColorBlack, ColorYello, Style1)
}

func Info(format string, other ...interface{}) {
	str := fmt.Sprintf(format, other...)
	logColorful("[I] "+str, ColorBlack, ColorUltramarine, Style1)
}

func Debug(format string, other ...interface{}) {
	str := fmt.Sprintf(format, other...)
	logColorful("[D] "+str, ColorBlack, ColorUltramarine, Style1)
}

func Output(out string) {
	_, file, line, ok := runtime.Caller(3)
	if ok {
		gopath, _ := os.LookupEnv("GOPATH")
		gopath = strings.Replace(gopath, "\\", "/", -1)
		file = strings.Replace(file, gopath+"/src/", "", -1)
		strs := strings.Split(file, "/")
		if len(strs) > 3 {
			strs = strs[len(strs)-2:]
		}
		file = strings.Join(strs, "/")
		log.Printf("[%s %d] %s", file, line, out)
	} else {
		fmt.Println("Dump caller 1 err")
	}
	return
}

func logColorful(out string, bg, fg, style int) {
	if runtime.GOOS == "linux" {
		bg = bg % 10
		fg = fg % 10
		style = style % 10

		bg = bg + 40
		bg = fg + 30

		format := fmt.Sprintf("%c[%d;%d;%dm%s%s%c[0m\n",
			0x1B, style, bg, fg, "", out, 0x1B)
		Output(format)
	} else {
		Output(out)
	}
}
