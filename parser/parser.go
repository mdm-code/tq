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
	"errors"
	"fmt"
	"strings"

	"github.com/mdm-code/tq/lexer"
)

// Visitor ...
type Visitor interface {
	Print(Expr) string
	visitRoot(Expr) string
	visitQuery(Expr) string
	visitFilter(Expr) string
	visitIdentity(Expr) string
	visitSelector(Expr) string
	visitIterator(Expr) string
	visitSpan(Expr) string
	visitString(Expr) string
	visitInteger(Expr) string
}

// AstPrinter ...
type AstPrinter struct{}

// Print ...
func (a AstPrinter) Print(e Expr) string {
	return e.accept(a)
}

func (a AstPrinter) visitRoot(e Expr) string {
	switch v := e.(type) {
	case *Root:
		return a.parenthesize("root", v.query)
	default:
		// error out
	}
	return ""
}

func (a AstPrinter) visitQuery(e Expr) string {
	switch v := e.(type) {
	case *Query:
		return a.parenthesize("query", v.filters...)
	default:
		// error out
	}
	return ""
}

func (a AstPrinter) visitFilter(e Expr) string {
	switch v := e.(type) {
	case *Filter:
		return a.parenthesize("filter", v.kind)
	default:
		// error out
	}
	return ""
}

func (a AstPrinter) visitIdentity(e Expr) string {
	switch e.(type) {
	case *Identity:
		return "(identity)"
	default:
		// error out
	}
	return ""
}

func (a AstPrinter) visitSelector(e Expr) string {
	switch v := e.(type) {
	case *Selector:
		return a.parenthesize("selector", v.value)
	default:
		// error out
	}
	return ""
}

func (a AstPrinter) visitIterator(e Expr) string {
	switch e.(type) {
	case *Iterator:
		return "(iterator)"
	default:
		// error out
	}
	return ""
}

func (a AstPrinter) visitSpan(e Expr) string {
	switch v := e.(type) {
	case *Span:
		return a.parenthesize("span", v.left, v.right)
	default:
		// error out
	}
	return ""
}

func (a AstPrinter) visitInteger(e Expr) string {
	switch v := e.(type) {
	case *Integer:
		return fmt.Sprintf("(integer %v)", v.value)
	default:
		// error out
	}
	return ""
}

func (a AstPrinter) visitString(e Expr) string {
	switch v := e.(type) {
	case *String:
		return fmt.Sprintf("(string %v)", v.value)
	default:
		// error out
	}
	return ""
}

func (a AstPrinter) parenthesize(name string, es ...Expr) string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(name)
	for _, e := range es {
		b.WriteString(" ")
		if e == nil {
			b.WriteString("(null)")
			continue
		}
		b.WriteString(e.accept(a))
	}
	b.WriteString(")")
	return b.String()
}

// Root ...
type Root struct {
	query Expr
}

func (r *Root) accept(v Visitor) string {
	return v.visitRoot(r)
}

// Query ...
type Query struct {
	filters []Expr
}

func (q *Query) accept(v Visitor) string {
	return v.visitQuery(q)
}

// Filter ...
type Filter struct {
	kind Expr
}

func (f *Filter) accept(v Visitor) string {
	return v.visitFilter(f)
}

// Identity ...
type Identity struct{}

func (i *Identity) accept(v Visitor) string {
	return v.visitIdentity(i)
}

// Selector ...
type Selector struct {
	value Expr
}

func (s *Selector) accept(v Visitor) string {
	return v.visitSelector(s)
}

// Span ...
type Span struct {
	left, right Expr
}

func (s *Span) accept(v Visitor) string {
	return v.visitSpan(s)
}

// Iterator ...
type Iterator struct{}

func (i *Iterator) accept(v Visitor) string {
	return v.visitIterator(i)
}

// String ...
type String struct {
	value string
}

func (s *String) accept(v Visitor) string {
	return v.visitString(s)
}

// Integer ...
type Integer struct {
	value string
}

func (i *Integer) accept(v Visitor) string {
	return v.visitInteger(i)
}

// Parser ...
type Parser struct {
	Buffer  []lexer.Token
	Current int
}

// Expr ...
type Expr interface {
	accept(v Visitor) string
}

// New ...
func New(l *lexer.Lexer) (*Parser, error) {
	buf := []lexer.Token{}
	buf, ok := l.ScanAll(true)
	if !ok {
		err := errors.Join(l.Errors...)
		return nil, err
	}
	p := Parser{
		Buffer:  buf,
		Current: 0,
	}
	return &p, nil
}

// Parse ...
func (p *Parser) Parse() (Expr, error) {
	return p.root()
}

func (p *Parser) root() (Expr, error) {
	q, err := p.query()
	e := &Root{query: q}
	return e, err
}

func (p *Parser) query() (Expr, error) {
	var err error
	var f Expr
	e := Query{}
	for !p.isAtEnd() {
		switch {
		case p.match(lexer.Dot):
			f, err = p.identity()
			e.filters = append(e.filters, f)
			if err != nil {
				return &e, err
			}
		case p.match(lexer.ArrayOpen):
			f, err = p.selector()
			e.filters = append(e.filters, f)
			if err != nil {
				return &e, err
			}
		default:
			return &e, fmt.Errorf("query error at: %v", p.previous())
		}
	}
	return &e, err
}

func (p *Parser) identity() (Expr, error) {
	return &Identity{}, nil
}

func (p *Parser) selector() (Expr, error) {
	e := Selector{}
	if p.match(lexer.String) {
		e.value = &String{value: p.previous().Lexeme()}
		_, err := p.consume(lexer.ArrayClose)
		return &e, err
	}
	if p.match(lexer.ArrayClose) {
		e.value = &Iterator{}
		return &e, nil
	}
	if p.match(lexer.Colon) {
		s := Span{}
		if p.match(lexer.Integer) {
			r := Integer{value: p.previous().Lexeme()}
			s.right = &r
		}
		e.value = &s
		_, err := p.consume(lexer.ArrayClose)
		return &e, err
	}
	if p.match(lexer.Integer) {
		l := Integer{value: p.previous().Lexeme()}
		if p.match(lexer.ArrayClose) {
			e.value = &l
			return &e, nil
		}
		if p.match(lexer.Colon) {
			if p.match(lexer.Integer) {
				r := Integer{value: p.previous().Lexeme()}
				e.value = &Span{left: &l, right: &r}
				_, err := p.consume(lexer.ArrayClose)
				return &e, err
			}
			e.value = &Span{left: &l}
			_, err := p.consume(lexer.ArrayClose)
			return &e, err
		}
	}
	return &e, fmt.Errorf("parser error at: %v", p.previous())
}

func (p *Parser) consume(t lexer.TokenType) (lexer.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	return lexer.Token{}, fmt.Errorf("consume error at: %v", p.previous())
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
