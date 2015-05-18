package render

import (
	"encoding/json"
	"net/http"
)

type (
	JSON struct {
		Data interface{}
	}

	IndentedJSON struct {
		Data interface{}
	}
)

func (r JSON) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(r.Data)
}

func (r IndentedJSON) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	jsonBytes, err := json.MarshalIndent(r.Data, "", "    ")
	if err != nil {
		return err
	}
	w.Write(jsonBytes)
	return nil
}
