package matcher

import "reflect"

func BeAny() *beAny {
	return &beAny{}
}

type beAny struct{}

var _ Matcher = beAny{}

func (b beAny) Match(v any) bool {
	return true
}

func (b beAny) Not() Matcher {
	return Not(b)
}

func (b beAny) Pointer() Matcher {
	return Ref(b)
}

func BeZero() *beZero {
	return &beZero{}
}

type beZero struct{}

var _ Matcher = beZero{}

func (b beZero) Match(v any) bool {
	if v == nil {
		return true
	}

	return v == reflect.Zero(reflect.TypeOf(v)).Interface()
}

func (b beZero) Not() Matcher {
	return Not(b)
}

func (b beZero) Pointer() Matcher {
	return Ref(b)
}
