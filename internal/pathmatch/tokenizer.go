package pathmatch

import (
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	tEOF = iota
)

type position struct {
	line int
	col  int
}

type tokenizer struct {
	src          *strings.Reader
	expectRegexp bool
	eof          bool
	offset       int
	lineHead     int
	line         int
	buf          rune // no backtracking, so 1 rune is enough
}

func newTokenizer(s string) *tokenizer {
	return &tokenizer{
		src: strings.NewReader(s),
		buf: utf8.RuneError,
	}
}

func (t *tokenizer) peek() rune {
	if t.eof {
		return utf8.RuneError
	}

	if t.buf == utf8.RuneError {
		r, _, err := t.src.ReadRune()
		if err != nil { // can only happen if at EOF
			t.eof = true
			return utf8.RuneError
		}
		t.buf = r
	}

	return t.buf
}

func (t *tokenizer) next() {
	if t.eof {
		return
	}
	t.offset++
	t.buf = utf8.RuneError
}

func (t *tokenizer) position() position {
	return position{
		line: t.line + 1,
		col:  t.offset - t.lineHead + 1,
	}
}

func (t *tokenizer) Token() (int, interface{}, position, error) {
	t.skipWhitespace()

	pos := t.position()
	var tok int
	var lit interface{}
	r := t.peek()
	switch r {
	case utf8.RuneError:
		return tEOF, nil, pos, io.EOF
	case '{':
		tok = tOpenBrace
		lit = "{"
		t.next()
	case '}':
		tok = tCloseBrace
		lit = "}"
		t.next()
	case ':':
		tok = tColon
		lit = ":"
		t.next()
		t.expectRegexp = true
	default:
		tok = tLiteral
		lit = t.literal()
	}

	return tok, lit, pos, nil
}

func (t *tokenizer) skipWhitespace() {
	for {
		r := t.peek()
		if r == utf8.RuneError || !unicode.IsSpace(r) {
			return
		}
		t.next()
	}
}
func (t *tokenizer) literal() string {
	var b strings.Builder
LOOP:
	for {
		r := t.peek()
		if t.expectRegexp {
			switch r {
			case '}', utf8.RuneError:
				break LOOP
			default:
				b.WriteRune(r)
				t.next()
			}
		} else {
			switch r {
			case ':', '{', '}', utf8.RuneError:
				break LOOP
			default:
				b.WriteRune(r)
				t.next()
			}
		}
	}
	t.expectRegexp = false
	return b.String()
}
