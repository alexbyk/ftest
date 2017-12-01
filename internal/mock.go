package internal

import (
	"fmt"
	"strings"
	"testing"
)

// This is for internal purposes

// MockT implements internal test interfaces
type MockT struct {
	t   *testing.T
	err string
}

// NewMock creates new mock
func NewMock(t *testing.T) *MockT { return &MockT{t: t} }

// Helper Mock
func (mt *MockT) Helper() {}

// Fatalf mock
func (mt *MockT) Fatalf(format string, args ...interface{}) {
	mt.err = fmt.Sprintf(format, args...)
	panic(mt.err)
}

// ShouldFail checks if a passed function fails with a substring
func (mt *MockT) ShouldFail(substr string, fn func()) {
	mt.t.Helper()
	mt.err = ""
	eval(fn)
	if mt.err == "" {
		mt.t.Fatalf("Should fail with '%v', but didn't", substr)
	}
	if !strings.Contains(mt.err, substr) {
		mt.t.Fatalf("Should fail with '%v', but failed with\n%v", substr, mt.err)
	}
}

// ShouldPass check if a passed function dosn't fail
func (mt *MockT) ShouldPass(fn func()) {
	mt.t.Helper()
	mt.err = ""
	eval(fn)
	if mt.err != "" {
		mt.t.Fatalf("Shoud pass, but failed with '%v'", mt.err)
	}
}

func eval(fn func()) {
	defer func() { recover() }()
	fn()
}
