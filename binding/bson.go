// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type bsonBinding struct{}

func (bsonBinding) Name() string {
	return "bson"
}

func (b bsonBinding) Bind(req *http.Request, obj any) error {
	buf, err := io.ReadAll(req.Body)
	if err == nil {
		err = b.BindBody(buf, obj)
	}
	return err
}

func (bsonBinding) BindBody(body []byte, obj any) error {
	return bson.Unmarshal(body, obj)
}
