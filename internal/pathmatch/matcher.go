//go:generate goyacc -l -o parser.go parser.go.y

package pathmatch

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

type Values map[string]string

func (v Values) Get(s string) string {
	ret, ok := v[s]
	if !ok {
		return ""
	}
	return ret
}

type Matcher struct {
	consumers []consumer
}

type consumer interface {
	Consume(string, Values) (string, error)
}

type literalConsumer string

func (c literalConsumer) Consume(s string, _ Values) (string, error) {
	ps := strings.TrimPrefix(s, string(c))
	if ps == s {
		return s, fmt.Errorf(`failed to match literal pattern %q`, c)
	}
	return ps, nil
}

type segmentConsumer struct {
	name string
}

func (c *segmentConsumer) Consume(s string, mv Values) (string, error) {
	var val string
	// ([^/]+)/...
	// (lastsegment)
	i := strings.IndexByte(s, '/')
	if i == -1 {
		// it's not an error if we still have something left
		if len(s) > 0 {
			val = s
			s = ""
		} else {
			return s, fmt.Errorf(`failed to match segment %q`, c.name)
		}
	} else {
		val = s[:i]
		s = s[i:]
	}

	mv[c.name] = val
	return s, nil
}

type regexpConsumer struct {
	name    string
	pattern *regexp.Regexp
}

func (c *regexpConsumer) Consume(s string, mv Values) (string, error) {
	loc := c.pattern.FindStringIndex(s)
	if loc == nil {
		return s, fmt.Errorf(`failed to match pattern %q`, c.name)
	}

	// although technically our match only exists between loc[0] and loc[1],
	// we're going to need to remove all path components that matched it.
	// this means "from the beginning of the string to loc[1], but also
	// everything up to EOF or the next '/'

	var val string
	// Find next '/'
	i := strings.IndexByte(s[loc[1]:], '/')
	if i == -1 { // read up to EOF
		val = s
		s = ""
	} else {
		val = s[:i+loc[1]]
		s = s[i+loc[1]:]
	}
	mv[c.name] = val
	return s, nil
}

func Parse(s string) (*Matcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	exprs, err := parse(ctx, s)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse path pattern: %w`, err)
	}

	consumers := make([]consumer, 0, len(exprs))
	for _, expr := range exprs {
		switch expr := expr.(type) {
		case *Literal:
			consumers = append(consumers, literalConsumer(expr.Lit))
		case *LiteralPattern:
			consumers = append(consumers, &segmentConsumer{
				name: expr.Name,
			})
		case *RegexpPattern:
			pat, err := regexp.Compile(expr.Pattern)
			if err != nil {
				return nil, fmt.Errorf(`failed to compile pattern for %q: %w`, expr.Name, err)
			}
			consumers = append(consumers, &regexpConsumer{
				name:    expr.Name,
				pattern: pat,
			})
		default:
			return nil, fmt.Errorf(`invalid expression %T`, expr)
		}
	}
	return &Matcher{
		consumers: consumers,
	}, nil
}

func (p *Matcher) Match(s string) (Values, error) {
	mv := make(Values)
	for _, c := range p.consumers {
		ps, err := c.Consume(s, mv)
		if err != nil {
			return nil, fmt.Errorf(`failed to match input: %w`, err)
		}
		s = ps
		if s == "" {
			break
		}
	}
	// we can't have anything unprocessed
	if s != "" {
		return nil, fmt.Errorf(`failed to match input (trailing input)`)
	}

	return mv, nil
}
