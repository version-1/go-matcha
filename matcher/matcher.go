package matcher

func equal(a, b any) bool {
	v, ok := a.(Matcher)
	if ok {
		return v.Match(b)
	}

	return a == b
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
