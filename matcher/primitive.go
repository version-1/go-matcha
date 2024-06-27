package matcher

// string
func BeString() *anyString {
	return &anyString{}
}

type anyString struct{}

var _ Matcher = anyString{}

func (s anyString) Match(v any) bool {
	return typeMatch[string](v)
}

func (s anyString) Not() Matcher {
	return Not(s)
}

func (s anyString) Pointer() Matcher {
	return Ref(s)
}

// int
func BeInt() *anyInt {
	return &anyInt{}
}

type anyInt struct{}

var _ Matcher = anyInt{}

func (e anyInt) Match(v any) bool {
	return typeMatch[int](v)
}

func (e anyInt) Not() Matcher {
	return Not(e)
}

func (e anyInt) Pointer() Matcher {
	return Ref(e)
}

// bool
type anyBool struct{}

var _ Matcher = anyBool{}

func BeBool() *anyBool {
	return &anyBool{}
}

func (e anyBool) Match(v any) bool {
	return typeMatch[bool](v)
}

func (e anyBool) Not() Matcher {
	return Not(e)
}

func (e anyBool) Pointer() Matcher {
	return Ref(e)
}
