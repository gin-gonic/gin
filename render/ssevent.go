package render

import (
	"net/http"

	"github.com/manucorporat/sse"
)

type sseRender struct{}

var SSEvent Render = sseRender{}

func (_ sseRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	eventName := data[0].(string)
	obj := data[1]
	return WriteSSEvent(w, eventName, obj)
}

func WriteSSEvent(w http.ResponseWriter, eventName string, data interface{}) error {
	header := w.Header()
	if len(header.Get("Content-Type")) == 0 {
		header.Set("Content-Type", sse.ContentType)
	}
	if len(header.Get("Cache-Control")) == 0 {
		header.Set("Cache-Control", "no-cache")
	}
	return sse.Encode(w, sse.Event{
		Event: eventName,
		Data:  data,
	})
}
