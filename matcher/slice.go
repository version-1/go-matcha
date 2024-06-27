package matcher

import (
	"reflect"
)

type anySlice struct{}

var _ Matcher = anySlice{}

func BeSlice() *anySlice {
	return &anySlice{}
}

func (a anySlice) Match(v any) bool {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

func (a anySlice) Not() Matcher {
	return Not(a)
}

func (a anySlice) Pointer() Matcher {
	return Ref(a)
}

func SliceOf(elements []any, opts ...func(m *SliceOfMatcherOptions)) Matcher {
	o := SliceOfMatcherOptions{
		Order:    true,
		Contains: false,
	}
	for _, opt := range opts {
		opt(&o)
	}

	return sliceOfMatcher{elements: elements, options: o}
}

func WithSliceOfPersistOrder(v bool) func(*SliceOfMatcherOptions) {
	return func(o *SliceOfMatcherOptions) {
		o.Order = v
	}
}

func WithSliceOfContains(v bool) func(*SliceOfMatcherOptions) {
	return func(o *SliceOfMatcherOptions) {
		o.Contains = v
	}
}

type SliceOfMatcherOptions struct {
	Order    bool
	Contains bool
}

type sliceOfMatcher struct {
	elements []any
	options  SliceOfMatcherOptions
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
