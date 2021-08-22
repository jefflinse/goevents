package goevents

import (
	"reflect"
	"strings"
)

func getEventType(v interface{}) string {
	name := ""
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	return strings.TrimSuffix(name, "Event")
}
