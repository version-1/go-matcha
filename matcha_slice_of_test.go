package matcha

import (
	"reflect"
	"testing"

	"github.com/version-1/go-matcha/assert"
	"github.com/version-1/go-matcha/internal/pointer"
	"github.com/version-1/go-matcha/matcher"
	"github.com/version-1/go-matcha/matcher/slices"
)

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
				slices.WithContains(true),
			),
			[]int{1, 2, 3, 4},
			true,
		},
		{
			"slice of matcher with contains: not match",
			matcher.SliceOf(
				[]any{1, 2, 3},
				slices.WithContains(true),
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

func TestSliceOfNotEqual(t *testing.T) {
	tests := []struct {
		name   string
		expect any
		target any
		ans    []matcher.Record
		assert func(expect, target any, ans []matcher.Record)
	}{
		{
			name:   "target is nil",
			expect: matcher.SliceOf([]any{}),
			target: nil,
			ans: []matcher.Record{
				{
					Code:   matcher.RecordCodeTargetIsNil,
					Actual: nil,
				},
			},
			assert: func(expect, target any, ans []matcher.Record) {
				Equal(expect, target)
				test := assert.New(t, expect, target)
				records := test.Records()

				for i, r := range records {
					if r.Code != ans[i].Code {
						t.Errorf("r.Code should be %s, got %s", ans[i].Code, r.Code)
					}
				}
			},
		},
		{
			name:   "target is not slice",
			expect: matcher.SliceOf([]any{}),
			target: 1,
			ans: []matcher.Record{
				{
					Code:   matcher.RecordCodeUnexpectedType,
					Actual: 1,
				},
			},
			assert: func(expect, target any, ans []matcher.Record) {
				Equal(expect, target)
				test := assert.New(t, expect, target)
				records := test.Records()

				for i, r := range records {
					if r.Code != ans[i].Code {
						t.Errorf("r.Code should be %s, got %s", ans[i].Code, r.Code)
					}
				}
			},
		},
		{
			name:   "unmatch slice length",
			expect: matcher.SliceOf([]any{1}, slices.WithContains(false)),
			target: []int{1, 2, 3},
			ans: []matcher.Record{
				{
					Code:   matcher.RecordCodeUnmatchLength,
					Expect: 1,
					Actual: 3,
				},
			},
			assert: func(expect, target any, ans []matcher.Record) {
				Equal(expect, target)
				test := assert.New(t, expect, target)
				records := test.Records()

				for i, r := range records {
					if r.Code != ans[i].Code {
						t.Errorf("r.Code should be %s, got %s", ans[i].Code, r.Code)
					}

					if r.Expect != ans[i].Expect {
						t.Errorf("r.Expect should be %s, got %s", ans[i].Expect, r.Expect)
					}

					if r.Actual != ans[i].Actual {
						t.Errorf("r.Actual should be %s, got %s", ans[i].Actual, r.Actual)
					}
				}
			},
		},
		{
			name:   "record not found",
			expect: matcher.SliceOf([]any{1, 2, 3}, slices.WithContains(true)),
			target: []int{1, 2},
			ans: []matcher.Record{
				{
					Code: matcher.RecordCodeNotFound,
					Key:  "2",
				},
			},
			assert: func(expect, target any, ans []matcher.Record) {
				Equal(expect, target)
				test := assert.New(t, expect, target)
				records := test.Records()

				for i, r := range records {
					if r.Code != ans[i].Code {
						t.Errorf("r.Code should be %s, got %s", ans[i].Code, r.Code)
					}

					if r.Key != ans[i].Key {
						t.Errorf("r.Key should be %s, got %s", ans[i].Key, r.Key)
					}
				}
			},
		},
		{
			name:   "record not equal",
			expect: matcher.SliceOf([]any{1, 10, 3}, slices.WithContains(true)),
			target: []int{1, 2, 3},
			ans: []matcher.Record{
				{
					Code:   matcher.RecordCodeNotEqual,
					Key:    "1",
					Expect: 10,
					Actual: 2,
				},
			},
			assert: func(expect, target any, ans []matcher.Record) {
				Equal(expect, target)
				test := assert.New(t, expect, target)
				records := test.Records()

				for i, r := range records {
					if r.Code != ans[i].Code {
						t.Errorf("r.Code should be %s, got %s", ans[i].Code, r.Code)
					}

					if r.Key != ans[i].Key {
						t.Errorf("r.Key should be %s, got %s", ans[i].Key, r.Key)
					}

					if r.Expect != ans[i].Expect {
						t.Errorf("r.Expect should be %s, got %s", ans[i].Expect, r.Expect)
					}

					if r.Actual != ans[i].Actual {
						t.Errorf("r.Actual should be %s, got %s", ans[i].Actual, r.Actual)
					}
				}
			},
		},
		{
			name:   "nested slice",
			expect: matcher.SliceOf([]any{matcher.BeInt(), []int{2, 3}, 3, []int{1, 2, 3}}, slices.WithContains(true)),
			target: []any{1, 2, 3, []int{}},
			ans: []matcher.Record{
				{
					Code:   matcher.RecordCodeNotEqual,
					Key:    "1",
					Expect: []int{2, 3},
					Actual: 2,
				},
				{
					Code:   matcher.RecordCodeNotEqual,
					Key:    "3",
					Expect: []int{1, 2, 3},
					Actual: []int{},
				},
			},
			assert: func(expect, target any, ans []matcher.Record) {
				Equal(expect, target)
				test := assert.New(t, expect, target)
				records := test.Records()
				if len(records) != len(ans) {
					t.Errorf("Length should be %d, got %d", len(ans), len(records))
				}

				for i, r := range records {
					if r.Code != ans[i].Code {
						t.Errorf("r.Code should be %s, got %s", ans[i].Code, r.Code)
					}

					if r.Key != ans[i].Key {
						t.Errorf("r.Key should be %s, got %s", ans[i].Key, r.Key)
					}

					if !reflect.DeepEqual(r.Expect, ans[i].Expect) {
						t.Errorf("r.Expect should be %s, got %s", ans[i].Expect, r.Expect)
					}

					if !reflect.DeepEqual(r.Actual, ans[i].Actual) {
						t.Errorf("r.Actual should be %s, got %s", ans[i].Actual, r.Actual)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assert(tt.expect, tt.target, tt.ans)
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
