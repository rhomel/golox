package check

import "reflect"

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func TypeOf(it interface{}) string {
	typ := reflect.TypeOf(it)
	switch typ.Kind() {
	case reflect.Ptr:
		return typ.Elem().Name() + " (ptr)"
	default:
		return typ.Name()
	}
}
