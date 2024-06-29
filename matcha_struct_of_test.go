package matcha

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/version-1/go-matcha/matcher"
	"github.com/version-1/go-matcha/matcher/structs"
)

type post struct {
	ID          uuid.UUID
	Title       string
	Content     string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type group struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type user struct {
	ID        uuid.UUID
	GroupID   uuid.UUID
	Name      string
	Age       int
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time

	Group *group
	Posts []post
}

func TestStructOfEqual(t *testing.T) {
	uid := uuid.New()
	tests := []struct {
		name   string
		expect any
		target any
		ans    bool
	}{
		{"struct matcher: match", matcher.StructOf(map[string]any{
			"ID":        uuid.Nil,
			"GroupID":   uuid.Nil,
			"Name":      "",
			"Age":       0,
			"Status":    "",
			"CreatedAt": time.Time{},
			"UpdatedAt": time.Time{},
		}), user{}, true},
		{"struct matcher: match 2", matcher.StructOf(map[string]any{
			"ID":        uid,
			"GroupID":   uuid.Nil,
			"Name":      "John Doe",
			"Age":       0,
			"Status":    "",
			"CreatedAt": time.Time{},
			"UpdatedAt": time.Time{},
		}), user{
			ID:   uid,
			Name: "John Doe",
		}, true},
		{"struct matcher: match 3", matcher.StructOf(map[string]any{
			"ID":        uid,
			"GroupID":   uuid.Nil,
			"Name":      matcher.BeString(),
			"Age":       matcher.BeInt(),
			"Status":    "",
			"CreatedAt": time.Time{},
			"UpdatedAt": time.Time{},
		}), user{
			ID:   uid,
			Name: "John Doe",
			Age:  25,
		}, true},
		{"struct matcher: nested matcher", matcher.StructOf(map[string]any{
			"ID":        uid,
			"GroupID":   uuid.Nil,
			"Name":      matcher.BeString(),
			"Age":       matcher.BeInt(),
			"Status":    "",
			"CreatedAt": time.Time{},
			"UpdatedAt": time.Time{},
		}), user{
			ID:   uid,
			Name: "John Doe",
			Age:  25,
		}, true},
		{"struct matcher: not match", matcher.StructOf(map[string]any{
			"ID":        uid,
			"GroupID":   uuid.Nil,
			"Name":      "Wrong Name",
			"Age":       0,
			"Status":    "",
			"CreatedAt": time.Time{},
			"UpdatedAt": time.Time{},
		}), user{
			ID:   uid,
			Name: "John Doe",
		}, false},
		{"struct matcher, contains fields: match", matcher.StructOf(map[string]any{
			"ID":      uuid.Nil,
			"GroupID": uuid.Nil,
			"Name":    "",
		}), user{}, true},
		{"struct matcher, contains fields: not match", matcher.StructOf(map[string]any{
			"ID":             uuid.Nil,
			"GroupID":        uuid.Nil,
			"WrongFieldName": "",
		}), user{}, false},
		{"struct matcher: not match", matcher.StructOf(map[string]any{
			"ID":         uuid.Nil,
			"GroupID":    uuid.Nil,
			"Name":       "",
			"Age":        0,
			"Status":     "",
			"CreatedAt":  time.Time{},
			"UpdatedAt":  time.Time{},
			"WrongField": "",
		}), user{}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Equal: %s", tt.name), func(t *testing.T) {
			if Equal(tt.expect, tt.target) != tt.ans {
				t.Errorf("Equal(%v, %v) should return %v", tt.expect, tt.target, tt.ans)
			}
		})
	}
}

