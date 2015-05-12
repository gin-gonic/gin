// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

//TODO
// func (engine *Engine) LoadHTMLGlob(pattern string) {
// func (engine *Engine) LoadHTMLFiles(files ...string) {
// func (engine *Engine) Run(addr string) error {
// func (engine *Engine) RunTLS(addr string, cert string, key string) error {

func init() {
	SetMode(TestMode)
}

func TestLogger(t *testing.T) {
	buffer := new(bytes.Buffer)
	router := New()
	router.Use(LoggerWithWriter(buffer))
	router.GET("/example", func(c *Context) {})

	performRequest(router, "GET", "/example")

	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
}
