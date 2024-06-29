package matcher

// int
func BeInt() *anyInt {
	return &anyInt{}
}

type anyInt struct {
	options MatcherOptions
}

func (m anyInt) Match(v any) bool {
	if !m.options.AllowZero && isZero(v) {
		return false
	}

	return typeMatch[int](v)
}

func (m anyInt) Not() Matcher {
	return Not(m)
}

func (m anyInt) Pointer() Matcher {
	return Ref(m)
}

func (m anyInt) AllowZero() Matcher {
	m.options.AllowZero = true
	return m
}

// bool
type anyBool struct{}

func BeBool() *anyBool {
	return &anyBool{}
}

// INFO: bool matcher allows zero by default
func (e anyBool) Match(v any) bool {
	return typeMatch[bool](v)
}

func (e anyBool) Not() Matcher {
	return Not(e)
}

func (e anyBool) Pointer() Matcher {
	return Ref(e)
}
