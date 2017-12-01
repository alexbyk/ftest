[![Build Status](https://travis-ci.org/alexbyk/ftest.svg?branch=master)](https://travis-ci.org/alexbyk/ftest)
# About

- `ftest` is a simple and easy to use go testing library with fluent design and exact failure messages.
- `fclient` is a simple http testing client, based on `ftest`.

# Documentation:
- [ftest](https://godoc.org/github.com/alexbyk/ftest)
- [fclient](https://godoc.org/github.com/alexbyk/ftest/fclient)

# Installation
```
go get -u github.com/alexbyk/ftest
```

or
```
dep ensure --add github.com/alexbyk/ftest
```

# Usage
## ftest
```go
import (
  "testing"

  "github.com/alexbyk/ftest"
)

func TestFoo(t *testing.T) {
  ftest.New(t).Eq(2, 2).
    Contains("FooBarBaz", "Bar").
    PanicsSubstr(func() { panic("Foo") }, "Foo")
}

```

You can also use your own label to group your test:
```go
ft := ftest.NewLabel(t, "MyLabel")
```

## fclient
```go

import (
  "fmt"
  "net/http"
  "testing"

  "github.com/alexbyk/ftest/fclient"
)

func hello(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte(`{"foo": "bar"}`))
}

func Test_hello(t *testing.T) {
  cl := fclient.New(t, hello)

  cl.Get("/hello").CodeEq(200).
    BodyContains("bar").
    JSONEq(`{"foo":"bar"}`)

  // using response directly
  res := cl.Get("/hello")
  fmt.Println(res.Body)
}

```

### How to use with multiple routes/framework app
You should provide a main `http.HandlerFunc` function, for example with `gin` it will be something like this:
```go
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
```

# Copyright
Copyright 2017, [alexbyk.com](https://alexbyk.com)
