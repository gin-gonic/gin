package render

import (
	"fmt"
	"net/http"
)

type plainTextRender struct{}

func (_ plainTextRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	format := data[0].(string)
	values := data[1].([]interface{})
	WritePlainText(w, code, format, values)
	return nil
}

func WritePlainText(w http.ResponseWriter, code int, format string, values []interface{}) {
	writeHeader(w, code, "text/plain; charset=utf-8")
	// we assume w.Write can not fail, is that right?
	if len(values) > 0 {
		fmt.Fprintf(w, format, values...)
	} else {
		w.Write([]byte(format))
	}
}
