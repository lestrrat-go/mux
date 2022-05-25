package pathmatch

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type TokReturn struct {
	Tok int
	Lit interface{}
	Pos position
	Err error
}

func TestTokenizer(t *testing.T) {
	testcases := []struct {
		Input    string
		Expected []TokReturn
	}{
		{
			Input: "/foo/bar/baz/{id}/view",
			Expected: []TokReturn{
				{
					Tok: tLiteral,
					Lit: "/foo/bar/baz/",
					Pos: position{
						line: 1,
						col:  1,
					},
				},
				{
					Tok: tOpenBrace,
					Lit: "{",
					Pos: position{
						line: 1,
						col:  14,
					},
				},
				{
					Tok: tLiteral,
					Lit: "id",
					Pos: position{
						line: 1,
						col:  15,
					},
				},
				{
					Tok: tCloseBrace,
					Lit: "}",
					Pos: position{
						line: 1,
						col:  17,
					},
				},
				{
					Tok: tLiteral,
					Lit: "/view",
					Pos: position{
						line: 1,
						col:  18,
					},
				},
				{
					Tok: tEOF,
					Pos: position{
						line: 1,
						col:  23,
					},
					Err: io.EOF,
				},
			},
		},
		{
			Input: "/foo/bar/baz/{id:^[0-9]+}/view",
			Expected: []TokReturn{
				{
					Tok: tLiteral,
					Lit: "/foo/bar/baz/",
					Pos: position{
						line: 1,
						col:  1,
					},
				},
				{
					Tok: tOpenBrace,
					Lit: "{",
					Pos: position{
						line: 1,
						col:  14,
					},
				},
				{
					Tok: tLiteral,
					Lit: "id",
					Pos: position{
						line: 1,
						col:  15,
					},
				},
				{
					Tok: tColon,
					Lit: ":",
					Pos: position{
						line: 1,
						col:  17,
					},
				},
				{
					Tok: tLiteral,
					Lit: "^[0-9]+",
					Pos: position{
						line: 1,
						col:  18,
					},
				},
				{
					Tok: tCloseBrace,
					Lit: "}",
					Pos: position{
						line: 1,
						col:  25,
					},
				},
				{
					Tok: tLiteral,
					Lit: "/view",
					Pos: position{
						line: 1,
						col:  26,
					},
				},
				{
					Tok: tEOF,
					Pos: position{
						line: 1,
						col:  31,
					},
					Err: io.EOF,
				},
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Input, func(t *testing.T) {
			tokenizer := newTokenizer(tc.Input)
			for i, exp := range tc.Expected {
				t.Run(fmt.Sprintf(`token %d`, i), func(t *testing.T) {
					tok, lit, pos, err := tokenizer.Token()
					require.Equal(t, exp.Tok, tok, `tokenizer.Token() token should match`)
					require.Equal(t, exp.Lit, lit, `tokenizer.Token() lit should match`)
					require.Equal(t, exp.Pos, pos, `tokenizer.Token() pos should match`)
					require.Equal(t, exp.Err, err, `tokenizer.Token() error should match`)
				})
			}
		})
	}
}
