/*Package fclient provides a testing client for http handlers.

It supports simple cookie storage, so can be used to test sites with sign-in functionality

  func hello(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"foo": "bar"}`))
  }

  func Test_hello(t *testing.T) {
    cl := fclient.New(t, hello)

    res := cl.Get("/hello").CodeEq(200).
      BodyContains("bar").
      JSONEq(`{"foo":"bar"}`)

    fmt.Println(res.Body)
  }

To test your final app with multiple routes/framework, use the following example(gin)

  func TestGin(t *testing.T) {

    // app
    app := gin.New()
    app.GET("/foo", func(c *gin.Context) {
      c.String(200, "Foo")
    })

    // test
    cl := fclient.New(t, app.ServeHTTP)
    cl.Get("/foo").BodyEq("Foo")
  }

*/
package fclient

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"

	"github.com/alexbyk/ftest"
)

type test interface {
	Fatalf(format string, args ...interface{})
	Helper()
}

// ----------- Client -----------

// Client is a struct which holds current test, LastResponse and LastRequest instances
type Client struct {
	t       test
	handler http.HandlerFunc

	// DefaultHeaders is a map for default headers
	DefaultHeaders map[string]string

	// Jar holds cookies. Can be set to nil to turn cookies off
	Jar http.CookieJar
}

// New builds a new http testing client
func New(t test, handler http.HandlerFunc) *Client {
	t.Helper()
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{t: t, Jar: jar, handler: handler,
		DefaultHeaders: map[string]string{},
	}
}

// Get makes a GET Request with nil body
func (cl *Client) Get(path string) *Response {
	cl.t.Helper()
	return cl.Do(cl.NewRequest("GET", path, nil))
}

// Post makes a POST Request
func (cl *Client) Post(path string, body interface{}) *Response {
	cl.t.Helper()
	reader := bytes.NewReader(toBytes(cl.t, body))
	req := cl.NewRequest("POST", path, reader)
	return cl.Do(req)
}

func urlFromReq(req *http.Request) *url.URL {
	u, err := url.Parse("https://example.com")
	if err != nil {
		panic(err)
	}
	u.Path = req.URL.Path
	return u
}

// Do invokes the handler, passing a given argument as a request,
// and stores cookies from the response in Jar. A ".Domain" field of every cookie
// will be cleared
func (cl *Client) Do(req *http.Request) *Response {
	cl.t.Helper()
	resp := &Response{ResponseRecorder: httptest.NewRecorder(), t: cl.t}
	cl.handler(resp, req)
	jar := cl.Jar
	if jar == nil {
		return resp
	}
	cookies := []*http.Cookie{}
	for _, c := range resp.Result().Cookies() {
		c.Domain = ""
		cookies = append(cookies, c)
	}
	jar.SetCookies(urlFromReq(req), cookies)
	return resp
}

// NewRequest creates a new request (r) by httptest.NewRequest,
// then fills r.Headers from cl.DefaultHeaders,
// then fills cookies from cl.Jar
// and returns r
func (cl *Client) NewRequest(method, path string, body io.Reader) *http.Request {
	cl.t.Helper()
	req := httptest.NewRequest(method, path, body)
	if cl.Jar == nil {
		return req
	}
	for k, v := range cl.DefaultHeaders {
		req.Header.Set(k, v)
	}
	for _, c := range cl.Jar.Cookies(urlFromReq(req)) {
		req.AddCookie(c)
	}
	return req
}

// ----------- Response -----------

// A Response represents the response from an HTTP request. It also inherits
// all methods from *httptest.ResponseRecorder
type Response struct {
	*httptest.ResponseRecorder
	t test
}

// CodeEq checks if a response code is equal to
func (resp *Response) CodeEq(expected int) *Response {
	resp.t.Helper()
	got := resp.Code
	ftest.NewLabel(resp.t, "CodeEq").Eq(got, expected)
	return resp
}

// BodyEq checks if a response body is equal to the given argument.
// Accepts []byte or string
func (resp *Response) BodyEq(expected interface{}) *Response {
	resp.t.Helper()
	expStr := string(toBytes(resp.t, expected))
	got := resp.Body.String()
	ftest.NewLabel(resp.t, "BodyEq").Eq(got, expStr)
	return resp
}

// BodyContains checks if a response body contains a given string
func (resp *Response) BodyContains(substr string) *Response {
	resp.t.Helper()
	ftest.NewLabel(resp.t, "BodyContains").
		Contains(resp.Body.String(), substr)
	return resp
}

// JSONEq checks if given argument is equal to response JSON.
// Argument can be a string, in which case it will be used as a raw JSON, or other kind,
// in which case it will be transformed to JSON
func (resp *Response) JSONEq(expected interface{}) *Response {
	resp.t.Helper()

	// Check if response is a valid json
	var respObj interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &respObj)
	if err != nil {
		resp.t.Fatalf("JSONEq: response body isn't a valid JSON:\n%s", resp.Body.String())
	}

	// Because an order is undefined, we convert all to bytes than to interface{}
	var expectedBytes []byte
	if v, ok := expected.(string); ok {
		expectedBytes = []byte(v)
	} else {
		expectedBytes, err = json.Marshal(expected)
		if err != nil {
			resp.t.Fatalf("Can't convert to JSON: %v", err)
		}
	}

	var expectedObj interface{}
	err = json.Unmarshal(expectedBytes, &expectedObj)
	if err != nil {
		resp.t.Fatalf("JSONEq: argument isn't a valid JSON %v", err)
	}

	ftest.NewLabel(resp.t, "JSONEq").Eqf(respObj, expectedObj,
		"got:\n%s\nexpected:\n%s", resp.Body.String(), expectedBytes)
	return resp
}

// HeaderEq checks if the first http header with given name is equal to the given value
func (resp *Response) HeaderEq(key, value string) *Response {
	resp.t.Helper()
	ftest.NewLabel(resp.t, "HeaderEq").
		Eq(resp.Result().Header.Get(key), value)
	return resp
}

func toBytes(t test, in interface{}) []byte {
	t.Helper()
	var data []byte
	switch in := in.(type) {
	case []byte:
		data = in
	case string:
		data = []byte(in)
	case nil:
		return []byte(nil)
	default:
		t.Fatalf("Unexpected type %T!", in)
	}
	return data
}
