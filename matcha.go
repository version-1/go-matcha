package matcha

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"testing"

	"github.com/version-1/go-matcha/matcher"
)

func Equal(expect, target any) bool {
	return matcher.Equal(expect, target)
}

func Test(t *testing.T, expect, target any) {
	Equal(expect, target)

	res := NewTesting(t, expect)
	res.Error()
}

func Records(r any) []matcher.Record {
	v, ok := r.(matcher.Recorder)
	if !ok {
		return []matcher.Record{}
	}
	return v.Records()
}

type Testing struct {
	t *testing.T
	r matcher.Recorder
}

func NewTesting(t *testing.T, mayRecords any) *Testing {
	tt := &Testing{t: t}
	v, ok := mayRecords.(matcher.Recorder)
	if ok {
		tt.r = v
	}

	return tt
}

func (t *Testing) Records() []matcher.Record {
	keys := []string{}
	for _, r := range t.r.Records() {
		keys = append(keys, r.Key)
	}

	sort.Strings(keys)

	res := []matcher.Record{}
	for _, k := range keys {
		for _, r := range t.r.Records() {
			if r.Key == k {
				res = append(res, r)
			}
		}
	}

	return res
}

func (t *Testing) Error() {
	t.PrintResult()
	t.t.FailNow()
}

func (t *Testing) PrintResult() {
	if len(t.r.Records()) == 0 {
		return
	}

	msg := []string{t.r.Title()}
	for _, r := range t.r.Records() {
		msg = append(msg, r.Error())
	}

	message := fmt.Sprintf("\n\n\n %s \n\n\n", strings.Join(msg, "\n\n"))
	log.Println(message)
}
