package fclient_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/alexbyk/ftest/fclient"
)

func makeHandlerCookie(name string, cookie *http.Cookie) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/get":
			c, _ := r.Cookie(name)
			if c == nil {
				w.Write([]byte("empty"))
				return
			}
			w.Write([]byte(cookie.Value))
			return
		case "/set":
			http.SetCookie(w, cookie)
			w.Write([]byte("Stored"))
		default:
			panic("unexpected: " + r.URL.Path)
		}
	}

}

func Test_CookiesGetSet(t *testing.T) {
	cookie := &http.Cookie{Name: "foo", Value: "bar"}
	cl := fclient.New(t, makeHandlerCookie("foo", cookie))
	cl.Do(cl.NewRequest("GET", "/get", nil)).BodyEq("empty")
	cl.Do(cl.NewRequest("GET", "/set", nil))
	cl.Do(cl.NewRequest("GET", "/get", nil)).BodyEq("bar")

	cookie.Expires = time.Unix(0, 0)
	cl.Do(cl.NewRequest("GET", "/set", nil))
	cl.Do(cl.NewRequest("GET", "/get", nil)).BodyEq("empty")
}

func Test_CookiesAnyDomains(t *testing.T) {
	cookie := &http.Cookie{Name: "foo", Value: "bar", Domain: "alexbyk.com"}
	cl := fclient.New(t, makeHandlerCookie("foo", cookie))
	cl.Do(cl.NewRequest("GET", "/set", nil))
	cl.Do(cl.NewRequest("GET", "/get", nil)).BodyEq("bar")
}

func Test_CookiesPath(t *testing.T) {
	cookie := &http.Cookie{Name: "foo", Value: "bar", Path: "/another"}
	cl := fclient.New(t, makeHandlerCookie("foo", cookie))
	cl.Do(cl.NewRequest("GET", "/set", nil))
	cl.Do(cl.NewRequest("GET", "/get", nil)).BodyEq("empty")

	cookie.Path = "/get"
	cl.Do(cl.NewRequest("GET", "/set", nil))
	cl.Do(cl.NewRequest("GET", "/get", nil)).BodyEq("bar")
}

func Test_EmptyJar(t *testing.T) {
	cookie := &http.Cookie{Name: "foo", Value: "bar"}
	cl := fclient.New(t, makeHandlerCookie("foo", cookie))
	cl.Jar = nil
	cl.Do(cl.NewRequest("GET", "/set", nil))
	cl.Do(cl.NewRequest("GET", "/get", nil)).BodyEq("empty")
}
