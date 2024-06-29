package matcher

import (
	"reflect"

	"github.com/version-1/go-matcha/matcher/structs"
)

func BeStruct() *anyStruct {
	return &anyStruct{}
}

type anyStruct struct {
	options MatcherOptions
}

func (a anyStruct) Match(v any) bool {
	if v == nil {
		return false
	}

	if !a.options.AllowZero && isZero(v) {
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

func (a anyStruct) AllowZero() Matcher {
	a.options.AllowZero = true
	return a
}

type StructMap map[string]any

func StructOf(fields StructMap, opts ...func(m *structs.MatcherOptions)) Matcher {
	o := structs.MatcherOptions{
		Contains: true,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &structOfMatcher{fields: fields, options: o}
}

type structOfMatcher struct {
	fields  StructMap
	options structs.MatcherOptions
	records []Record
}

func (m structOfMatcher) Title() string {
	return "StructOfMatcher got errors."
}

func (m structOfMatcher) Records() []Record {
	return m.records
}

func (m *structOfMatcher) Match(v any) bool {
	if v == nil {
		r := recordTargetIsNil(m, v)
		m.records = append(m.records, r)
		return false
	}

	s := MayStruct(v)
	if !s.IsStruct() {
		r := recordUnexpectedType(m, "Struct", v)
		m.records = append(m.records, r)
		return false
	}

	fields := reflect.VisibleFields(*s.t)
	if !m.options.Contains && len(m.fields) != len(fields) {
		r := recordUnmatchLength(m, len(m.fields), len(fields))
		m.records = append(m.records, r)
		return false
	}

	for k, v := range m.fields {
		f := s.v.FieldByName(k)
		if !f.IsValid() {
			r := recordNotFound(m, k)
			m.records = append(m.records, r)
			continue
		}

		if !Equal(v, f.Interface()) {
			r := recordNotEqual(m, k, v, f.Interface())
			m.records = append(m.records, r)

			continue
		}
	}

	return len(m.records) == 0
}

func (m structOfMatcher) Not() Matcher {
	return Not(&m)
}

func (m structOfMatcher) Pointer() Matcher {
	return Ref(&m)
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
