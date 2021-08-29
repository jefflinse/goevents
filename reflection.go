package goevents

import (
	"reflect"
)

// Returns the type name of the given value.
func typeName(v interface{}) string {
	name := ""
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	return name
}
