package matcher

import "reflect"

func Equal(expect, target any) bool {
	v, ok := expect.(Matcher)
	if ok {
		return v.Match(target)
	}

	r := reflect.TypeOf(expect)
	if r.Kind() == reflect.Slice || r.Kind() == reflect.Array {
		return reflect.DeepEqual(expect, target)
	}

	return expect == target
}

type MatcherOptions struct {
	AllowZero bool
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

var _ Matcher = &RefMatcher{}
var _ Matcher = &notMatcher{}
var _ Matcher = beZero{}
var _ Matcher = beAny{}
var _ Matcher = anyBool{}
var _ Matcher = anyInt{}
var _ Matcher = anyString{}
var _ Matcher = anySlice{}
var _ Matcher = anyStruct{}
var _ Matcher = sliceLenMatcher{}
var _ Matcher = &sliceOfMatcher{}
var _ Matcher = &structOfMatcher{}
