package pathmatch

import (
	"context"
	"fmt"
)

type lexer struct {
	expr      Expression
	exprs     []Expression
	tokenizer *tokenizer
	err       chan error
	recentLit interface{}
	recentPos position

	done <-chan struct{}
}

func newLexer(ctx context.Context, s string) *lexer {
	return &lexer{
		tokenizer: newTokenizer(s),
		err:       make(chan error, 1),
		done:      ctx.Done(),
	}
}

type yyToken struct {
	tok int
	lit interface{}
	pos position
}

func (l *lexer) Lex(lval *yySymType) int {
	tok, lit, pos, err := l.tokenizer.Token()
	if err != nil {
		go l.emitError(err)
		return -1
	}
	if tok == TokEOF {
		return 0
	}

	lval.token = &yyToken{
		tok: tok,
		lit: lit,
		pos: pos,
	}
	l.recentLit = lit
	l.recentPos = pos

	return tok
}

func (l *lexer) makeError(e interface{}) error {
	switch e := e.(type) {
	case error:
		return fmt.Errorf(`parse error: line %d, column %d: %q: %w`, l.recentPos.line, l.recentPos.col, l.recentLit, e)
	default:
		return fmt.Errorf(`parse error: line %d, column %d: %q: %s`, l.recentPos.line, l.recentPos.col, l.recentLit, e)
	}
}

func (l *lexer) emitError(e error) {
	select {
	case l.err <- e:
	case <-l.done:
	}
}

func (l *lexer) Error(s string) {
	l.emitError(l.makeError(s))
}
