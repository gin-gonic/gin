// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// Server-Sent Events
// W3C Working Draft 29 October 2009
// http://www.w3.org/TR/2009/WD-eventsource-20091029/

const ContentType = "text/event-stream;charset=utf-8"

var (
	contentType = []string{ContentType}
	noCache     = []string{"no-cache"}
)

var fieldReplacer = strings.NewReplacer(
	"\n", "\\n",
	"\r", "\\r")

var dataReplacer = strings.NewReplacer(
	"\n", "\ndata: ",
	"\r", "\\r")

type Event struct {
	Event string
	Id    string
	Retry uint
	Data  interface{}
}

func Encode(writer io.Writer, event Event) error {
	w := checkWriter(writer)
	writeId(w, event.Id)
	writeEvent(w, event.Event)
	writeRetry(w, event.Retry)
	return writeData(w, event.Data)
}

func writeId(w stringWriter, id string) {
	if len(id) > 0 {
		_, _ = w.WriteString("id: ")
		_, _ = fieldReplacer.WriteString(w, id)
		_, _ = w.WriteString("\n")
	}
}

func writeEvent(w stringWriter, event string) {
	if len(event) > 0 {
		_, _ = w.WriteString("event: ")
		_, _ = fieldReplacer.WriteString(w, event)
		_, _ = w.WriteString("\n")
	}
}

func writeRetry(w stringWriter, retry uint) {
	if retry > 0 {
		_, _ = w.WriteString("retry: ")
		_, _ = w.WriteString(strconv.FormatUint(uint64(retry), 10))
		_, _ = w.WriteString("\n")
	}
}

func writeData(w stringWriter, data interface{}) error {
	_, _ = w.WriteString("data: ")

	bData, ok := data.([]byte)
	if ok {
		_, _ = dataReplacer.WriteString(w, string(bData))
		_, _ = w.WriteString("\n\n")
		return nil
	}

	switch kindOfData(data) { //nolint:exhaustive
	case reflect.Struct, reflect.Slice, reflect.Map:
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			return err
		}
		_, _ = w.WriteString("\n")
	default:
		_, _ = dataReplacer.WriteString(w, fmt.Sprint(data))
		_, _ = w.WriteString("\n\n")
	}
	return nil
}

func (r Event) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return Encode(w, r)
}

func (r Event) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	header["Content-Type"] = contentType

	if _, exist := header["Cache-Control"]; !exist {
		header["Cache-Control"] = noCache
	}
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
