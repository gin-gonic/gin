

package base64x

import (
	"github.com/cloudwego/base64x/internal/native"
)

// HACK: maintain these only to prevent breakchange, because sonic-go linkname these
var (
	_subr__b64decode uintptr
	_subr__b64encode uintptr
)

func init() {
	_subr__b64decode = native.S_b64decode
	_subr__b64encode = native.S_b64encode
}
