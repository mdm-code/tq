package parser

import (
	"errors"

	"github.com/mdm-code/tq/internal/lexer"
)

// Parser encapsulates the logic of parsing tq queries into valid expressions.
type Parser struct {
	buffer  []lexer.Token
	current int
}

// New returns a new Parser with the buffer populated with lexer tokens read
// from the Lexer l.
func New(l *lexer.Lexer) (*Parser, error) {
	buf := []lexer.Token{}
	buf, ok := l.ScanAll(true)
	if !ok {
		err := errors.Join(l.Errors...)
		return nil, err
	}
	p := Parser{
		buffer:  buf,
		current: 0,
	}
	return &p, nil
}

// Parse the abstract syntax tree given the buffer of tq lexer tokens.
func (p *Parser) Parse() (*Root, error) {
	root, err := p.root()
	return &root, err
}

func (p *Parser) root() (Root, error) {
	q, err := p.query()
	expr := Root{query: &q}
	return expr, err
}

func (p *Parser) query() (Query, error) {
	var expr Query
	var err error
	for !p.isAtEnd() {
		var f Filter
		f, err = p.filter()
		expr.filters = append(expr.filters, &f)
		if err != nil {
			break
		}
	}
	return expr, err
}

func (p *Parser) filter() (Filter, error) {
	var expr Filter
	var err error
	switch {
	case p.match(lexer.Dot):
		var i Identity
		i, err = p.identity()
		expr.kind = &i
	case p.match(lexer.ArrayOpen):
		var s Selector
		s, err = p.selector()
		expr.kind = &s
	default:
		prev := p.previous()
		err = &Error{prev.Lexeme(), ErrQueryElement}
	}
	return expr, err
}

func (p *Parser) identity() (Identity, error) {
	return Identity{}, nil
}

func (p *Parser) selector() (Selector, error) {
	var expr Selector
	var err error
	switch {
	case p.check(lexer.ArrayClose):
		i, _ := p.iterator()
		expr.value = &i
	case p.match(lexer.String):
		s, _ := p.string()
		expr.value = &s
	case p.match(lexer.Colon):
		s, _ := p.span(nil)
		expr.value = &s
	case p.match(lexer.Integer):
		i, _ := p.integer()
		if p.match(lexer.Colon) {
			s, _ := p.span(&i)
			expr.value = &s
		} else {
			expr.value = &i
		}
	}
	_, err = p.consume(lexer.ArrayClose, ErrSelectorUnterminated)
	return expr, err
}

func (p *Parser) iterator() (Iterator, error) {
	return Iterator{}, nil
}

func (p *Parser) string() (String, error) {
	return String{value: p.previous().Lexeme()}, nil
}

func (p *Parser) integer() (Integer, error) {
	return Integer{value: p.previous().Lexeme()}, nil
}

func (p *Parser) span(left *Integer) (Span, error) {
	s := Span{left: left}
	if p.match(lexer.Integer) {
		r, _ := p.integer()
		s.right = &r
	}
	return s, nil
}

func (p *Parser) consume(t lexer.TokenType, e error) (lexer.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	curr, err := p.peek()
	var lexeme string
	if err != nil && errors.Is(err, ErrParserBufferOutOfRange) {
		lexeme = "EOF"
	} else {
		lexeme = curr.Lexeme()
	}
	err = &Error{lexeme, e}
	return curr, err
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
	other, err := p.peek()
	if err != nil {
		return false
	}
	return other.Type == t
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	if p.current > len(p.buffer)-1 {
		return true
	}
	return false
}

func (p *Parser) previous() lexer.Token {
	return p.buffer[p.current-1]
}

func (p *Parser) peek() (lexer.Token, error) {
	if p.isAtEnd() {
		return p.previous(), ErrParserBufferOutOfRange
	}
	return p.buffer[p.current], nil
}
