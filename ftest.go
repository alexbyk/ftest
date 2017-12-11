/*Package ftest is a simple easy to use testing library.
It prints only the exact location of failed test (without scary stack trace).
It uses a fluent desing. It stops a test on the first failure

	ftest.New(t).Eq(2, 2).
		Contains("FooBarBaz", "Bar").
    PanicsSubstr(func() { panic("Foo") }, "Foo")

Also to make testing simpler, this package performs extra nil checks, so nil, non initialized slice and empty pointer
- all will be considered equal (and that's what you are expecting)
*/
package ftest

import (
	"fmt"
	"reflect"
	"strings"
)

// Test is an interface with FatalF and Helper methods, which are required by the Client
type test interface {
	Fatalf(format string, args ...interface{})
	Helper()
}

// ----------- Assertion -----------

// Assertion represents an assertion which holds current a *testing.T object
type Assertion struct {
	t     test
	label string
}

// NewLabel creates an Assertion instance with a label
func NewLabel(t test, label string) *Assertion {
	return &Assertion{t: t, label: label}
}

// New creates an Assertion instance with label "Assertion"
func New(t test) *Assertion {
	return NewLabel(t, "Assertion")
}

// TODO: avoid defer/recover by checking a kind
func isNil(v interface{}) (ret bool) {
	defer func() { recover() }()
	ret = reflect.ValueOf(v).IsNil()
	return
}

// NotEq tests if 2 arguments are equal
func (ass *Assertion) NotEq(got, expected interface{}) *Assertion {
	ass.t.Helper()
	return ass.NotEqf(got, expected, "are equal: %v(%v)", reflect.TypeOf(got), got)
}

func (ass *Assertion) deepEq(got, expected interface{}) bool {
	ass.t.Helper()
	gotV := reflect.ValueOf(got)
	expectedV := reflect.ValueOf(expected)
	switch {
	case !gotV.IsValid():
		return !expectedV.IsValid() || isNil(expected)
	case !expectedV.IsValid():
		return !gotV.IsValid() || isNil(got)
	}

	return reflect.DeepEqual(got, expected) || isNil(got) && isNil(expected)
}

// NotEqf is an f version of NotEq
func (ass *Assertion) NotEqf(got, expected interface{}, format string, args ...interface{}) *Assertion {
	ass.t.Helper()
	if ass.deepEq(got, expected) {
		ass.fail(format, args...)
	}
	return ass
}

// Eq tests if 2 arguments are equal
func (ass *Assertion) Eq(got, expected interface{}) *Assertion {
	ass.t.Helper()
	var gotNilS, expNilS string
	if isNil(got) {
		gotNilS = "*nil*"
	}
	if isNil(expected) {
		expNilS = "*nil*"
	}

	return ass.Eqf(got, expected, "got: %v(%s%v), expected: %v(%s%v)",
		reflect.TypeOf(got), gotNilS, got, reflect.TypeOf(expected), expNilS, expected)
}

// Eqf is an f version of Eq
func (ass *Assertion) Eqf(got, expected interface{}, format string, args ...interface{}) *Assertion {
	ass.t.Helper()
	if !ass.deepEq(got, expected) {
		ass.fail(format, args...)
	}
	return ass
}

// Contains checks if the first argument contains a second one
func (ass *Assertion) Contains(str, substr string) *Assertion {
	ass.t.Helper()
	return ass.Containsf(str, substr, `"%s" doesn't contain "%s"`,
		str, substr)

}

// Containsf is like Contains but fails with given format and arguments
func (ass *Assertion) Containsf(str, substr, format string, args ...interface{}) *Assertion {
	ass.t.Helper()
	if !strings.Contains(str, substr) {
		ass.fail(format, args...)
	}
	return ass
}

// PanicsSubstr tests if a given function causes and panic
// and an error contains given substring.
// You can pass an empty string("") if the error doesn't matter
func (ass *Assertion) PanicsSubstr(fn func(), substr string) (ret *Assertion) {
	ass.t.Helper()
	ret = ass

	defer func() {
		ass.t.Helper()
		e := recover()
		if e == nil {
			ass.fail("Function %#v didn't panic as expected", fn)
		}
		errStr := fmt.Sprintf("%s", e)
		if !strings.Contains(errStr, substr) {
			ass.fail("Error \"%s\" doesn't contain substring \"%s\"", errStr, substr)
		}
	}()

	fn()
	return
}

// NotNil tests if a given argument isn't nil
func (ass *Assertion) NotNil(got interface{}) *Assertion {
	ass.t.Helper()
	return ass.NotNilf(got, "%v(%v) is nil", reflect.TypeOf(got), got)
}

// NotNilf is an f vetsion of NotNil
func (ass *Assertion) NotNilf(got interface{}, format string, args ...interface{}) *Assertion {
	ass.t.Helper()
	return ass.NotEqf(got, nil, format, args...)
}

// Nil tests if a given argument is nil
func (ass *Assertion) Nil(got interface{}) *Assertion {
	ass.t.Helper()
	return ass.Nilf(got, "%v(%v) isn't nil", reflect.TypeOf(got), got)
}

// Nilf is an f vetsion of Nil
func (ass *Assertion) Nilf(got interface{}, format string, args ...interface{}) *Assertion {
	ass.t.Helper()
	return ass.Eqf(got, nil, format, args...)
}

// False tests if a given argument is false
func (ass *Assertion) False(got bool) *Assertion {
	ass.t.Helper()
	return ass.Falsef(got, "Not false")
}

// Falsef is an f version of False
func (ass *Assertion) Falsef(got bool, format string, args ...interface{}) *Assertion {
	ass.t.Helper()
	return ass.Eqf(got, false, format, args...)
}

// True tests if a given argument is true
func (ass *Assertion) True(got bool) *Assertion {
	ass.t.Helper()
	return ass.Truef(got, "Not true")
}

// Truef is an f version of True
func (ass *Assertion) Truef(got bool, format string, args ...interface{}) *Assertion {
	ass.t.Helper()
	return ass.Eqf(got, true, format, args...)
}

func (ass *Assertion) fail(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	ass.t.Helper()
	ass.t.Fatalf("[%s] %s", ass.label, msg)
}
