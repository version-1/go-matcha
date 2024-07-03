package assert

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/version-1/go-matcha/matcher"
)

type Testing interface {
	FailNow()
}

type assertion struct {
	t      Testing
	expect any
	target any
	r      matcher.Recorder
}

func New(t Testing, expect, target any) *assertion {
	tt := &assertion{t: t}
	tt.expect = expect
	tt.target = target

	r, ok := expect.(matcher.Recorder)
	if ok {
		tt.r = r
	}

	return tt
}

func (a assertion) Records() []matcher.Record {
	if a.r == nil {
		return []matcher.Record{}
	}

	keys := []string{}

	for _, r := range a.r.Records() {
		keys = append(keys, r.Key)
	}

	sort.Strings(keys)

	res := []matcher.Record{}
	for _, k := range keys {
		for _, r := range a.r.Records() {
			if r.Key == k {
				res = append(res, r)
			}
		}
	}

	return res
}

func (a assertion) Assert() {
	if a.r == nil {
		log.Printf("expect %s but got %s", a.expect, Stringify(a.target))
		a.t.FailNow()
		return
	}

	a.PrintResult()
	if len(a.Records()) > 0 {
		a.t.FailNow()
	}
}

func (a assertion) PrintResult() {
	if len(a.r.Records()) == 0 {
		return
	}

	msg := []string{a.r.Title()}
	for _, r := range a.r.Records() {
		msg = append(msg, r.Error())
	}

	message := fmt.Sprintf("\n\n\n %s \n\n\n", strings.Join(msg, "\n\n"))
	log.Println(message)
}

type stringer interface {
	String() string
}

func Stringify(s any) string {
	if s == nil {
		return ""
	}

	switch vv := s.(type) {
	case string:
		return vv
	case *string:
		return *vv
	case int:
		return strconv.Itoa(vv)
	case *int:
		return strconv.Itoa(*vv)
	case bool:
		return strconv.FormatBool(vv)
	case *bool:
		return strconv.FormatBool(*vv)
	case stringer:
		return vv.String()
	}

	t := reflect.TypeOf(s)

	var vv any
	if t.Kind() == reflect.Ptr {
		vv = t.Elem()
	} else {
		vv = s
	}

	if t.Kind() == reflect.Struct {
		return fmt.Sprintf("%#v", vv)
	}

	if t.Kind() == reflect.Slice {
		return fmt.Sprintf("%v", vv)
	}

	return fmt.Sprintf("%s", s)
}
