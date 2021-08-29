package goevents

import (
	"reflect"
)

func typeName(v interface{}) string {
	name := ""
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	return name
}
