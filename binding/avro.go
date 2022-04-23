// Copyright 2022 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin/internal/json"
	"github.com/hamba/avro"
)

type avroBinding struct {
	s string
}

func (avroBinding) Name() string {
	return "avro"
}

func (r avroBinding) Bind(req *http.Request, obj any) error {
	return decodeAvro(req.Body, r.s, obj)
}

func (avroBinding) BindBody(body []byte, s string, obj any) error {
	return decodeAvro(bytes.NewReader(body), s, obj)
}

func decodeAvro(r io.Reader, s string, obj any) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		return err
	}
	schema, err := avro.Parse(s)
	if err != nil {
		return err
	}
	data, err := avro.Marshal(schema, obj)
	if err != nil {
		return err
	}
	decoder, err := avro.NewDecoder(s, bytes.NewReader(data))
	if err != nil {
		return err
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
