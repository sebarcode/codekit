package codekit

import "reflect"

func IsNilOrEmpty(x interface{}) bool {
	if x == nil {
		return true
	}

	v := reflect.Indirect(reflect.ValueOf(x))
	return v.IsZero()
}
