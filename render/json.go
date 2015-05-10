package render

import (
	"encoding/json"
	"net/http"
)

type (
	jsonRender struct{}

	indentedJSON struct{}
)

func (_ jsonRender) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	return WriteJSON(w, code, data[0])
}

func (_ indentedJSON) Render(w http.ResponseWriter, code int, data ...interface{}) error {
	return WriteIndentedJSON(w, code, data[0])
}

func WriteJSON(w http.ResponseWriter, code int, data interface{}) error {
	writeHeader(w, code, "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(data)
}

func WriteIndentedJSON(w http.ResponseWriter, code int, data interface{}) error {
	writeHeader(w, code, "application/json; charset=utf-8")
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(jsonData)
	return err
}
