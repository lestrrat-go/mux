package pathmatch

type Literal struct {
	Lit string
}

type LiteralPattern struct {
	Name string
}

type RegexpPattern struct {
	Name    string
	Pattern string
}

func NewLiteralPattern(s string) Expression {
	return &LiteralPattern{Name: s}
}

func NewRegexpPattern(name string, pattern string) Expression {
	return &RegexpPattern{
		Name:    name,
		Pattern: pattern,
	}
}

func NewLiteral(s string) Expression {
	return &Literal{
		Lit: s,
	}
}
