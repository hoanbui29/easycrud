package validator

import (
	"reflect"
)

func IsStruct(t reflect.Type) bool {
	if t.Kind() == reflect.Struct {
		return true
	}

	return false
}
