package gin

import (
	"reflect"
	"runtime"
)

type Function struct {
	Name string
	File string
	Line int
}

func parseFunction(f interface{}) Function {
	ptr := reflect.ValueOf(f).Pointer()
	fu := runtime.FuncForPC(ptr)
	file, line := fu.FileLine(ptr)
	return Function{
		Name: fu.Name(),
		File: file,
		Line: line,
	}
}

func nameOfFunction(f interface{}) string {
	return parseFunction(f).Name
}
