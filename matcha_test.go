package matcha

import (
	"testing"

	"github.com/version-1/go-matcha/internal/pointer"
	"github.com/version-1/go-matcha/matcher"
)

type dummy struct {
	a int
}

func TestZeroValueEqual(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		// zero value
		{"any matcher with zero", matcher.BeAny(), 0, false},
		{"any matcher with empty struct", matcher.BeAny(), dummy{}, false},
		{"any matcher with string array", matcher.BeAny(), []string{}, false},
		{"any matcher with int array", matcher.BeAny(), []int{}, false},
		{"string matcher with zero string", matcher.BeString(), "", false},
		{"int matcher with zero int", matcher.BeInt(), 0, false},
		{"slice matcher with zero string slice", matcher.BeSlice(), []string{}, false},
		{"slice matcher with zero any slice", matcher.BeSlice(), []any{}, false},
		{"struct matcher with zero struct", matcher.BeStruct(), dummy{}, false},
		// (allow zero)
		{"any matcher with zero", matcher.BeAny().AllowZero(), 0, true},
		{"any matcher with empty struct", matcher.BeAny().AllowZero(), dummy{}, true},
		{"any matcher with string array", matcher.BeAny().AllowZero(), []string{}, true},
		{"any matcher with int array", matcher.BeAny().AllowZero(), []int{}, true},
		{"string matcher with zero string", matcher.BeString().AllowZero(), "", true},
		{"int matcher with zero int", matcher.BeInt().AllowZero(), 0, true},
		{"slice matcher with zero string slice", matcher.BeSlice().AllowZero(), []string{}, true},
		{"slice matcher with zero any slice", matcher.BeSlice().AllowZero(), []any{}, true},
		{"struct matcher with zero struct", matcher.BeStruct().AllowZero(), dummy{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
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
		{"bool equal", true, true, true},
		{"bool not equal", true, false, false},
		{"slice equal", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"slice not equal", []int{1, 2}, []int{1, 2, 3}, false},
		{"struct equal", dummy{}, dummy{}, true},
		{"slice not equal", dummy{}, dummy{1}, false},
		// any
		{"any matcher with num", matcher.BeAny(), 1, true},
		{"any matcher with string", matcher.BeAny(), "abca", true},
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

func TestAnyStructEqual(t *testing.T) {
	var dNil *dummy
	var dValid *dummy = &dummy{1}

	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		{"struct matcher: match", matcher.BeStruct(), dummy{}, false},
		{"struct matcher: match 2", matcher.BeStruct(), dummy{1}, true},
		{"struct matcher: not match", matcher.BeStruct(), []string{"a", "b", "c"}, false},
		{"struct matcher: not match 2", matcher.BeStruct(), nil, false},
		{"pointer struct matcher: match", matcher.BeStruct().Pointer(), dValid, true},
		{"pointer struct matcher: not match", matcher.BeStruct().Pointer(), dNil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
}
