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

	"github.com/mdm-code/tq/lexer"
)

// Error ...
type Error struct {
	Token lexer.Token
	Err   error
}

// Error ...
func (e Error) Error() string {
	if e.Err != nil {
		return e.Err.Error() + fmt.Sprintf(" but got '%s'", e.Token.Lexeme())
	}
	return "null"
}

// Visitor ...
type Visitor interface {
	Interpret(Expr)
	visitRoot(Expr)
	visitQuery(Expr)
	visitFilter(Expr)
	visitIdentity(Expr)
	visitSelector(Expr)
	visitIterator(Expr)
	// visitSpan(Expr)
	// visitString(Expr) string
	// visitInteger(Expr) string
}

// FilterFunc ...
type FilterFunc func(data interface{}) interface{}

// QueryConstructor ...
type QueryConstructor struct {
	Filters []FilterFunc
}

func (q *QueryConstructor) interpret(es ...Expr) {
	for _, e := range es {
		e.accept(q)
	}
}

// Interpret ...
func (q *QueryConstructor) Interpret(e Expr) {
	e.accept(q)
}

func (q *QueryConstructor) visitRoot(e Expr) {
	switch v := e.(type) {
	case *Root:
		q.interpret(v.query)
	default:
		// error out
	}
}

func (q *QueryConstructor) visitQuery(e Expr) {
	switch v := e.(type) {
	case *Query:
		q.interpret(v.filters...)
	default:
		// error out
	}
}

func (q *QueryConstructor) visitFilter(e Expr) {
	switch v := e.(type) {
	case *Filter:
		q.interpret(v.kind)
	default:
		// error out
	}
}
func (q *QueryConstructor) visitIdentity(e Expr) {
	switch v := e.(type) {
	case *Identity:
		fmt.Println(*v)
		q.Filters = append(q.Filters, func(data interface{}) interface{} { return data })
	default:
		// error out
	}
}

func (q *QueryConstructor) visitSelector(e Expr) {
	switch v := e.(type) {
	case *Selector:
		fmt.Println(*v)
		q.Filters = append(q.Filters, func(data interface{}) interface{} { return data })
	default:
		// error out
	}
}

func (q *QueryConstructor) visitIterator(e Expr) {
	switch v := e.(type) {
	case *Iterator:
		fmt.Println(*v)
		q.Filters = append(q.Filters, func(data interface{}) interface{} { return data })
	default:
		// error out
	}
}

// // AstPrinter ...
// type AstPrinter struct{}

// // Print ...
// func (a AstPrinter) Print(e Expr) string {
// 	return e.accept(a)
// }

// func (a AstPrinter) visitRoot(e Expr) string {
// 	switch v := e.(type) {
// 	case *Root:
// 		return a.parenthesize("root", v.query)
// 	default:
// 		// error out
// 	}
// 	return ""
// }

// func (a AstPrinter) visitQuery(e Expr) string {
// 	switch v := e.(type) {
// 	case *Query:
// 		return a.parenthesize("query", v.filters...)
// 	default:
// 		// error out
// 	}
// 	return ""
// }

// func (a AstPrinter) visitFilter(e Expr) string {
// 	switch v := e.(type) {
// 	case *Filter:
// 		return a.parenthesize("filter", v.kind)
// 	default:
// 		// error out
// 	}
// 	return ""
// }

// func (a AstPrinter) visitIdentity(e Expr) string {
// 	switch e.(type) {
// 	case *Identity:
// 		return "(identity)"
// 	default:
// 		// error out
// 	}
// 	return ""
// }

// func (a AstPrinter) visitSelector(e Expr) string {
// 	switch v := e.(type) {
// 	case *Selector:
// 		return a.parenthesize("selector", v.value)
// 	default:
// 		// error out
// 	}
// 	return ""
// }

// func (a AstPrinter) visitIterator(e Expr) string {
// 	switch e.(type) {
// 	case *Iterator:
// 		return "(iterator)"
// 	default:
// 		// error out
// 	}
// 	return ""
// }

