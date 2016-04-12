package render

import (
	"github.com/golang/protobuf/proto"
	"net/http"
)

type Protobuf struct {
	Data interface{}
}

var pbContentType = []string{"application/x-protobuf"}

func (r Protobuf) Render(w http.ResponseWriter) error {
	writeContentType(w, pbContentType)
	encoded, err := proto.Marshal(r.Data.(proto.Message))
	if err != nil {
		return err
	}
	if _, err = w.Write(encoded); err != nil {
		return err
	}
	return nil
}
