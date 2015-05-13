package render

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
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
		w.Header().Set("Content-Type", "text/event-stream")
	}
	var stringData string
	switch typeOfData(data) {
	case reflect.Struct, reflect.Slice, reflect.Map:
		if jsonBytes, err := json.Marshal(data); err == nil {
			stringData = string(jsonBytes)
		} else {
			return err
		}
	case reflect.Ptr:
		stringData = escape(fmt.Sprintf("%v", &data))
	default:
		stringData = escape(fmt.Sprintf("%v", data))
	}
	_, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", escape(eventName), stringData)
	return err
}

func typeOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		newValue := value.Elem().Kind()
		fmt.Println(newValue)
		if newValue == reflect.Struct ||
			newValue == reflect.Slice ||
			newValue == reflect.Map {
			return newValue
		} else {
			return valueType
		}
	}
	return valueType
}

func escape(str string) string {
	return str
}
