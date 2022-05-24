mux
===

Package `github.com/lestrrat-go/mux` implements a very simple HTTP mux.

```go
var r mux.Router

r.Handler(http.MethodGet, `/foo/bar/baz/{id}`, http.Handler(...))
r.Get(`/foo/bar/baz/{id}`, http.Handler(...))

http.ListenAndServe(":8080", r)
```
