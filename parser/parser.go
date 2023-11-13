/*
Package parser ...

GRAMMAR
-------
root        = query
query       = *(filter)
filter      = identity / selector
identity    = DOT
selector    = '[' *( STRING / INTEGER / span ) ']'
span        = *INTEGER ':' *INTEGER
*/
package parser

import (
	"fmt"

	"github.com/mdm-code/tq/lexer"
)

// Acceptor ...
type Acceptor struct{}

func (a Acceptor) accept() {}

// Root ...
type Root struct {
	Acceptor
	query Expr
}

func (r Root) String() string {
	return fmt.Sprintf("root: %v", r.query)
}

// Query ...
type Query struct {
	Acceptor
	filters []Expr
}

func (q Query) String() string {
	return fmt.Sprintf("query: %v", q.filters)
}

// Filter ...
type Filter struct {
	Acceptor
	kind Expr
}

func (f Filter) String() string {
	return fmt.Sprintf("filter: %v", f.kind)
}

// Identity ...
type Identity struct {
	Acceptor
}

func (i Identity) String() string {
	return fmt.Sprintf("identity")
}

// Selector ...
type Selector struct {
	Acceptor
	value Expr
}

func (s Selector) String() string {
	return fmt.Sprintf("selector: %v", s.value)
}

// Span ...
type Span struct {
	Acceptor
	left, right Expr
}

func (s Span) String() string {
	return fmt.Sprintf("span: %v:%v", s.left, s.right)
}

// Iterator ...
type Iterator struct {
	Acceptor
}

func (i Iterator) String() string {
	return fmt.Sprintf("iterator")
}

// String ...
type String struct {
	Acceptor
	value string
}

func (s String) String() string { return fmt.Sprintf("str: %v", s.value) }

// Integer ...
type Integer struct {
	Acceptor
	value string
}

func (i Integer) String() string { return fmt.Sprintf("int: %v", i.value) }

// Parser ...
type Parser struct {
	Buffer  []lexer.Token
	Current int
}

// Expr ...
type Expr interface {
	accept()
}

// New ...
func New(l *lexer.Lexer, ignoreWhitespace bool) (*Parser, error) {
	buf := []lexer.Token{}
	for l.Next() {
		if ignoreWhitespace && l.Token().Type == lexer.Whitespace {
			continue
		}
		buf = append(buf, l.Token())
	}
	p := Parser{
		Buffer:  buf,
		Current: 0,
	}
	return &p, nil
}

func (p *Parser) root() Expr {
	e := Root{
		query: p.query(),
	}
	return e
}

func (p *Parser) query() Expr {
	e := Query{}
	for !p.isAtEnd() {
		if p.match(lexer.Dot) {
			f := p.identity()
			e.filters = append(e.filters, f)
			continue
		}
		if p.match(lexer.ArrayOpen) {
			f := p.selector()
			e.filters = append(e.filters, f)
			continue
		}
	}
	return e
}

func (p *Parser) identity() Expr {
	return Identity{}
}

func (p *Parser) selector() Expr {
	e := Selector{}
	if p.match(lexer.String) {
		e.value = String{value: p.previous().Lexeme()}
		p.consume(lexer.ArrayClose)
		return e
	}
	if p.match(lexer.Colon) {
		s := Span{}
		if p.match(lexer.Integer) {
			r := Integer{value: p.previous().Lexeme()}
			s.right = r
		}
		e.value = s
		p.consume(lexer.ArrayClose)
		return e
	}
	if p.match(lexer.Integer) {
		l := Integer{value: p.previous().Lexeme()}
		if p.match(lexer.ArrayClose) {
			e.value = l
			return e
		}
		if p.match(lexer.Colon) {
			if p.match(lexer.Integer) {
				r := Integer{value: p.previous().Lexeme()}
				e.value = Span{left: l, right: r}
				p.consume(lexer.ArrayClose)
				return e
			}
			e.value = Span{left: l}
			p.consume(lexer.ArrayClose)
			return e
		}
	}
	e.value = Iterator{}
	p.consume(lexer.ArrayClose)
	return e
}

func (p *Parser) consume(t lexer.TokenType) (lexer.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return lexer.Token{}, fmt.Errorf("can't consume")
}

func (p *Parser) match(tt ...lexer.TokenType) bool {
	for _, t := range tt {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	if p.Current > len(p.Buffer)-1 {
		return true
	}
	return false
}

func (p *Parser) previous() lexer.Token {
	return p.Buffer[p.Current-1]
}

func (p *Parser) peek() lexer.Token {
	return p.Buffer[p.Current]
}
