package parser

import "fmt"

const (
	colorFgHiBlack int = iota + 90
	colorFgHiRed
	colorFgHiGreen
	colorFgHiYellow
	colorFgHiBlue
	colorFgHiMagenta
	colorFgHiCyan
)

var colorTable = []int{
	colorFgHiRed,
	colorFgHiGreen,
	colorFgHiYellow,
	colorFgHiBlue,
	colorFgHiMagenta,
	colorFgHiCyan,
}

func colorize(idx int, content string) string {
	colorIdx := idx % len(colorTable)
	color := colorTable[colorIdx]
	return fmt.Sprintf("\x1b[1;%dm", color) + content + "\x1b[22;0m"
}
