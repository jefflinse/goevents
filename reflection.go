package goevents

import (
	"reflect"
	"strings"
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

func EventName(v interface{}) string {
	return strings.TrimSuffix(typeName(v), "Event")
}

func CommandName(v interface{}) string {
	return strings.TrimSuffix(typeName(v), "Command")
}
