package matcher

import (
	"net/mail"
	"regexp"
)

// string
func BeString() *anyString {
	return &anyString{}
}

type anyString struct {
	options MatcherOptions
}

var _ Matcher = anyString{}

func (m anyString) Match(v any) bool {
	if !m.options.AllowZero && isZero(v) {
		return false
	}

	return typeMatch[string](v)
}

func (m anyString) Not() Matcher {
	return Not(m)
}

func (m anyString) Pointer() Matcher {
	return Ref(m)
}

func (m anyString) AllowZero() Matcher {
	m.options.AllowZero = true
	return m
}

type regExpMatcher struct {
	regexp *regexp.Regexp
}

func RegExp(r string) Matcher {
	m := regexp.MustCompile(r)

	return &regExpMatcher{m}
}

func (m regExpMatcher) Match(v any) bool {
	if v == nil {
		return false
	}

	switch vv := v.(type) {
	case string:
		return m.regexp.MatchString(vv)
	case *string:
		return m.regexp.MatchString(*vv)
	default:
		return false
	}
}

func (m regExpMatcher) Not() Matcher {
	return Not(m)
}

func (m regExpMatcher) Pointer() Matcher {
	return Ref(m)
}

type emailMatcher struct{}

func Email() Matcher {
	return &emailMatcher{}
}

func (m emailMatcher) Match(v any) bool {
	if v == nil {
		return false
	}

	var target string
	switch vv := v.(type) {
	case string:
		target = vv
	case *string:
		target = *vv
	default:
		return false
	}

	_, err := mail.ParseAddress(target)
	if err != nil {
		return false
	}

	return true
}

func (m emailMatcher) Not() Matcher {
	return Not(m)
}

func (m emailMatcher) Pointer() Matcher {
	return Ref(m)
}
