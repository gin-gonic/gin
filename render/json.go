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
	WriteHeader(w, code, "application/json")
	jsonData, err := json.MarshalIndent(data[0], "", "    ")
	if err != nil {
		return err
	}
	_, err = w.Write(jsonData)
	return err
}

func WriteJSON(w http.ResponseWriter, code int, data interface{}) error {
	WriteHeader(w, code, "application/json")
	return json.NewEncoder(w).Encode(data)
}
