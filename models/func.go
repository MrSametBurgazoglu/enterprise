package models

import (
	"reflect"
	"runtime"
	"strings"
)

type FuncStruct struct {
	PackageAddress string
	PackageFunc    string
}

func (f *FuncStruct) DefaultFunc(v any) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Func {
		panic("not a function")
	}
	pc := val.Pointer()
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		panic("unknown function")
	}
	fullName := fn.Name()

	i1 := strings.LastIndex(fullName, ".")
	f.PackageAddress = fullName[:i1]
	i2 := strings.LastIndex(fullName, "/")
	f.PackageFunc = fullName[i2+1:]
}
