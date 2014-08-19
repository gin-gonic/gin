package gin

import (
	"os"
)

const GIN_MODE = "GIN_MODE"

const (
	DebugMode   string = "debug"
	ReleaseMode string = "release"
)
const (
	debugCode   = iota
	releaseCode = iota
)

var gin_mode int = debugCode

func SetMode(value string) {
	switch value {
	case DebugMode:
		gin_mode = debugCode
	case ReleaseMode:
		gin_mode = releaseCode
	default:
		panic("gin mode unknown, the allowed modes are: " + DebugMode + " and " + ReleaseMode)
	}
}

func init() {
	value := os.Getenv(GIN_MODE)
	if len(value) == 0 {
		SetMode(DebugMode)
	} else {
		SetMode(value)
	}
}
