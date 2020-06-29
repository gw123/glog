package driver

import (
	"fmt"
	"runtime"
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

func logColorful(out string, bg, fg, style int) {
	if runtime.GOOS == "linux" {
		bg = bg % 10
		fg = fg % 10
		style = style % 10

		bg = bg + 40
		bg = fg + 30

		format := fmt.Sprintf("%c[%d;%d;%dm%s%s%c[0m\n",
			0x1B, style, bg, fg, "", out, 0x1B)
		fmt.Println(format)
	} else {
	}
}
