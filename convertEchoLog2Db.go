package glog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type LogItem struct {
	Http struct {
		ByteIn  string `json:"byte_in"`
		ByteOut string `json:"byte_out"`
		Host    string `json:"host"`
		Latency string `json:"latency"`
		Method  string `json:"method"`
		Remote  string `json:"remote"`
		Status  int    `json:"status"`
		Uri     string `json:"uri"`
	} `json:"http"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Scope string `json:"scope"`
	Time  string `json:"time"`
}

type LogModel struct {
	gorm.Model
	Latency int       `json:"latency"`
	Method  string    `json:"method"`
	Remote  string    `json:"remote"`
	Status  int       `json:"status"`
	Uri     string    `json:"uri"`
	Level   string    `json:"level"`
	Msg     string    `json:"msg"`
	Time    time.Time `json:"time"`
	ByteIn  string    `json:"bytes_in"`
	ByteOut string    `json:"bytes_out"`
}

var gDbInstance *gorm.DB

func insertDb(logitem *LogItem) error {
	if logitem == nil {
		return errors.New("logitem is nil")
	}
	loc, _ := time.LoadLocation("Local")
	reqtime, err := time.ParseInLocation("2006-01-02T15:04:05+08:00", logitem.Time, loc)
	if err != nil {
		return err
	}
	nums := strings.Split(logitem.Http.Latency, ".")
	latency, err := strconv.Atoi(nums[0])
	if err != nil {
		return err
	}
	newmodel := &LogModel{
		Latency: latency,
		Method:  logitem.Http.Method,
		Remote:  logitem.Http.Remote,
		Status:  logitem.Http.Status,
		Uri:     logitem.Http.Uri,
		Level:   logitem.Level,
		Msg:     logitem.Msg,
		Time:    reqtime,
		ByteIn:  logitem.Http.ByteIn,
		ByteOut: logitem.Http.ByteOut,
	}
	return gDbInstance.Save(newmodel).Error
}

func NewDb(drive, db_host, db_database, db_username, db_pwd string) (*gorm.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_username, db_pwd, db_host, db_database)
	dbInstance, err := gorm.Open(drive, connStr)
	if err != nil {
		return nil, err
	}
	gDbInstance = dbInstance
	return dbInstance, nil
}

func parseFile(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	buf := bufio.NewReader(f)
	i := 0
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		i++
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		logItem := &LogItem{}
		err = json.Unmarshal([]byte(line), logItem)
		if err != nil {
			fmt.Println("=====================")
			fmt.Println("Error", err)
		} else {
			if logItem.Http.Latency != "" {
				fmt.Println("+++++++++++++++++++++ ", i)
				fmt.Println(logItem)
				err := insertDb(logItem)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

	}
	return nil
}

//func main() {
//	_, err := NewDb("mysql", "192.168.30.138", "pos", "gw", "gao123456")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	g_dbInstance.AutoMigrate(&LogModel{})
//
//	fileName := "C:/Users/gw/Desktop/lh/pos机器/pos/pos.log"
//	parseFile(fileName)
//}