func TestStructOfNotMatch(t *testing.T) {
	uid := uuid.New()
	tests := []struct {
		name   string
		expect any
		target any
		ans    []matcher.Record
		assert func(expect, target any, ans []matcher.Record)
	}{
		{
			name: "target is nil",
			expect: matcher.StructOf(matcher.StructMap{
				"ID": uuid.Nil,
			}),
			target: nil,
			ans: []matcher.Record{
				{
					Code:   matcher.RecordCodeTargetIsNil,
					Actual: nil,
				},
			},
			assert: func(expect, target any, ans []matcher.Record) {
				Equal(expect, target)
				test := NewTesting(t, expect)
				records := test.Records()

				for i, r := range records {
					if r.Code != ans[i].Code {
						t.Errorf("r.Code should be %s, got %s", ans[i].Code, r.Code)
					}

					if r.Actual != ans[i].Actual {
						t.Errorf("r.Actual should be %s, got %s", ans[i].Actual, r.Actual)
					}
				}
			},
		},
		{
			name: "target is not struct",
			expect: matcher.StructOf(matcher.StructMap{
				"ID": uuid.Nil,
			}),
			target: 1,
			ans: []matcher.Record{
				{
					Code:   matcher.RecordCodeUnexpectedType,
					Expect: "Struct",
					Actual: 1,
				},
			},
			assert: func(expect, target any, ans []matcher.Record) {
				Equal(expect, target)
				test := NewTesting(t, expect)
				records := test.Records()

				for i, r := range records {
					if r.Code != ans[i].Code {
						t.Errorf("r.Code should be %s, got %s", ans[i].Code, r.Code)
					}

					if r.Actual != ans[i].Actual {
						t.Errorf("r.Actual should be %s, got %s", ans[i].Actual, r.Actual)
					}
				}
			},
		},
		{
			name: "wrong field error",
			expect: matcher.StructOf(matcher.StructMap{
				"WrongField": uuid.Nil,
			}),
			target: user{
				ID: uid,
			},
			ans: []matcher.Record{
				{
					Key:  "WrongField",
					Code: matcher.RecordCodeNotFound,
				},
			},
			assert: func(expect, target any, ans []matcher.Record) {
				Equal(expect, target)
				test := NewTesting(t, expect)
				records := test.Records()

				for i, r := range records {
					if r.Key != ans[i].Key {
						t.Errorf("r.Key should be %s, got %s", ans[i].Key, r.Key)
					}

					if r.Code != ans[i].Code {
						t.Errorf("r.Code should be %s, got %s", ans[i].Code, r.Code)
					}
				}
			},
		},
		{
			name: "unmatch length error",
			expect: matcher.StructOf(matcher.StructMap{
				"ID": uuid.Nil,
			}, structs.WithContains(false)),
			target: user{
				ID: uid,
			},
			ans: []matcher.Record{
				{
					Key:    "",
					Expect: 1,
					Actual: 9,
					Code:   matcher.RecordCodeUnmatchLength,
				},
			},
			assert: func(expect, target any, ans []matcher.Record) {
				Equal(expect, target)
				test := NewTesting(t, expect)
				records := test.Records()

				for i, r := range records {
					if r.Expect != ans[i].Expect {
						t.Errorf("r.Expect should be %s, got %s", ans[i].Expect, r.Expect)
					}

					if r.Actual != ans[i].Actual {
						t.Errorf("r.Actual should be %s, got %s", ans[i].Actual, r.Actual)
					}

					if r.Code != ans[i].Code {
						t.Errorf("r.Code should be %s, got %s", ans[i].Code, r.Code)
					}
				}
			},
		},
		{
			name: "nested struct error",
			expect: matcher.StructOf(matcher.StructMap{
				"ID":        uuid.Nil,
				"GroupID":   uuid.Nil,
				"Name":      "",
				"Age":       0,
				"Status":    "",
				"CreatedAt": time.Time{},
				"UpdatedAt": time.Time{},
				"Group": matcher.StructOf(matcher.StructMap{
					"ID":        uuid.Nil,
					"Name":      "",
					"CreatedAt": time.Time{},
					"UpdatedAt": time.Time{},
				}).Pointer(),
			}),
			target: user{
				ID:  uid,
				Age: 24,
				Group: &group{
					ID: uuid.New(),
				},
			},
			ans: []matcher.Record{
				{
					Key:  "Age",
					Code: matcher.RecordCodeNotEqual,
				},
				{
					Key:  "Group",
					Code: matcher.RecordCodeNotEqual,
					Children: []matcher.Record{
						{
							Key:  "ID",
							Code: matcher.RecordCodeNotEqual,
						},
					},
				},
				{
					Key:  "ID",
					Code: matcher.RecordCodeNotEqual,
				},
			},
			assert: func(expect, target any, ans []matcher.Record) {
				Equal(expect, target)
				test := NewTesting(t, expect)
				records := test.Records()

				for i, r := range records {
					if r.Key != ans[i].Key {
						t.Errorf("r.Key should be %s, got %s", ans[i].Key, r.Key)
					}

					if r.Code != ans[i].Code {
						t.Errorf("r.Code should be %s, got %s", ans[i].Code, r.Code)
					}

					if i == 1 {
						for j, child := range r.Children {
							if child.Key != ans[i].Children[j].Key {
								t.Errorf("child.Key should be %s, got %s", ans[i].Children[j].Key, child.Key)
							}

							if child.Code != ans[i].Children[j].Code {
								t.Errorf("child.Code should be %s, got %s", ans[i].Children[j].Code, child.Code)
							}
						}
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Equal: %s", tt.name), func(t *testing.T) {
			tt.assert(tt.expect, tt.target, tt.ans)
		})
	}
}
