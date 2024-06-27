package matcher

import (
	"reflect"
)

func BeAny() *beAny {
	return &beAny{}
}

type beAny struct {
	options MatcherOptions
}

func (m beAny) Match(v any) bool {
	if !m.options.AllowZero && isZero(v) {
		return false
	}

	return true
}

func (m beAny) Not() Matcher {
	return Not(m)
}

func (m beAny) Pointer() Matcher {
	return Ref(m)
}

func (m beAny) AllowZero() Matcher {
	m.options.AllowZero = true
	return m
}

func BeZero() *beZero {
	return &beZero{}
}

type beZero struct{}

func (b beZero) Match(v any) bool {
	return isZero(v)
}

func (b beZero) Not() Matcher {
	return Not(b)
}

func (b beZero) Pointer() Matcher {
	return Ref(b)
}

type notMatcher struct {
	m Matcher
}

func (m notMatcher) Match(v any) bool {
	return !m.m.Match(v)
}

func (m notMatcher) Pointer() Matcher {
	return Ref(m)
}

func (m notMatcher) Not() Matcher {
	return Not(m)
}

func Not(m Matcher) Matcher {
	return &notMatcher{m: m}
}

type RefMatcher struct {
	m Matcher
}

func (r RefMatcher) Title() string {
	v, ok := r.m.(Recorder)
	if ok {
		return v.Title()
	}

	return ""
}

func (r RefMatcher) Records() []Record {
	v, ok := r.m.(Recorder)
	if ok {
		return v.Records()
	}

	return []Record{}
}

func (r RefMatcher) Match(v any) bool {
	if v == nil {
		return r.m.Match(nil)
	}

	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr {
		return false
	}

	e := vv.Elem()
	if !e.IsValid() {
		return r.m.Match(nil)
	}

	return r.m.Match(e.Interface())
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

func ExtractIfPossible(mayRef any) any {
	p, ok := mayRef.(*RefMatcher)
	if ok {
		return p.m
	}

	return mayRef
}
