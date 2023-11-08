// Package lexer ...
//
// It's difficult but let's decide that all keys should be quoted. This makes
// allowed toml only-digits keys explicit. Dotted keys in toml will be handled
// with ["foo"]["bar"] syntax of a query.
package lexer

import (
	"errors"
	"strings"

	"github.com/mdm-code/scanner"
)

const (
	// String ...
	String TokenType = iota

	// Integer ...
	Integer

	// Dot ...
	Dot

	// Comma ...
	Comma

	// Colon ...
	Colon

	// ArrayOpen ...
	ArrayOpen

	// ArrayClose ...
	ArrayClose

	// Whitespace ...
	Whitespace

	// Undef ...
	Undef
)

// KeyCharMap ...
var KeyCharMap = map[rune]TokenType{
	'.': Dot,
	',': Comma,
	':': Colon,
	'[': ArrayOpen,
	']': ArrayClose,
}

// ErrNilScanner ...
var ErrNilScanner = errors.New("provided Scanner is nil")

// ErrKeyCharUnsupported ...
var ErrKeyCharUnsupported = errors.New("unsupported key character")

// ErrUnterminatedString ...
var ErrUnterminatedString = errors.New("unterminated string literal")

// ErrDisallowedChar ...
var ErrDisallowedChar = errors.New("disallowed character")

// TokenType ...
type TokenType uint8

// Token ...
type Token struct {
	Buffer     *[]scanner.Token
	Type       TokenType
	Lexeme     string
	Start, End int
}

// Lexer ...
type Lexer struct {
	Buffer []scanner.Token
	Errors []error
	Offset int
	Curr   Token
}

// New ...
func New(s *scanner.Scanner) (*Lexer, error) {
	if s == nil {
		return nil, ErrNilScanner
	}
	buf := []scanner.Token{}
	buf, ok := s.ScanAll()
	if !ok {
		err := errors.Join(s.Errors...)
		return nil, err
	}
	l := Lexer{
		Offset: 0,
		Buffer: buf,
		Curr: Token{
			Buffer: nil,
			Type:   Undef,
			Lexeme: "",
			Start:  0,
			End:    0,
		},
	}
	return &l, nil
}

// Token ...
func (l *Lexer) Token() Token {
	return l.Curr
}

// Next ...
func (l *Lexer) Next() bool {
	if l.Offset > len(l.Buffer)-1 {
		return false
	}
	t := l.Buffer[l.Offset]
	switch r := t.Rune; {
	case IsKeyChar(r):
		return l.nextKeyChar()
	case IsQuote(r):
		return l.nextString()
	case IsDigit(r):
		return l.nextInteger()
	case IsWhitespace(r):
		return l.nextWhitespace()
	default:
		l.Errors = append(l.Errors, ErrDisallowedChar)
		return false
	}
}

func (l *Lexer) nextKeyChar() bool {
	t := l.Buffer[l.Offset]
	tp, ok := KeyCharMap[t.Rune]
	if !ok {
		l.Errors = append(l.Errors, ErrKeyCharUnsupported)
		return false
	}
	l.Curr = Token{
		Buffer: &l.Buffer,
		Type:   tp,
		Lexeme: string(t.Rune),
		Start:  l.Offset,
		End:    l.Offset + 1,
	}
	if ok {
		l.Offset++
	}
	return true
}

func (l *Lexer) nextString() bool {
	t := l.Buffer[l.Offset]
	tq := t.Rune
	ss := []string{}
	ss = append(ss, string(t.Rune))
	start := l.Offset
	l.Offset++
	for {
		if l.Offset > len(l.Buffer)-1 {
			lx := strings.Join(ss, "")
			l.Curr = Token{
				Buffer: &l.Buffer,
				Type:   Undef,
				Lexeme: lx,
				Start:  start,
				End:    l.Offset + 1,
			}
			l.Errors = append(l.Errors, ErrUnterminatedString)
			return false
		}
		t = l.Buffer[l.Offset]
		if IsNewline(t.Rune) {
			lx := strings.Join(ss, "")
			l.Curr = Token{
				Buffer: &l.Buffer,
				Type:   Undef,
				Lexeme: lx,
				Start:  start,
				End:    l.Offset + 1,
			}
			l.Errors = append(l.Errors, ErrDisallowedChar)
			return false
		}
		if t.Rune == tq {
			ss = append(ss, string(t.Rune))
			l.Offset++
			break
		}
		ss = append(ss, string(t.Rune))
		l.Offset++
	}
	lx := strings.Join(ss, "")
	l.Curr = Token{
		Buffer: &l.Buffer,
		Type:   String,
		Lexeme: lx,
		Start:  start,
		End:    l.Offset,
	}
	return true
}

func (l *Lexer) nextInteger() bool {
	t := l.Buffer[l.Offset]
	ts := []string{}
	ts = append(ts, string(t.Rune))
	start := l.Offset
	l.Offset++
	for {
		if l.Offset > len(l.Buffer)-1 {
			break
		}
		t = l.Buffer[l.Offset]
		if !IsDigit(t.Rune) {
			break
		}
		ts = append(ts, string(t.Rune))
		l.Offset++
	}
	lx := strings.Join(ts, "")
	l.Curr = Token{
		Buffer: &l.Buffer,
		Type:   Integer,
		Lexeme: lx,
		Start:  start,
		End:    l.Offset,
	}
	return true
}

func (l *Lexer) nextWhitespace() bool {
	t := l.Buffer[l.Offset]
	ss := []string{}
	ss = append(ss, string(t.Rune))
	start := l.Offset
	l.Offset++
	for {
		if l.Offset > len(l.Buffer)-1 {
			break
		}
		t = l.Buffer[l.Offset]
		if !IsWhitespace(t.Rune) {
			break
		}
		ss = append(ss, string(t.Rune))
		l.Offset++
	}
	lx := strings.Join(ss, "")
	l.Curr = Token{
		Buffer: &l.Buffer,
		Type:   Whitespace,
		Lexeme: lx,
		Start:  start,
		End:    l.Offset,
	}
	return true
}
