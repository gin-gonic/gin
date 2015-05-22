package render

import "net/http"

type Data struct {
	ContentType string
	Data        []byte
}

func (r Data) Write(w http.ResponseWriter) error {
	if len(r.ContentType) > 0 {
		w.Header().Set("Content-Type", r.ContentType)
	}
	w.Write(r.Data)
	return nil
}
