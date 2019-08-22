package glog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

func SetOutput() {

}

func Dump(input interface{}) {
	out1, err := json.MarshalIndent(input, "", " ")
	if err != nil {
		fmt.Println("Dump", err.Error())
		return
	}
	output(string(out1))
}

func Error(out string) {
	log.Printf("%s", out)
}

func output(out string) {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		gopath, _ := os.LookupEnv("GOPATH")
		gopath = strings.Replace(gopath, "\\", "/", -1)
		file = strings.Replace(file, gopath+"/src/", "", -1)
		log.Printf("[%s %d] %s", file, line, out)
	} else {
		fmt.Println("Dump caller 1 err")
	}
	return
}
