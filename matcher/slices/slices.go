package slices

type MatcherOptions struct {
	Order    bool
	Contains bool
}

func WithPersistOrder(v bool) func(*MatcherOptions) {
	return func(o *MatcherOptions) {
		o.Order = v
	}
}

func WithContains(v bool) func(*MatcherOptions) {
	return func(o *MatcherOptions) {
		o.Contains = v
	}
}
