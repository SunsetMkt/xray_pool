package node

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"strings"
)

// RepeatChar 生成num个相同字符组成的字符串
func RepeatChar(ch byte, num int) string {
	return strings.Repeat(string(ch), num)
}

// MaxWidth 所有字符串中最大宽度
func MaxWidth(str ...string) int {
	max := 0
	for _, s := range str {
		width := runewidth.StringWidth(s)
		if width > max {
			max = width
		}
	}
	return max
}

// ShowTopBottomSepLine 添加上下的分割线
func ShowTopBottomSepLine(ch byte, str ...string) {
	width := MaxWidth(str...)
	fmt.Println(RepeatChar(ch, width))
	fmt.Println(strings.Join(str, "\n"))
	fmt.Println(RepeatChar(ch, width))
}
