package json

import (
	"fmt"

	"github.com/goccy/go-json/internal/encoder"
)

type (
	ColorFormat = encoder.ColorFormat
	ColorScheme = encoder.ColorScheme
)

const escape = "\x1b"

type colorAttr int

//nolint:deadcode,varcheck
const (
	fgBlackColor colorAttr = iota + 30
	fgRedColor
	fgGreenColor
	fgYellowColor
	fgBlueColor
	fgMagentaColor
	fgCyanColor
	fgWhiteColor
)

//nolint:deadcode,varcheck
const (
	fgHiBlackColor colorAttr = iota + 90
	fgHiRedColor
	fgHiGreenColor
	fgHiYellowColor
	fgHiBlueColor
	fgHiMagentaColor
	fgHiCyanColor
	fgHiWhiteColor
)

func createColorFormat(attr colorAttr) ColorFormat {
	return ColorFormat{
		Header: wrapColor(attr),
		Footer: resetColor(),
	}
}

func wrapColor(attr colorAttr) string {
	return fmt.Sprintf("%s[%dm", escape, attr)
}

func resetColor() string {
	return wrapColor(colorAttr(0))
}

var (
	DefaultColorScheme = &ColorScheme{
		Int:       createColorFormat(fgHiMagentaColor),
		Uint:      createColorFormat(fgHiMagentaColor),
		Float:     createColorFormat(fgHiMagentaColor),
		Bool:      createColorFormat(fgHiYellowColor),
		String:    createColorFormat(fgHiGreenColor),
		Binary:    createColorFormat(fgHiRedColor),
		ObjectKey: createColorFormat(fgHiCyanColor),
		Null:      createColorFormat(fgBlueColor),
	}
)
