package matcher

import (
	"fmt"
	"strings"
)

type Record struct {
	Matcher  Matcher
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

func (r Record) String() string {
	indent := strings.Repeat(" ", (r.depth+1)*padding)
	chIndent := strings.Repeat(indent, 2)

	if r.Code == RecordCodeTargetIsNil {
		return fmt.Sprintf("%starget is nil. expect %T but got nil", indent, r.Matcher)
	}

	if r.Code == RecordCodeWrongFieldCount {
		// TODO: diff fields and print
		return fmt.Sprintf("%sField count is unmatched. expect %d but got %d", indent, r.Expect, r.Actual)
	}

	if r.Code == RecordCodeNotStruct {
		return fmt.Sprintf("%sTarget is not struct. expect struct but got %T", indent, r.Actual)
	}

	if r.Code == RecordCodeFieldNotFound {
		return fmt.Sprintf("%sField is not found. field: %s", indent, r.Path())
	}

	if r.Code == RecordCodeNotEqual {
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

		return fmt.Sprintf("%sField ( %s ) didn't match.\n\n%sexpect: %s\n\n%sgot: %s", indent, r.Path(), chIndent, r.Expect, chIndent, r.Actual)
	}

	return string(r.Code)
}

type RecordCode string

const (
	RecordCodeTargetIsNil     RecordCode = "target_is_nil"
	RecordCodeWrongFieldCount RecordCode = "wrong_field_count"
	RecordCodeNotStruct       RecordCode = "not_struct"
	RecordCodeFieldNotFound   RecordCode = "field_not_found"
	RecordCodeNotEqual        RecordCode = "not_equal"
)

type Recorder interface {
	Title() string
	Records() []Record
}

var _ Recorder = &RefMatcher{}
var _ Recorder = &structOfMatcher{}
