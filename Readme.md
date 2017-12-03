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
package app_test

import (
  "fmt"
  "net/http"
  "testing"

  "github.com/alexbyk/ftest/fclient"
)

type MyApp struct{}

func (app MyApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte(`{"foo": "bar"}`))
}

func Test_hello(t *testing.T) {
  app := MyApp{}
  cl := fclient.New(t, app)

  cl.Get("/hello").CodeEq(200).
    BodyContains("bar").
    JSONEq(`{"foo":"bar"}`)

  // using a response directly
  res := cl.Get("/hello")
  fmt.Println(res.Body)
}
```

# Copyright
Copyright 2017, [alexbyk.com](https://alexbyk.com)
