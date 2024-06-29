package matcha

import (
	"testing"

	"github.com/version-1/go-matcha/internal/pointer"
	"github.com/version-1/go-matcha/matcher"
)

func TestStringEqual(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		// primitive
		{"string equal", "abc", "abc", true},
		{"string not qual", "abca", "abc", false},
		{"slice equal", []string{"a", "b", "c"}, []string{"a", "b", "c"}, true},
		{"slice not equal", []string{"a", "b"}, []string{"a", "b", "c"}, false},
		// string
		{"string matcher with string", matcher.BeString(), "123", true},
		{"string matcher with not string", matcher.BeString(), 123, false},
		{"string matcher with string ref", matcher.BeString(), pointer.Ref("123"), false},
		{"string ref matcher with string", matcher.BeString().Pointer(), "123", false},
		{"string ref matcher with not string", matcher.BeString().Pointer(), 123, false},
		{"string ref matcher with not string pointer", matcher.BeString().Pointer(), pointer.Ref(123), false},
		{"string ref matcher with string ref", matcher.BeString().Pointer(), pointer.Ref("123"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
}

func TestRegexpEqual(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		// string
		{"regexp matcher with string", matcher.RegExp("^[0-9]+$"), "123", true},
		{"not matcher with string", matcher.RegExp("1234"), "123", false},
		{"regexp matcher with nil", matcher.RegExp("^[0-9]+$"), nil, false},
		{"regexp matcher with not string", matcher.RegExp("^[0-9]+$"), 123, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
}

func TestEmailEqual(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		// string
		{"email matcher with email string", matcher.Email(), "hoge@example.com", true},
		{"email matcher with email string 2", matcher.Email(), "Hoge <hoge@example.com>", true},
		{"not matcher with not-email string", matcher.Email(), "hoge", false},
		{"not matcher with not-email string 2", matcher.Email(), "example.com", false},
		{"not matcher with not string", matcher.Email(), 123, false},
		{"not matcher with nil", matcher.Email(), nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
}
