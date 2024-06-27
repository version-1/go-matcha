package matcha

import (
	"testing"

	"github.com/version-1/go-matcha/internal/pointer"
	"github.com/version-1/go-matcha/matcher"
)

type dummy struct {
	a int
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		// primitive
		{"num equal", 1, 1, true},
		{"num not equal", 1, 2, false},
		{"string equal", "abc", "abc", true},
		{"string not qual", "abca", "abc", false},
		{"bool equal", true, true, true},
		{"bool not qual", true, false, false},
		// any
		{"any matcher with num", matcher.BeAny(), 1, true},
		{"any matcher with string", matcher.BeAny(), "abca", true},
		// string
		{"string matcher with string", matcher.BeString(), "123", true},
		{"string matcher with not string", matcher.BeString(), 123, false},
		{"string matcher with string ref", matcher.BeString(), pointer.Ref("123"), false},
		{"string ref matcher with string", matcher.BeString().Pointer(), "123", false},
		{"string ref matcher with not string", matcher.BeString().Pointer(), 123, false},
		{"string ref matcher with not string pointer", matcher.BeString().Pointer(), pointer.Ref(123), false},
		{"string ref matcher with string ref", matcher.BeString().Pointer(), pointer.Ref("123"), true},
		// int
		{"int matcher with int", matcher.BeInt(), 123, true},
		{"int matcher with not int", matcher.BeInt(), "123", false},
		{"int matcher with not int ref", matcher.BeInt(), pointer.Ref("123"), false},
		{"int matcher with int ref", matcher.BeInt(), pointer.Ref(123), false},
		{"int ref matcher with int", matcher.BeInt().Pointer(), 123, false},
		{"int ref matcher with not int", matcher.BeInt().Pointer(), "123", false},
		{"int ref matcher with not int ref", matcher.BeInt().Pointer(), pointer.Ref("123"), false},
		{"int ref matcher with int ref", matcher.BeInt().Pointer(), pointer.Ref(123), true},
		// bool
		{"bool matcher with bool", matcher.BeBool(), false, true},
		{"bool matcher with bool", matcher.BeBool(), true, true},
		{"bool matcher with not bool", matcher.BeBool(), "", false},
		{"bool matcher with bool ref", matcher.BeBool(), pointer.Ref(true), false},
		{"bool ref matcher with bool", matcher.BeBool().Pointer(), true, false},
		{"bool ref matcher with not bool", matcher.BeBool().Pointer(), true, false},
		{"bool ref matcher with bool ref", matcher.BeBool().Pointer(), pointer.Ref(true), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		{"not zero matcher with int zero", matcher.BeZero().Not(), 0, false},
		{"not zero matcher with int non zero", matcher.BeZero().Not(), 1, true},
		{"not zero matcher with string zero", matcher.BeZero().Not(), "", false},
		{"not zero matcher with string non zero", matcher.BeZero().Not(), "a", true},
		{"not zero matcher with nil", matcher.BeZero().Not(), nil, false},
		{"not zero matcher with zero struct", matcher.BeZero().Not(), dummy{}, false},
		{"not zero matcher with non zero struct", matcher.BeZero().Not(), dummy{1}, true},
		// string
		{"not string matcher with string", matcher.BeString().Not(), "123", false},
		{"not string matcher with not string", matcher.BeString().Not(), 123, true},
		{"not string matcher with string ref", matcher.BeString().Not(), pointer.Ref("123"), true},
		{"not string ref matcher with string", matcher.BeString().Pointer().Not(), "123", true},
		{"not string ref matcher with not string", matcher.BeString().Pointer().Not(), 123, true},
		{"not string ref matcher with not string pointer", matcher.BeString().Pointer().Not(), pointer.Ref(123), true},
		{"not string ref matcher with string ref", matcher.BeString().Pointer().Not(), pointer.Ref("123"), false},
		// int
		{"not int matcher with int", matcher.BeInt().Not(), 123, false},
		{"not int matcher with not int", matcher.BeInt().Not(), "123", true},
		{"not int matcher with not int ref", matcher.BeInt().Not(), pointer.Ref("123"), true},
		{"not int matcher with int ref", matcher.BeInt().Not(), pointer.Ref(123), true},
		{"not int ref matcher with int", matcher.BeInt().Pointer().Not(), 123, true},
		{"not int ref matcher with not int", matcher.BeInt().Pointer().Not(), "123", true},
		{"not int ref matcher with not int ref", matcher.BeInt().Pointer().Not(), pointer.Ref("123"), true},
		{"not int ref matcher with int ref", matcher.BeInt().Pointer().Not(), pointer.Ref(123), false},
		// bool
		{"not bool matcher with bool", matcher.BeBool().Not(), false, false},
		{"not bool matcher with bool", matcher.BeBool().Not(), true, false},
		{"not bool matcher with not bool", matcher.BeBool().Not(), "", true},
		{"not bool matcher with bool ref", matcher.BeBool().Not(), pointer.Ref(true), true},
		{"not bool ref matcher with bool", matcher.BeBool().Pointer().Not(), true, true},
		{"not bool ref matcher with not bool", matcher.BeBool().Pointer().Not(), "1", true},
		{"not bool ref matcher with bool ref", matcher.BeBool().Pointer().Not(), pointer.Ref(true), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
}

func TestZeroEqual(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		{"zero matcher with int zero", matcher.BeZero(), 0, true},
		{"zero matcher with int non zero", matcher.BeZero(), 1, false},
		{"zero matcher with string zero", matcher.BeZero(), "", true},
		{"zero matcher with string non zero", matcher.BeZero(), "a", false},
		{"zero matcher with nil", matcher.BeZero(), nil, true},
		{"zero matcher with zero struct", matcher.BeZero(), dummy{}, true},
		{"zero matcher with non zero struct", matcher.BeZero(), dummy{1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
}

func TestSliceEqual(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		{"array matcher with string array", matcher.BeSlice(), []string{"a", "b", "c"}, true},
		{"array matcher with int array", matcher.BeSlice(), []int{1, 2, 3}, true},
		{"array matcher with not array", matcher.BeSlice(), 1, false},
		{"array ref matcher with string array ref", matcher.BeSlice().Pointer(), &[]string{"a", "b", "c"}, true},
		{"array ref matcher with string array", matcher.BeSlice().Pointer(), []string{"a", "b", "c"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
}

func TestSliceOfEqual(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		// with primitive
		{"slice of matcher: match", matcher.SliceOf([]any{
			"a", "b", "c",
		}), []string{"a", "b", "c"}, true},
		{"slice of matcher: not match", matcher.SliceOf([]any{
			"a", "b", "c",
		}), []string{"a", "b", "c", "b"}, false},
		{"slice of matcher: not match", matcher.SliceOf([]any{
			1, 2, 3,
		}), []string{"a", "b", "c", "b"}, false},
		{"slice of matcher: order not match", matcher.SliceOf([]any{
			1, 2, 3,
		}), []int{3, 2, 1}, false},
		// contains
		{
			"slice of matcher with contains: match",
			matcher.SliceOf(
				[]any{1, 2, 3},
				matcher.WithSliceOfContains(true),
			),
			[]int{1, 2, 3, 4},
			true,
		},
		{
			"slice of matcher with contains: not match",
			matcher.SliceOf(
				[]any{1, 2, 3},
				matcher.WithSliceOfContains(true),
			),
			[]int{1, 2, 4, 5},
			false,
		},
		// with matcher
		{
			"slice of matcher with matcher: match",
			matcher.SliceOf(
				[]any{matcher.BeAny(), matcher.BeAny(), matcher.BeAny()},
			),
			[]int{1, 2, 4},
			true,
		},
		{
			"slice of matcher with matcher: match 2",
			matcher.SliceOf(
				[]any{
					matcher.BeInt(),
					matcher.BeString(),
					matcher.BeBool(),
					matcher.BeInt().Pointer(),
					matcher.BeString().Pointer(),
					matcher.BeBool().Pointer(),
				},
			),
			[]any{1, "abc", true, pointer.Ref(1), pointer.Ref("abc"), pointer.Ref(true)},
			true,
		},
		{
			"slice of matcher with matcher: not match",
			matcher.SliceOf(
				[]any{
					matcher.BeInt(),
					matcher.BeString().Not(),
					matcher.BeBool(),
				},
			),
			[]any{1, "abc", true},
			false,
		},
		{
			"slice of matcher with nest slice: match",
			matcher.SliceOf(
				[]any{
					matcher.SliceOf([]any{1, 2, 3}),
					matcher.SliceOf([]any{"abc", "cde", "fgh"}),
					true,
				},
			),
			[]any{[]int{1, 2, 3}, []string{"abc", "cde", "fgh"}, true},
			true,
		},
		{
			"slice of matcher with nest slice: not match",
			matcher.SliceOf(
				[]any{
					matcher.SliceOf([]any{1, 2, 3}),
					matcher.SliceOf([]any{"abc", "cde", "fgh"}),
					true,
				},
			),
			[]any{[]int{1, 2, 3}, []string{"abc", "123cde", "fgh"}, true},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
}

func TestSliceLenEqual(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		{"slice length matcher: match", matcher.SliceLen(3), []string{"a", "b", "c"}, true},
		{"slice length matcher: not match", matcher.SliceLen(3), []string{"a", "b", "c", "d"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
}