/*
Package parser ...

GRAMMAR
-------
root        = query
query       = *(filter)
filter      = identity / selector
identity    = DOT
selector    = '[' *( STRING / INTEGER / range ) ']'
range       = *INTEGER ':' *INTEGER
*/
package parser

import "github.com/mdm-code/tq/lexer"

// Parser ...
type Parser struct {
	Buffer []lexer.Token
}

// New ...
func New(l *lexer.Lexer) (*Parser, error) {
	buf := []lexer.Token{}
	for l.Next() {
		buf = append(buf, l.Token())
	}
	p := Parser{
		Buffer: buf,
	}
	return &p, nil
}

func (p *Parser) parse() {
}
