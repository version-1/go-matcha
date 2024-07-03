package matcher

import (
	"reflect"
)

func typeMatch[T any](v any) bool {
	switch v.(type) {
	case T:
		return true
	default:
		return false
	}
}

func isZero(v any) bool {
	if v == nil {
		return true
	}

	vv := reflect.ValueOf(v)

	if !vv.IsValid() {
		return true
	}

	vt := vv.Type()
	if isSlice(vt) {
		return vv.Len() == 0
	}

	return v == reflect.Zero(vt).Interface()
}
