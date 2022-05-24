package mux_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lestrrat-go/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMux(t *testing.T) {
	type MuxRequestResponseStatusPair struct {
		Method string
		Path   string
		Header http.Header
		Body   io.Reader
		Status int
	}
	testcases := []struct {
		Method   string
		Pattern  string
		Handler  http.Handler
		Requests []MuxRequestResponseStatusPair
	}{
		{
			Method:  http.MethodGet,
			Pattern: `/foo/bar/baz/{id:^[0-9a-z]+$}`,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				name := r.Header.Get(`x-var-name`)
				value := r.Header.Get(`x-var-value`)
				if name != "" {
					if vars.Get(name) == value {
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, `expected var %q to be %q, got %q`, name, value, vars.Get(name))
					}
				}
			}),
			Requests: []MuxRequestResponseStatusPair{
				{
					Method: http.MethodGet,
					Header: http.Header{
						`x-var-name`:  []string{`id`},
						`x-var-value`: []string{`abcdef`},
					},
					Path: `/foo/bar/baz/abcdef`,
				},
				{
					Method: http.MethodHead,
					Path:   `/foo/bar/baz/abcdef`,
					Status: http.StatusNotFound,
				},
				{
					Method: http.MethodGet,
					Path:   `/foo/bar/baz`,
					Status: http.StatusNotFound,
				},
				{
					Method: http.MethodGet,
					Path:   `/foo/bar/baz/012345`,
				},
				{
					Method: http.MethodGet,
					Path:   `/foo/bar/baz/012abc`,
				},
			},
		},
		{
			Method:  http.MethodGet,
			Pattern: `/foo/bar/baz/{id:[0-9a-z]+}/view`,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				name := r.Header.Get(`x-var-name`)
				value := r.Header.Get(`x-var-value`)
				if name != "" {
					if vars.Get(name) == value {
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, `expected var %q to be %q, got %q`, name, value, vars.Get(name))
					}
				}
			}),
			Requests: []MuxRequestResponseStatusPair{
				{
					Method: http.MethodGet,
					Header: http.Header{
						`x-var-name`:  []string{`id`},
						`x-var-value`: []string{`abcdef`},
					},
					Path: `/foo/bar/baz/abcdef/view`,
				},
				{
					Method: http.MethodHead,
					Path:   `/foo/bar/baz/abcdef/view`,
					Status: http.StatusNotFound,
				},
				{
					Method: http.MethodGet,
					Path:   `/foo/bar/baz`,
					Status: http.StatusNotFound,
				},
				{
					Method: http.MethodGet,
					Path:   `/foo/bar/baz/012345/view`,
				},
				{
					Method: http.MethodGet,
					Path:   `/foo/bar/baz/012abc/view`,
				},
			},
		},
		{
			Method:  http.MethodGet,
			Pattern: `/foo/bar/baz/{id:.*$}`,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				name := r.Header.Get(`x-var-name`)
				value := r.Header.Get(`x-var-value`)
				if name != "" {
					if vars.Get(name) == value {
						w.WriteHeader(http.StatusOK)
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, `expected var %q to be %q, got %q`, name, value, vars.Get(name))
					}
				}
			}),
			Requests: []MuxRequestResponseStatusPair{
				{
					Method: http.MethodGet,
					Header: http.Header{
						`x-var-name`:  []string{`id`},
						`x-var-value`: []string{`abc123/hello/world`},
					},
					Path: `/foo/bar/baz/abc123/hello/world`,
				},
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Pattern, func(t *testing.T) {
			var r mux.Router
			require.NoError(t, r.Handler(tc.Method, tc.Pattern, tc.Handler), `r.Handler should succeed`)
			srv := httptest.NewServer(&r)
			defer srv.Close()

			for _, r := range tc.Requests {
				r := r
				t.Run(fmt.Sprintf("%s %s", r.Method, r.Path), func(t *testing.T) {
					req, err := http.NewRequest(r.Method, srv.URL+r.Path, r.Body)
					require.NoError(t, err, `http.NewRequest should succeed`)
					if hdr := r.Header; hdr != nil {
						req.Header = hdr
					}

					res, err := http.DefaultClient.Do(req)
					require.NoError(t, err, `http.DefaultClient.Do should succeed`)

					if r.Status == 0 {
						r.Status = http.StatusOK
					}
					if !assert.Equal(t, r.Status, res.StatusCode, `status code should match`) {
						buf, err := io.ReadAll(res.Body)
						require.NoError(t, err, `io.ReadAll should succeed`)
						t.Logf("%s", buf)
					}
				})
			}
		})
	}
}
