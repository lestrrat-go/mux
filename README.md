mux
===

Package `github.com/lestrrat-go/mux` implements a very simple HTTP mux.

```go
var r mux.Router

r.Handler(http.MethodGet, `/foo/bar/baz/{id}`, http.Handler(...))
r.Get(`/foo/bar/baz/{id}`, http.Handler(...))

http.ListenAndServe(":8080", r)
```

# FAQ

## Who is this for?

People who want to write a basic HTTP server, no bells in whistles, but at the
same time don't want to have to manually write all the dispatching rules
by comparing `req.Method` or `req.URL.Path` and doing regular expression matches, etc.

If you need a more full-fledged support from a framework, this is not for you.
We also don't expect you to write a real production server that can handle
complex routing. This is mostly useful for PoCs and sample apps.

## Why don't you use a framework X? Why reinvent the wheel?

We just wanted something to dispatch requests to `http.Handler`s, and nothing
major like a framework -- we did not feel like importing a big framework just
for this task.

There were some lightweight contenders, but one of our favorite tools for this
sort of task was in a transitional state where they were looking for maintainers
and nobody was actively maintaining it. Also, even if it was properly maintained,
it did do a bit more than what we needed.
