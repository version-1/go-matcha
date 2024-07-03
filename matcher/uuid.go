package matcher

import "github.com/google/uuid"

func BeUUID() *anyUUID {
	return &anyUUID{}
}

type anyUUID struct {
	options MatcherOptions
}

var _ Matcher = anyUUID{}

func (m anyUUID) Match(v any) bool {
	if !m.options.AllowZero && (uuid.Nil == v || isZero(v)) {
		return false
	}

	return typeMatch[uuid.UUID](v)
}

func (m anyUUID) Not() Matcher {
	return Not(m)
}

func (m anyUUID) Pointer() Matcher {
	return Ref(m)
}

func (m anyUUID) AllowZero() Matcher {
	m.options.AllowZero = true
	return m
}
