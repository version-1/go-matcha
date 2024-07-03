package matcha

import (
	"github.com/version-1/go-matcha/assert"
	"github.com/version-1/go-matcha/matcher"
)

func Equal(expect, target any) bool {
	return matcher.Equal(expect, target)
}

func Test(t assert.Testing, expect any, target any) {
	assertion := Equal(expect, target)
	if assertion {
		return
	}

	res := assert.New(t, expect, target)
	res.Assert()
}

func Records(r any) []matcher.Record {
	v, ok := r.(matcher.Recorder)
	if !ok {
		return []matcher.Record{}
	}
	return v.Records()
}
