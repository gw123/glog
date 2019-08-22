package main

import (
	"fmt"
	"github.com/gw123/glog"
)

func plattle() {
	// 前景 背景 颜色
	// ---------------------------------------
	// 30  40  黑色
	// 31  41  红色
	// 32  42  绿色
	// 33  43  黄色
	// 34  44  蓝色
	// 35  45  紫红色
	// 36  46  青蓝色
	// 37  47  白色
	//
	// 代码 意义
	// -------------------------
	//  0  终端默认设置
	//  1  高亮显示
	//  4  使用下划线
	//  5  闪烁
	//  7  反白显示
	//  8  不可见

	for b := 40; b <= 47; b++ { // 背景色彩 = 40-47
		for f := 30; f <= 37; f++ { // 前景色彩 = 30-37
			for d := range []int{0, 1, 4, 5, 7, 8} { // 显示方式 = 0,1,4,5,7,8
				fmt.Printf(" %c[%d;%d;%dm%s(f=%d,b=%d,d=%d)%c[0m ",
					0x1B, d, b, f, "", f, b, d, 0x1B)
			}
			fmt.Println("")
		}
		fmt.Println("")
	}
}

func main() {
	glog.Error("name:%s", "wwww")
	glog.Info("name:%s", "wwww")
	glog.Warn("name:%s", "wwww")
	glog.Debug("name:%s", "wwww")
}
