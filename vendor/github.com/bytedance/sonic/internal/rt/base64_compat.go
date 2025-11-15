// +build !amd64 !go1.17 go1.26

package rt

import (
	"encoding/base64"
)

func DecodeBase64(raw []byte) ([]byte, error) {
	ret := make([]byte, base64.StdEncoding.DecodedLen(len(raw)))
	n, err := base64.StdEncoding.Decode(ret, raw)
	if err != nil {
		return nil, err
	}
	return ret[:n], nil
}

func EncodeBase64ToString(src []byte) string {
    return base64.StdEncoding.EncodeToString(src)
}

func EncodeBase64(buf []byte, src []byte) []byte {
	if len(src) == 0 {
		return append(buf, '"', '"')
	}
	buf = append(buf, '"')
	need := base64.StdEncoding.EncodedLen(len(src))
	if cap(buf) - len(buf) < need {
		tmp := make([]byte, len(buf), len(buf) + need*2)
		copy(tmp, buf)
		buf = tmp
	}
	base64.StdEncoding.Encode(buf[len(buf):cap(buf)], src)
	buf = buf[:len(buf) + need]
	buf = append(buf, '"')
	return buf
}
