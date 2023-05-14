package util

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

func FunctionFullName(i interface{}) string {
	pointer := reflect.ValueOf(i).Pointer()
	funcForPC := runtime.FuncForPC(pointer)
	if funcForPC == nil {
		return ""
	}

	file, line := funcForPC.FileLine(pointer)

	return fmt.Sprintf("%s:%d", file, line)
}

func FunctionName(i interface{}) string {
	pointer := reflect.ValueOf(i).Pointer()
	funcForPC := runtime.FuncForPC(pointer)
	if funcForPC == nil {
		return ""
	}

	file, line := funcForPC.FileLine(pointer)

	files := strings.Split(file, "/")

	return fmt.Sprintf("%s:%d", files[len(files)-1], line)
}
