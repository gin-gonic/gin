package render

import "net/http"

type dataRender struct{}

func (_ dataRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	contentType := data[0].(string)
	bytes := data[1].([]byte)
	WriteData(w, code, contentType, bytes)
	return nil
}

func WriteData(w http.ResponseWriter, code int, contentType string, data []byte) {
	if len(contentType) > 0 {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(code)
	w.Write(data)
}
