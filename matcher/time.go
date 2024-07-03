package matcher

import "time"

func BeTime() *anyTime {
	return &anyTime{}
}

type anyTime struct {
	options MatcherOptions
}

var _ Matcher = anyTime{}

func (m anyTime) Match(v any) bool {
	if !m.options.AllowZero && isZero(v) {
		return false
	}

	return typeMatch[time.Time](v)
}

func (m anyTime) Not() Matcher {
	return Not(m)
}

func (m anyTime) Pointer() Matcher {
	return Ref(m)
}

func (m anyTime) AllowZero() Matcher {
	m.options.AllowZero = true
	return m
}
