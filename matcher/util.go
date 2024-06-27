package matcher

import "reflect"

func typeMatch[T any](v any) bool {
	switch v.(type) {
	case T:
		return true
	default:
		return false
	}
}

func isMatcher(v any) bool {
	_, ok := v.(Matcher)
	return ok
}

func equal(a, b any) bool {
	v, ok := a.(Matcher)
	if ok {
		return v.Match(b)
	}

	return a == b
}

func MaySlice(raw any) *maySlice {
	v := reflect.ValueOf(raw)
	t := v.Type()
	return &maySlice{raw: raw, v: &v, t: &t}
}

type maySlice struct {
	raw any
	t   *reflect.Type
	v   *reflect.Value
}

func (w maySlice) Length() int {
	if !w.IsSlice() {
		return 0
	}

	return w.v.Len()
}

func (w maySlice) Index(n int) (any, bool) {
	if !w.IsSlice() {
		return nil, false
	}

	if n < 0 || n >= w.v.Len() {
		return nil, false
	}

	res := w.v.Index(n).Interface()

	return res, true
}

func (w maySlice) FindIndex(target any, excludes map[int]bool) int {
	for i := 0; i < w.Length(); i++ {
		v, ok := w.Index(i)
		if !ok {
			return -1
		}

		if _, ok := excludes[i]; ok {
			continue
		}

		if equal(v, target) {
			return i
		}
	}

	return -1
}

func (w maySlice) IsSlice() bool {
	if w.t == nil {
		return false
	}
	return isSlice(*w.t)
}

func isSlice(v reflect.Type) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}
