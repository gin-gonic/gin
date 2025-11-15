// This source inspired by https://github.com/fatih/color.
package printer

import (
	"fmt"
	"strings"
)

type ColorAttribute int

const (
	ColorReset ColorAttribute = iota
	ColorBold
	ColorFaint
	ColorItalic
	ColorUnderline
	ColorBlinkSlow
	ColorBlinkRapid
	ColorReverseVideo
	ColorConcealed
	ColorCrossedOut
)

const (
	ColorFgHiBlack ColorAttribute = iota + 90
	ColorFgHiRed
	ColorFgHiGreen
	ColorFgHiYellow
	ColorFgHiBlue
	ColorFgHiMagenta
	ColorFgHiCyan
	ColorFgHiWhite
)

const (
	ColorResetBold ColorAttribute = iota + 22
	ColorResetItalic
	ColorResetUnderline
	ColorResetBlinking

	ColorResetReversed
	ColorResetConcealed
	ColorResetCrossedOut
)

const escape = "\x1b"

var colorResetMap = map[ColorAttribute]ColorAttribute{
	ColorBold:         ColorResetBold,
	ColorFaint:        ColorResetBold,
	ColorItalic:       ColorResetItalic,
	ColorUnderline:    ColorResetUnderline,
	ColorBlinkSlow:    ColorResetBlinking,
	ColorBlinkRapid:   ColorResetBlinking,
	ColorReverseVideo: ColorResetReversed,
	ColorConcealed:    ColorResetConcealed,
	ColorCrossedOut:   ColorResetCrossedOut,
}

func format(attrs ...ColorAttribute) string {
	format := make([]string, 0, len(attrs))
	for _, attr := range attrs {
		format = append(format, fmt.Sprint(attr))
	}
	return fmt.Sprintf("%s[%sm", escape, strings.Join(format, ";"))
}

func unformat(attrs ...ColorAttribute) string {
	format := make([]string, len(attrs))
	for _, attr := range attrs {
		v := fmt.Sprint(ColorReset)
		reset, exists := colorResetMap[attr]
		if exists {
			v = fmt.Sprint(reset)
		}
		format = append(format, v)
	}
	return fmt.Sprintf("%s[%sm", escape, strings.Join(format, ";"))
}

func colorize(msg string, attrs ...ColorAttribute) string {
	return format(attrs...) + msg + unformat(attrs...)
}
