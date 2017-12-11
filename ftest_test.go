package ftest_test

import (
	"testing"

	"github.com/alexbyk/ftest"
	"github.com/alexbyk/ftest/internal"
)

func Test_NewLabel(t *testing.T) {
	mt := internal.NewMock(t)
	ass := ftest.NewLabel(mt, "MyLabel")
	mt.ShouldFail("MyLabel", func() { ass.Eq(22, "22") })
}

func Test_NotEq(t *testing.T) {
	ass, mt := buildAssMt(t)
	mt.ShouldFail("are equal:", func() { ass.NotEq(22, 22) })
	mt.ShouldPass(func() { ass.NotEq(22, "22") })
}

func Test_Eq(t *testing.T) {
	ass, mt := buildAssMt(t)
	mt.ShouldFail("expected:", func() { ass.Eq(22, "22") })
	mt.ShouldFail("expected:", func() {
		ass.Eq(map[string]string{"foo": "23"}, map[string]string{"foo": "22"})
	})
	mt.ShouldPass(func() { ass.Eq("22", "22") })
	mt.ShouldPass(func() {
		ass.Eq(map[string]string{"foo": "22"}, map[string]string{"foo": "22"})
	})

	// nil pointers
	var nilArr []rune
	var nilPtr *[]rune
	mt.ShouldFail("*nil*", func() { ass.Eq([]rune{}, nilArr) })
	mt.ShouldFail("*nil*", func() { ass.Eq(nilArr, []rune{}) })
	mt.ShouldPass(func() { ass.Eq(nilArr, nilArr) })
	mt.ShouldPass(func() { ass.Eq(nilPtr, nil) })
	mt.ShouldPass(func() { ass.Eq(nilArr, nil) })
	mt.ShouldPass(func() { ass.Eq(nil, nil) })
	mt.ShouldPass(func() { ass.Eq(nil, nilPtr) })
	mt.ShouldPass(func() { ass.Eq(nil, nilArr) })
}

func Test_Contains(t *testing.T) {
	ass, mt := buildAssMt(t)
	mt.ShouldPass(func() { ass.Contains("Foo", "oo") })
	mt.ShouldFail("doesn't contain", func() { ass.Contains("Foo", "bar") })
}

func Test_PanicsSubstr(t *testing.T) {
	ass, mt := buildAssMt(t)
	mt.ShouldPass(func() { ass.PanicsSubstr(func() { panic("BfooE") }, "foo") })
	mt.ShouldPass(func() { ass.PanicsSubstr(func() { panic("Any") }, "") })
	mt.ShouldFail("didn't panic",
		func() { ass.PanicsSubstr(func() {}, "") })
	mt.ShouldFail("doesn't contain",
		func() { ass.PanicsSubstr(func() { panic("NOT") }, "foo") })
}

func Test_TrueFalseNilNotNil(t *testing.T) {

	// nil pointers
	var nilArr []rune
	var nilPtr *[]rune

	ass, mt := buildAssMt(t)
	mt.ShouldFail("true", func() { ass.True(false) })
	mt.ShouldFail("false", func() { ass.False(true) })
	mt.ShouldFail("nil", func() { ass.Nil(33) })
	mt.ShouldFail("nil", func() { ass.NotNil(nil) })
	mt.ShouldFail("nil", func() { ass.NotNil(nilPtr) })
	mt.ShouldFail("nil", func() { ass.NotNil(nilArr) })
	mt.ShouldPass(func() { ass.True(true) })
	mt.ShouldPass(func() { ass.False(false) })
	mt.ShouldPass(func() { ass.Nil(nil) })
	mt.ShouldPass(func() { ass.Nil(nilArr) })
	mt.ShouldPass(func() { ass.Nil(nilPtr) })
	mt.ShouldPass(func() { ass.NotNil(0) })
	var inil interface{}
	mt.ShouldPass(func() { ass.Nil(inil) })
	mt.ShouldFail("nil", func() { ass.NotNil(inil) })
}

func buildAssMt(t *testing.T) (*ftest.Assertion, *internal.MockT) {
	mt := internal.NewMock(t)
	return ftest.New(mt), mt
}
