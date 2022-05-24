package mux

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/lestrrat-go/mux/internal/pathmatch"
)

type identMatchValues struct{}

// Values is the interface that allows users to access the
// variable path components in the given path. Use `mux.Vars`
// to access this structure
type Values interface {
	Get(string) string
}

// Vars returns the variable path components matched during
// dispatch. The returned value is always non-nil
func Vars(req *http.Request) Values {
	v := req.Context().Value(identMatchValues{})

	switch v := v.(type) {
	case Values:
		return v
	default:
		return pathmatch.Values{}
	}
}

type path struct {
	method  string
	matcher *pathmatch.Matcher
	handler http.Handler
}

// Router is the component that allows users to dispatch requests based on
// HTTP method and path, which may include variable components in the
// form of `/foo/bar/{id}` or `/foo/bar/{id:^[0-9]$}`
//
// The zero value is safe to be used, but may not be copied.
type Router struct {
	mu    sync.RWMutex
	paths []*path
}

// Handler is the generic way to associate an http.Handler to
// an HTTP method and a path.
//
// The method may be an empty string, in which case any HTTP verbs
// will match. Otherwise, the handlers will only respond to
// specific HTTP verbs that was specified
//
// The path must start with a slash (`/`).
//
// The path may contain variable components. The variable components
// are denoted by `{...}`. The pattern may either be a single string
// which signifies that anything in a path segment (a block of text
// in between slashes `/` or EOF), or a string followed by a colon (`:`)
// and a regular expression pattern.
//
// When the form `{name}` is used, any byte sequence excluding slashes
// are captured.
//
// `/foo/bar/{id}` matches `/foo/bar/123` or `/foo/bar/%31%32%33` but not `/foo/bar/123/`
// `/foo/bar/{id}/view` matches `/foo/bar/123/view` but not `/foo/bar//view`.
//
// Wehn the form `{name:regexp}` is used, the regular expression is matched
// against all remaining segments of the path, including slashes.
//
// `/foo/bar/{id:^[0-9]+}` matches `/foo/bar/123abc` but not `/foo/bar/abc123`
// `/foo/bar/{id:[0-9]+$}` matches `/foo/bar/abc123` but not `/foo/bar/123abc`
// `/foo/bar/{id:^[0-9]+$}` matches `/foo/bar/123` but not `/foo/bar/abc`
// `/foo/bar/{rest:.*$}` matches anything under `/foo/bar/`
//

func (r *Router) Handler(method string, pattern string, hh http.Handler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	m, err := pathmatch.Parse(pattern)
	if err != nil {
		return fmt.Errorf(`failed to parse path pattern: %w`, err)
	}

	r.paths = append(r.paths, &path{
		method:  method,
		matcher: m,
		handler: hh,
	})
	return nil
}

// Any declares an endpoint that responds to HTTP requests with
// any HTTP verbs in the specified path pattern
func (r *Router) Any(pattern string, hh http.Handler) error {
	return r.Handler("", pattern, hh)
}

// Get declares an endpoint that responds to HTTP GET requests
// in the specified path pattern
func (r *Router) Get(pattern string, hh http.Handler) error {
	return r.Handler(http.MethodGet, pattern, hh)
}

// Head declares an endpoint that responds to HTTP HEAD requests
// in the specified path pattern
func (r *Router) Head(pattern string, hh http.Handler) error {
	return r.Handler(http.MethodHead, pattern, hh)
}

// Post declares an endpoint that responds to HTTP POST requests
// in the specified path pattern
func (r *Router) Post(pattern string, hh http.Handler) error {
	return r.Handler(http.MethodPost, pattern, hh)
}

// Put declares an endpoint that responds to HTTP PUT requests
// in the specified path pattern
func (r *Router) Put(pattern string, hh http.Handler) error {
	return r.Handler(http.MethodPut, pattern, hh)
}

// Patch declares an endpoint that responds to HTTP PATCH requests
// in the specified path pattern
func (r *Router) Patch(pattern string, hh http.Handler) error {
	return r.Handler(http.MethodPatch, pattern, hh)
}

// Delete declares an endpoint that responds to HTTP DELETE requests
// in the specified path pattern
func (r *Router) Delete(pattern string, hh http.Handler) error {
	return r.Handler(http.MethodDelete, pattern, hh)
}

// Connect declares an endpoint that responds to HTTP CONNECT requests
// in the specified path pattern
func (r *Router) Connect(pattern string, hh http.Handler) error {
	return r.Handler(http.MethodConnect, pattern, hh)
}

// Options declares an endpoint that responds to HTTP OPTIONS requests
// in the specified path pattern
func (r *Router) Options(pattern string, hh http.Handler) error {
	return r.Handler(http.MethodOptions, pattern, hh)
}

// Trace declares an endpoint that responds to HTTP TRACE requests
// in the specified path pattern
func (r *Router) Trace(pattern string, hh http.Handler) error {
	return r.Handler(http.MethodTrace, pattern, hh)
}

// ServeHTTP implements the http.Handler interface, allowing `*Router`
// to be passed to anythign that expects an http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, path := range r.paths {
		if method := path.method; method != "" {
			if req.Method != method {
				continue
			}
		}

		mv, err := path.matcher.Match(req.URL.Path)
		if err != nil {
			continue
		}

		ctx := context.WithValue(req.Context(), identMatchValues{}, mv)
		path.handler.ServeHTTP(w, req.WithContext(ctx))
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
