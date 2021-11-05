// Copyright 2020 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build nomsgpack
// +build nomsgpack

package binding

import (
	"context"
	"net/http"
)

// Content-Type MIME of the most common data formats.
const (
	MIMEJSON              = "application/json"
	MIMEHTML              = "text/html"
	MIMEXML               = "application/xml"
	MIMEXML2              = "text/xml"
	MIMEPlain             = "text/plain"
	MIMEPOSTForm          = "application/x-www-form-urlencoded"
	MIMEMultipartPOSTForm = "multipart/form-data"
	MIMEPROTOBUF          = "application/x-protobuf"
	MIMEYAML              = "application/x-yaml"
)

// Binding describes the interface which needs to be implemented for binding the
// data present in the request such as JSON request body, query parameters or
// the form POST.
type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

// ContextBinding enables contextual validation by adding BindContext to Binding.
// Custom validators can take advantage of the information in the context.
type ContextBinding interface {
	Binding
	BindContext(context.Context, *http.Request, interface{}) error
}

// BindingBody adds BindBody method to Binding. BindBody is similar to Bind,
// but it reads the body from supplied bytes instead of req.Body.
type BindingBody interface {
	Binding
	BindBody([]byte, interface{}) error
}

// ContextBindingBody enables contextual validation by adding BindBodyContext to BindingBody.
// Custom validators can take advantage of the information in the context.
type ContextBindingBody interface {
	BindingBody
	BindContext(context.Context, *http.Request, interface{}) error
	BindBodyContext(context.Context, []byte, interface{}) error
}

// BindingUri is similar to Bind, but it read the Params.
type BindingUri interface {
	Name() string
	BindUri(map[string][]string, interface{}) error
}

// ContextBindingUri enables contextual validation by adding BindUriContext to BindingUri.
// Custom validators can take advantage of the information in the context.
type ContextBindingUri interface {
	BindingUri
	BindUriContext(context.Context, map[string][]string, interface{}) error
}

// StructValidator is the minimal interface which needs to be implemented in
// order for it to be used as the validator engine for ensuring the correctness
// of the request. Gin provides a default implementation for this using
// https://github.com/go-playground/validator/tree/v10.6.1.
type StructValidator interface {
	// ValidateStruct can receive any kind of type and it should never panic, even if the configuration is not right.
	// If the received type is a slice/array/map, the validation should be performed on every element.
	// If the received type is not a struct or slice/array/map, any validation should be skipped and nil must be returned.
	// If the received type is a pointer to a struct/slice/array/map, the validation should be performed.
	// If the struct is not valid or the validation itself fails, a descriptive error should be returned.
	// Otherwise nil must be returned.
	ValidateStruct(interface{}) error

	// Engine returns the underlying validator engine which powers the
	// StructValidator implementation.
	Engine() interface{}
}

// ContextStructValidator is an extension of StructValidator that requires implementing
// context-aware validation.
// Custom validators can take advantage of the information in the context.
type ContextStructValidator interface {
	StructValidator
	ValidateStructContext(context.Context, interface{}) error
}

// Validator is the default validator which implements the StructValidator
// interface. It uses https://github.com/go-playground/validator/tree/v10.6.1
// under the hood.
var Validator StructValidator = &defaultValidator{}

// These implement the Binding interface and can be used to bind the data
// present in the request to struct instances.
var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	ProtoBuf      = protobufBinding{}
	YAML          = yamlBinding{}
	Uri           = uriBinding{}
	Header        = headerBinding{}
)

// Default returns the appropriate Binding instance based on the HTTP method
// and the content type.
func Default(method, contentType string) Binding {
	if method == http.MethodGet {
		return Form
	}

	switch contentType {
	case MIMEJSON:
		return JSON
	case MIMEXML, MIMEXML2:
		return XML
	case MIMEPROTOBUF:
		return ProtoBuf
	case MIMEYAML:
		return YAML
	case MIMEMultipartPOSTForm:
		return FormMultipart
	default: // case MIMEPOSTForm:
		return Form
	}
}

func validateContext(ctx context.Context, obj interface{}) error {
	if Validator == nil {
		return nil
	}
	if v, ok := Validator.(ContextStructValidator); ok {
		return v.ValidateStructContext(ctx, obj)
	}
	return Validator.ValidateStruct(obj)
}
