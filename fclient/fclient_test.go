package fclient_test

import (
	"net/http"
	"testing"

	"github.com/alexbyk/ftest"
	"github.com/alexbyk/ftest/fclient"
	"github.com/alexbyk/ftest/internal"
)

func Test_Request(t *testing.T) {
	cl := fclient.New(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Method + " " + r.URL.String()))
	})
	ass := ftest.New(t)
	ass.Eq(cl.Get("http://foo.bar/baz").Body.String(), "GET http://foo.bar/baz")

	ass.Eq(cl.Post("http://foo.bar/", nil).Body.String(), "POST http://foo.bar/")
	ass.Eq(cl.Do(cl.NewRequest("HEAD", "http://foo.bar/baz", nil)).Body.String(),
		"HEAD http://foo.bar/baz")
}

func Test_CodeEq(t *testing.T) {
	cl, mt := buildClientMt(t, makeBodyResp(201, ""))
	mt.ShouldFail("200", func() { cl.Get("/").CodeEq(201) })
	mt.ShouldPass(func() { cl.Get("/").CodeEq(200) })
}

func Test_BodyEq(t *testing.T) {
	cl, mt := buildClientMt(t, makeBodyResp(201, "OK"))
	mt.ShouldFail("BodyEq", func() { cl.Get("/").BodyEq("Foo") })
	mt.ShouldFail("Unexpected type int", func() { cl.Get("/").BodyEq(22) })
	mt.ShouldPass(func() { cl.Get("/").BodyEq("OK") })
	mt.ShouldPass(func() { cl.Get("/").BodyEq([]byte("OK")) })
}

func Test_BodyContains(t *testing.T) {
	cl, mt := buildClientMt(t, makeBodyResp(201, "Foo"))
	mt.ShouldFail("BodyContains", func() { cl.Get("/").BodyContains("ar") })
	mt.ShouldPass(func() { cl.Get("/").BodyContains("oo") })
}

func Test_JSONEq(t *testing.T) {

	type prodT struct {
		Name  string `json:"name,omitempty"`
		Price int    `json:"price,omitempty"`
	}

	invalidMsg := "isn't a valid JSON"
	notEqMsg := "JSONEq"
	testCases := []struct {
		resp string
		in   interface{}
		err  string
	}{
		// bad resp
		{`{"foo":  "ok"`, `{"foo": "ok"}`, invalidMsg},

		// string
		{`{"foo":  "ok"}`, `{"foo": "ok"`, invalidMsg},
		{`{"foo":  "ok"}`, `{"foo": "NOT"}`, notEqMsg},
		{`{"foo":  "ok", "bar": "ok"}`, `{"foo": "ok"}`, notEqMsg},

		{`{"foo":  "ok", "bar": "ok"}`, `{"bar":"ok","foo": "ok"}`, ""},
		{`{"foo":  "ok", "bar": "ok"}`, `{"foo":  "ok", "bar": "ok"}`, ""},

		//// map
		{`{"foo":  "ok"`, map[string]string{"foo": "NOT"}, invalidMsg},
		{`{"foo":  "ok"}`, map[string]string{"foo": "NOT"}, notEqMsg},

		{`{"foo":  "ok", "bar": "ok"}`, map[string]string{"bar": "ok", "foo": "ok"}, ""},
		{`{"foo":  "ok", "bar": "ok"}`, map[string]string{"foo": "ok", "bar": "ok"}, ""},

		// obj
		{`{"Name":  "ok", "price": 22}`, prodT{"ok", 22}, notEqMsg},
		{`{"name":  "ok", "price": 22}`, prodT{"NOT", 22}, notEqMsg},
		{`{"name":  "ok", "price": "22"}`, prodT{"NOT", 22}, notEqMsg},

		{`{"name":  "ok", "price": 22}`, prodT{"ok", 22}, ""},
		{`{"price": 22, "name":  "ok"}`, prodT{"ok", 22}, ""},

		{`null`, nil, ""},

		// in json `22` is number, `"22"` is string
		{`22`, 22, ""},
		{`22`, "22", ""},
		{`"22"`, `"22"`, ""},
		{`"22"`, 22, notEqMsg},
		{`22`, `"22"`, notEqMsg},
	}
	for _, tc := range testCases {
		t.Run(tc.resp, func(t *testing.T) {

			cl, mt := buildClientMt(t, makeBodyResp(200, tc.resp))
			if tc.err == "" {
				mt.ShouldPass(func() { cl.Get("/").JSONEq(tc.in) })
			} else {
				mt.ShouldFail(tc.err, func() { cl.Get("/").JSONEq(tc.in) })
			}

		})
	}

}

func Test_DefaultHeaders(t *testing.T) {
	cl := fclient.New(t, func(w http.ResponseWriter, r *http.Request) {})
	ft := ftest.New(t)
	cl.DefaultHeaders["foo"] = "FOO"
	req := cl.NewRequest("GET", "/", nil)
	ft.Eq(req.Header.Get("foo"), "FOO")
}

func Test_HeaderEq(t *testing.T) {
	cl, mt := buildClientMt(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("foo", "FOO")
	})
	mt.ShouldFail("HeaderEq", func() { cl.Get("/").HeaderEq("foo", "bad") })
	mt.ShouldPass(func() { cl.Get("/").HeaderEq("foo", "FOO") })
}

func buildClientMt(t *testing.T, handler http.HandlerFunc) (*fclient.Client, *internal.MockT) {
	tt := internal.NewMock(t)
	return fclient.New(tt, handler), tt
}

func makeBodyResp(code int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	}
}
