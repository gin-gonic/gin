package yaml

// StructValidator need to implement Struct method only
// ( see https://pkg.go.dev/github.com/go-playground/validator/v10#Validate.Struct )
type StructValidator interface {
	Struct(interface{}) error
}

// FieldError need to implement StructField method only
// ( see https://pkg.go.dev/github.com/go-playground/validator/v10#FieldError )
type FieldError interface {
	StructField() string
}
