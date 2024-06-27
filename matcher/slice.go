package matcher

import (
	"reflect"
	"testing"

	"github.com/version-1/go-matcha/matcher/slices"
)

type anySlice struct{}

var _ Matcher = anySlice{}

func BeSlice() *anySlice {
	return &anySlice{}
}

func (a anySlice) Match(v any) bool {
	s := MaySlice(v)
	return s.IsSlice()
}

func (a anySlice) Not() Matcher {
	return Not(a)
}

func (a anySlice) Pointer() Matcher {
	return Ref(a)
}

func SliceOf(elements []any, opts ...func(m *slices.MatcherOptions)) Matcher {
	o := slices.MatcherOptions{
		Order:    true,
		Contains: false,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return sliceOfMatcher{elements: elements, options: o}
}

type sliceOfMatcher struct {
	elements []any
	options  slices.MatcherOptions
}

var _ Matcher = sliceOfMatcher{}

func (m sliceOfMatcher) Match(v any) bool {
	vw := MaySlice(v)
	if !vw.IsSlice() {
		return false
	}

	if !m.options.Contains && len(m.elements) != vw.Length() {
		return false
	}

	if m.options.Order {
		for i := range m.elements {
			ele, ok := vw.Index(i)
			if !ok {
				return false
			}

			if !equal(m.elements[i], ele) {
				return false
			}
		}
		return true
	}

	maps := map[int]bool{}
	for i := 0; i < len(m.elements); i++ {
		idx := vw.FindIndex(m.elements[i], maps)
		if idx < 0 {
			return false
		}
		maps[idx] = true
	}

	return true
}

func (m sliceOfMatcher) Not() Matcher {
	return Not(m)
}

func (m sliceOfMatcher) Pointer() Matcher {
	return Ref(m)
}

type sliceLenMatcher struct {
	n int
}

func SliceLen(n int) Matcher {
	return sliceLenMatcher{n: n}
}

func (m sliceLenMatcher) Match(v any) bool {
	vw := MaySlice(v)
	if !vw.IsSlice() {
		return false
	}

	return vw.Length() == m.n
}

func (m sliceLenMatcher) Not() Matcher {
	return Not(m)
}

func (m sliceLenMatcher) Pointer() Matcher {
	return Ref(m)
}

func MaySlice(raw any) *maySlice {
	v := reflect.ValueOf(raw)
	t := v.Type()
	return &maySlice{raw: raw, v: &v, t: &t}
}

type maySlice struct {
	raw  any
	t    *reflect.Type
	v    *reflect.Value
	test *testing.T
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
