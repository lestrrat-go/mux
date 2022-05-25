package pathmatch_test

import (
	"testing"

	"github.com/lestrrat-go/mux/internal/pathmatch"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testcases := []struct {
		Pattern string
		Input   string
	}{
		{
			Pattern: "/foo/bar/baz/{id}/view",
			Input:   "/foo/bar/baz/abc123/view",
		},
		{
			Pattern: "/foo/bar/baz/{id:^[0-9]+}/view",
			Input:   "/foo/bar/baz/0123456/view",
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Pattern, func(t *testing.T) {
			p, err := pathmatch.Parse(tc.Pattern)
			require.NoError(t, err, `path.Parse should succeed`)

			mv, err := p.Match(tc.Input)
			require.NoError(t, err, `p.Match should succeed`)
			t.Logf("%#v", mv)
		})
	}
}