// func (a AstPrinter) visitSpan(e Expr) string {
// 	switch v := e.(type) {
// 	case *Span:
// 		return a.parenthesize("span", v.left, v.right)
// 	default:
// 		// error out
// 	}
// 	return ""
// }

// func (a AstPrinter) visitInteger(e Expr) string {
// 	switch v := e.(type) {
// 	case *Integer:
// 		return fmt.Sprintf("(integer %v)", v.value)
// 	default:
// 		// error out
// 	}
// 	return ""
// }

// func (a AstPrinter) visitString(e Expr) string {
// 	switch v := e.(type) {
// 	case *String:
// 		return fmt.Sprintf("(string %v)", v.value)
// 	default:
// 		// error out
// 	}
// 	return ""
// }

// func (a AstPrinter) parenthesize(name string, es ...Expr) string {
// 	var b strings.Builder
// 	b.WriteString("(")
// 	b.WriteString(name)
// 	for _, e := range es {
// 		b.WriteString(" ")
// 		if e == nil {
// 			b.WriteString("(null)")
// 			continue
// 		}
// 		b.WriteString(e.accept(a))
// 	}
// 	b.WriteString(")")
// 	return b.String()
// }

// Root ...
type Root struct {
	query Expr
}

func (r *Root) accept(v Visitor) {
	v.visitRoot(r)
}

// Query ...
type Query struct {
	filters []Expr
}

func (q *Query) accept(v Visitor) {
	v.visitQuery(q)
}

// Filter ...
type Filter struct {
	kind Expr
}

func (f *Filter) accept(v Visitor) {
	v.visitFilter(f)
}

// Identity ...
type Identity struct{}

func (i *Identity) accept(v Visitor) {
	v.visitIdentity(i)
}

// Selector ...
type Selector struct {
	value Expr
}

func (s *Selector) accept(v Visitor) {
	v.visitSelector(s)
}

// Span ...
type Span struct {
	left, right Expr
}

func (s *Span) accept(v Visitor) {
	// v.visitSpan(s)
}

// Iterator ...
type Iterator struct{}

func (i *Iterator) accept(v Visitor) {
	v.visitIterator(i)
}

// String ...
type String struct {
	value string
}

func (s *String) accept(v Visitor) {
	// v.visitString(s)
}

// Integer ...
type Integer struct {
	value string
}

func (i *Integer) accept(v Visitor) {
	// v.visitInteger(i)
}

// Parser ...
type Parser struct {
	Buffer  []lexer.Token
	Current int
}

// Expr ...
type Expr interface {
	accept(v Visitor)
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
			err = Error{
				p.previous(),
				fmt.Errorf("expected '.' or '[' to parse query element"),
			}
			return &e, err
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
		_, err := p.consume(lexer.ArrayClose, "expected ']' to terminate selector")
		return &e, err
	}
	if p.match(lexer.Colon) {
		s := Span{}
		if p.match(lexer.Integer) {
			r := Integer{value: p.previous().Lexeme()}
			s.right = &r
		}
		e.value = &s
		_, err := p.consume(lexer.ArrayClose, "expected ']' to terminate selector")
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
				_, err := p.consume(lexer.ArrayClose, "expected ']' to terminate selector")
				return &e, err
			}
			e.value = &Span{left: &l}
			_, err := p.consume(lexer.ArrayClose, "expected ']' to terminate selector")
			return &e, err
		}
	}
	if p.match(lexer.ArrayClose) {
		e.value = &Iterator{}
		return &e, nil
	}
	err := Error{
		p.previous(),
		fmt.Errorf("expected ']' to terminate the selector"),
	}
	return &e, err
}

func (p *Parser) consume(t lexer.TokenType, msg string) (lexer.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	// NOTE: or possibly p.previous()
	err := Error{p.peek(), fmt.Errorf(msg)}
	return p.peek(), err
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
