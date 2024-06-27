package matcher

import "testing"

type logger struct {
	t *testing.T
}

func NewLogger(t *testing.T) *logger {
	return &logger{t: t}
}

func (l *logger) Logf(format string, args ...interface{}) {
	if l.t == nil {
		return
	}

	l.t.Logf(format, args...)
}

func (l *logger) Log(args ...interface{}) {
	if l.t == nil {
		return
	}

	l.t.Log(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if l.t == nil {
		return
	}

	l.t.Errorf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	if l.t == nil {
		return
	}

	l.t.Error(args...)
}
