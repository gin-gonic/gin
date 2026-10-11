// Copyright 2025 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type bsonBinding struct{}

func (bsonBinding) Name() string {
	return "bson"
}

func (b bsonBinding) Bind(req *http.Request, obj any) error {
	body, err := io.ReadAll(io.LimitReader(req.Body, MaxBodySize+1))
	if err != nil {
		return err
	}
	if int64(len(body)) > MaxBodySize {
		return errors.New("request body too large")
	}
	return b.BindBody(body, obj)
}

func (bsonBinding) BindBody(body []byte, obj any) error {
	return bson.Unmarshal(body, obj)
}
