package matcher

import (
	"reflect"

	"github.com/version-1/go-matcha/matcher/structs"
)

func BeStruct() *anyStruct {
	return &anyStruct{}
}

type anyStruct struct{}

var _ Matcher = anyStruct{}

func (a anyStruct) Match(v any) bool {
	if v == nil {
		return false
	}

	s := MayStruct(v)
	return s.IsStruct()
}

func (a anyStruct) Not() Matcher {
	return Not(a)
}

func (a anyStruct) Pointer() Matcher {
	return Ref(a)
}

func StructOf(fields map[string]any, opts ...func(m *structs.MatcherOptions)) Matcher {
	o := structs.MatcherOptions{
		Contains: true,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return structFieldsMatcher{fields: fields, options: o}
}

type structFieldsMatcher struct {
	fields  map[string]any
	options structs.MatcherOptions
}

var _ Matcher = structFieldsMatcher{}

func (m structFieldsMatcher) Match(v any) bool {
	if v == nil {
		return false
	}

	s := MayStruct(v)
	if !s.IsStruct() {
		return false
	}

	fields := reflect.VisibleFields(*s.t)
	if !m.options.Contains && len(m.fields) != len(fields) {
		return false
	}

	for k, v := range m.fields {
		f := s.v.FieldByName(k)
		if !f.IsValid() {
			return false
		}

		if !equal(v, f.Interface()) {
			return false
		}
	}

	return true
}

func (a structFieldsMatcher) Not() Matcher {
	return Not(a)
}

func (a structFieldsMatcher) Pointer() Matcher {
	return Ref(a)
}

func MayStruct(raw any) *mayStruct {
	v := reflect.ValueOf(raw)
	t := v.Type()
	return &mayStruct{raw: raw, v: &v, t: &t}
}

type mayStruct struct {
	raw any
	t   *reflect.Type
	v   *reflect.Value
}

func (m mayStruct) IsStruct() bool {
	if m.t == nil {
		return false
	}

	return m.v.Kind() == reflect.Struct
}
