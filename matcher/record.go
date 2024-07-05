package matcher

import (
	"fmt"
	"reflect"
	"strings"
)

type Record struct {
	Matcher  Matcher
	Root     Matcher
	Code     RecordCode
	Key      string
	Expect   any
	Actual   any
	Parent   *Record
	Children []Record
	depth    int
}

func (r *Record) SetChildren(list []Record) {
	r.Children = list
	for i := range r.Children {
		if r.Root == nil {
			r.Children[i].Root = r.Matcher
		} else {
			r.Children[i].Root = r.Root
		}
		r.Children[i].Parent = r
		r.Children[i].depth = r.depth + 1
	}
}

func (r Record) Path() string {
	path := []string{}
	n := r.Parent
	for n != nil {
		path = append([]string{n.Key}, path...)
		n = n.Parent
	}
	path = append(path, r.Key)

	return strings.Join(path, " > ")
}

func (r Record) Error() string {
	return r.String()
}

var padding int = 4

func isSliceOfMatcher(m Matcher) bool {
	_, ok := m.(*sliceOfMatcher)
	return ok
}

func isStructOfMatcher(m Matcher) bool {
	_, ok := m.(*structOfMatcher)
	return ok
}

type recordPrinter struct {
	indent string
}

func (r Record) String() string {
	indent := strings.Repeat(" ", (r.depth+1)*padding)
	chIndent := strings.Repeat(indent, 2)
	keyName := "Field"
	isSliceMatcher := isSliceOfMatcher(r.Matcher)
	if isSliceMatcher {
		keyName = "Slice"
	}

	switch r.Code {
	case RecordCodeTargetIsNil:
		return fmt.Sprintf("%sTarget is nil. expect %T but got nil", indent, r.Matcher)
	case RecordCodeUnmatchLength:
		// TODO: diff fields and print
		return fmt.Sprintf("%s%s length is unmatched. expect %d but got %d", indent, keyName, r.Expect, r.Actual)
	case RecordCodeUnexpectedType:
		return fmt.Sprintf("%sTarget is unexpected type. expect %s but got %T", indent, r.Expect, r.Actual)
	case RecordCodeNotFound:
		if isSliceMatcher {
			return fmt.Sprintf("%sIndex: %s is not found.", indent, r.Path())
		}
		return fmt.Sprintf("%s%s is not found. field: %s", indent, keyName, r.Path())
	case RecordCodeNotEqual:
		v := ExtractIfPossible(r.Expect)
		som, ok := v.(*structOfMatcher)
		if ok {
			msgfmt := "%sField ( %s ) didn't match.\n\n%sexpect: %#v\n\n%sgot: %#v"

			msg := fmt.Sprintf(msgfmt, indent, r.Path(), chIndent, som.fields, chIndent, r.Actual)
			msg += "\n\n"
			for _, c := range r.Children {
				msg += c.String()
			}

			return msg
		}

		slm, ok := v.(*sliceOfMatcher)
		if ok {
			msgfmt := "%sIndex ( %s ) didn't match.\n\n%sexpect: %v\n\n%sgot: %v"

			msg := fmt.Sprintf(msgfmt, indent, r.Path(), chIndent, slm.elements, chIndent, r.Actual)
			msg += "\n\n"
			for _, c := range r.Children {
				msg += c.String()
			}

			return msg
		}

		expectType := reflect.TypeOf(r.Expect)
		if isSlice(expectType) {
			msgfmt := "%sIndex ( %s ) didn't match.\n\n%sexpect: %v\n\n%sgot: %v"

			msg := fmt.Sprintf(msgfmt, indent, r.Path(), chIndent, r.Expect, chIndent, r.Actual)
			msg += "\n\n"
			for _, c := range r.Children {
				msg += c.String()
			}

			return msg
		}

		if r.Root != nil {
			if isSliceOfMatcher(r.Root) {
				return fmt.Sprintf("%sIndex ( %s ) didn't match.\n\n%sexpect: %v\n\n%sgot: %v", indent, r.Path(), chIndent, r.Expect, chIndent, r.Actual)
			}
		}

		return fmt.Sprintf("%sField ( %s ) didn't match.\n\n%sexpect: %s\n\n%sgot: %s", indent, r.Path(), chIndent, r.Expect, chIndent, r.Actual)
	default:
		return fmt.Sprintf("%sField ( %s ) didn't match.\n\n%sgot %s error", indent, r.Path(), chIndent, r.Code)
	}
}

type RecordCode string

const (
	RecordCodeTargetIsNil    RecordCode = "target_is_nil"
	RecordCodeUnexpectedType RecordCode = "unexpected_type"
	RecordCodeNotFound       RecordCode = "not_found"
	RecordCodeNotEqual       RecordCode = "not_equal"
	RecordCodeUnmatchLength  RecordCode = "unmatch_length"
)

type Recorder interface {
	Title() string
	Records() []Record
}

var _ Recorder = &RefMatcher{}
var _ Recorder = &structOfMatcher{}

func recordNotEqual(m Matcher, key string, expect, actual any) Record {
	r := Record{
		Matcher: m,
		Root:    m,
		Key:     key,
		Expect:  expect,
		Actual:  actual,
		Code:    RecordCodeNotEqual,
	}

	rr, ok := expect.(Recorder)
	if ok {
		r.SetChildren(rr.Records())
	}

	return r
}

func recordUnmatchLength(m Matcher, expect, actual int) Record {
	return Record{
		Matcher: m,
		Expect:  expect,
		Actual:  actual,
		Code:    RecordCodeUnmatchLength,
	}
}

func recordNotFound(m Matcher, key string) Record {
	return Record{
		Matcher: m,
		Key:     key,
		Code:    RecordCodeNotFound,
	}
}

func recordTargetIsNil(m Matcher, actual any) Record {
	return Record{
		Matcher: m,
		Code:    RecordCodeTargetIsNil,
		Actual:  actual,
	}
}

func recordUnexpectedType(m Matcher, expect, actual any) Record {
	return Record{
		Matcher: m,
		Expect:  expect,
		Actual:  actual,
		Code:    RecordCodeUnexpectedType,
	}
}
