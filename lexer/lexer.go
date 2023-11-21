// Package lexer ...
package lexer

import (
	"errors"
	"fmt"
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

	// Colon ...
	Colon

	// ArrayOpen ...
	ArrayOpen

	// ArrayClose ...
	ArrayClose

	// Whitespace ...
	Whitespace

	// Undefined ...
	Undefined
)

// KeyCharMap ...
var KeyCharMap = map[rune]TokenType{
	'.': Dot,
	':': Colon,
	'[': ArrayOpen,
	']': ArrayClose,
}

var (
	// ErrNilScanner ...
	ErrNilScanner = errors.New("provided Scanner is nil")

	// ErrKeyCharUnsupported ...
	ErrKeyCharUnsupported = errors.New("unsupported key character")

	// ErrUnterminatedString ...
	ErrUnterminatedString = errors.New("unterminated string literal")

	// ErrDisallowedChar ...
	ErrDisallowedChar = errors.New("disallowed character")
)

// TokenType ...
type TokenType uint8

// Error ...
type Error struct {
	Buffer *[]scanner.Token
	Offset int
	Err    error
}

// Token ...
type Token struct {
	Buffer     *[]scanner.Token
	Type       TokenType
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
			Type:   Undefined,
			Start:  0,
			End:    0,
		},
	}
	return &l, nil
}

// Error ...
func (e *Error) Error() string {
	if e.Buffer == nil {
		return e.Err.Error()
	}
	var b strings.Builder
	b.Grow(e.Offset*2 + 1)
	for _, t := range *e.Buffer {
		b.WriteString(string(t.Rune))
	}
	b.WriteString("\n")
	for i := 0; i < e.Offset; i++ {
		b.WriteString(" ")
	}
	b.WriteString("^")
	b.WriteString("\n")
	var errMsg string
	if e.Err != nil {
		errMsg = e.Err.Error()
	} else {
		errMsg = "null"
	}
	b.WriteString(fmt.Sprintf("Lexer error: %s", errMsg))
	return b.String()
}

// Lexeme ...
func (t Token) Lexeme() string {
	if t.Buffer == nil {
		return ""
	}
	chars := make([]string, t.End-t.Start)
	for _, t := range (*t.Buffer)[t.Start:t.End] {
		chars = append(chars, string(t.Rune))
	}
	return strings.Join(chars, "")
}

// Token ...
func (l *Lexer) Token() Token {
	return l.Curr
}

// Scan ...
func (l *Lexer) Scan() bool {
	if l.Offset > len(l.Buffer)-1 {
		return false
	}
	t := l.Buffer[l.Offset]
	switch r := t.Rune; {
	case IsKeyChar(r):
		return l.scanKeyChar()
	case IsQuote(r):
		return l.scanString()
	case IsDigit(r):
		return l.scanInteger()
	case IsWhitespace(r):
		return l.scanWhitespace()
	default:
		l.setToken(Undefined, l.Offset, l.Offset+1)
		l.pushErr(ErrDisallowedChar, l.Offset)
		return false
	}
}

func (l *Lexer) scanKeyChar() bool {
	t := l.Buffer[l.Offset]
	tp, ok := KeyCharMap[t.Rune]
	if !ok {
		l.pushErr(ErrKeyCharUnsupported, l.Offset)
		return false
	}
	l.setToken(tp, l.Offset, l.Offset+1)
	if ok {
		l.advance()
	}
	return true
}

func (l *Lexer) scanString() bool {
	t := l.Buffer[l.Offset]
	tq := t.Rune
	start := l.Offset
	l.advance()
	for {
		if l.Offset > len(l.Buffer)-1 {
			l.setToken(Undefined, start, l.Offset+1)
			l.pushErr(ErrUnterminatedString, start)
			return false
		}
		t = l.Buffer[l.Offset]
		if IsNewline(t.Rune) {
			l.setToken(Undefined, start, l.Offset+1)
			l.pushErr(ErrDisallowedChar, start)
			return false
		}
		if t.Rune == tq {
			l.advance()
			break
		}
		l.advance()
	}
	l.setToken(String, start, l.Offset)
	return true
}

func (l *Lexer) scanInteger() bool {
	t := l.Buffer[l.Offset]
	start := l.Offset
	l.advance()
	for {
		if l.Offset > len(l.Buffer)-1 {
			break
		}
		t = l.Buffer[l.Offset]
		if !IsDigit(t.Rune) {
			break
		}
		l.advance()
	}
	l.setToken(Integer, start, l.Offset)
	return true
}

func (l *Lexer) scanWhitespace() bool {
	t := l.Buffer[l.Offset]
	start := l.Offset
	l.advance()
	for {
		if l.Offset > len(l.Buffer)-1 {
			break
		}
		t = l.Buffer[l.Offset]
		if !IsWhitespace(t.Rune) {
			break
		}
		l.advance()
	}
	l.setToken(Whitespace, start, l.Offset)
	return true
}

func (l *Lexer) setToken(tp TokenType, start, end int) {
	l.Curr = Token{
		Buffer: &l.Buffer,
		Type:   tp,
		Start:  start,
		End:    end,
	}
}

func (l *Lexer) pushErr(err error, offset int) {
	e := Error{
		Buffer: &l.Buffer,
		Offset: offset,
		Err:    err,
	}
	l.Errors = append(l.Errors, &e)
}

func (l *Lexer) advance() {
	l.Offset++
}

// ScanAll ...
func (l *Lexer) ScanAll(ignoreWhitespace bool) ([]Token, bool) {
	result := []Token{}
	for l.Scan() {
		if ignoreWhitespace && l.Token().Type == Whitespace {
			continue
		}
		t := l.Token()
		result = append(result, t)
	}
	if l.Errored() {
		return result, false
	}
	return result, true
}

// Errored ...
func (l *Lexer) Errored() bool {
	return len(l.Errors) > 0
}
