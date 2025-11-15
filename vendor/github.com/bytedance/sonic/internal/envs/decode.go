package envs

import (
	"os"
)

var UseOptDec  = os.Getenv("SONIC_USE_OPTDEC")  == "1" 
var UseFastMap = os.Getenv("SONIC_USE_FASTMAP") == "1" 

func EnableOptDec() {
	UseOptDec = true
}

func DisableOptDec() {
	UseOptDec = false
}

func EnableFastMap() {
	UseFastMap = true
}

func DisableFastMap() {
	UseFastMap = false
}