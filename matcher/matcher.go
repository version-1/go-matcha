package matcher

import "reflect"

func equal(a, b any) bool {
	v, ok := a.(Matcher)
	if ok {
		return v.Match(b)
	}

	return a == b
}

func IsMatcher(v any) bool {
	_, ok := v.(Matcher)
	return ok
}

type Matcher interface {
	Match(v any) bool
	Pointer() Matcher
	Not() Matcher
}

type RefMatcher struct {
	m Matcher
}

var _ Matcher = RefMatcher{}

func (r RefMatcher) Match(v any) bool {
	if v == nil {
		return r.m.Match(nil)
	}
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr {
		return false
	}

	return r.m.Match(vv.Elem().Interface())
}

func (r RefMatcher) Not() Matcher {
	return Not(r)
}

func (r RefMatcher) Pointer() Matcher {
	return Ref(r)
}

func Ref(m Matcher) Matcher {
	return &RefMatcher{m: m}
}

type notMatcher struct {
	m Matcher
}

var _ Matcher = notMatcher{}

func (e notMatcher) Match(v any) bool {
	return !e.m.Match(v)
}

func (e notMatcher) Pointer() Matcher {
	return Ref(e)
}

func (e notMatcher) Not() Matcher {
	return Not(e)
}

func Not(m Matcher) Matcher {
	return &notMatcher{m: m}
}
