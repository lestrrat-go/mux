// Syntax:
//
// path: expr ...
// expr: literal | pattern
// pattern: patname (colon pattern)

%{
package pathmatch

import (
	"context"
	"fmt"
)

var _ = fmt.Printf

type Expression interface{}

func parse(ctx context.Context, s string) ([]Expression, error) {
	l := newLexer(ctx, s)
	if yyRet := yyParse(l); yyRet != 0 {
		return nil, <-l.err
	}
	return l.exprs, nil
}

%}

%union{
	token *yyToken
	expr  Expression
}

%type<expr> path
%type<expr> pattern
%type<expr> exprs
%type<expr> expr
%token<token> tLiteral tOpenBrace tCloseBrace tColon

%%

path
	: exprs

exprs
	: expr
	{
		$$ = $1
		if l, ok := yylex.(*lexer); ok {
			l.exprs = append([]Expression{$$}, l.exprs...)
		}
	}
	| expr exprs
	{
		$$ = $1
		if l, ok := yylex.(*lexer); ok {
			l.exprs = append([]Expression{$$}, l.exprs...)
		}
	}

expr
	: tOpenBrace pattern tCloseBrace
	{
		$$ = $2
	}
	| tLiteral
	{
		$$ = NewLiteral($1.lit.(string))
	}

pattern
	: tLiteral tColon tLiteral
	{
		$$ = NewRegexpPattern($1.lit.(string), $3.lit.(string))
	}
	| tLiteral
	{
		$$ = NewLiteralPattern($1.lit.(string))
	}

%%
