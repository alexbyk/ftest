# Changelog

## v0.2.2 (2017-12-09)
- Improved nil checks
- Added fclient.NewResponse builder function

## v0.2.0 (2017-12-03)
- [Breaking] `fclient.New` now accepts `http.Handler` instead of `a function`. You can use `http.HandlerFunc` as an adapter
- make `Client.Handler` public
