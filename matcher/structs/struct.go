package structs

type MatcherOptions struct {
	Contains bool
}

func WithContains(b bool) func(*MatcherOptions) {
	return func(o *MatcherOptions) {
		o.Contains = b
	}
}
