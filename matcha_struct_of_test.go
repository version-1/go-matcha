package matcha

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/version-1/go-matcha/matcher"
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
