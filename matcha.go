package matcha

import "github.com/version-1/go-matcha/matcher"

func Equal(expect, target any) bool {
	switch m := expect.(type) {
	case matcher.Matcher:
		return m.Match(target)
	default:
		return expect == target
	}
}
