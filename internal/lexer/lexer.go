// Package lexer provides a Lexer struct that converts a string of scanned
// characters into allowed tq lexemes.
package lexer

import (
	"errors"

	"github.com/mdm-code/scanner"
)

const (
	// lexerOffsetStart declares the initial Lexer offset.
	lexerOffsetStart = 0

	// lexerLineOffsetStart declares the initial Lexer line offset.
	lexerLineOffsetStart = 0
)

// Lexer is the struct that tokenizes character input into tq lexemes.
type Lexer struct {
	buffer     []scanner.Token
	Errors     []error // errors encountered in the course of lexer execution
	offset     int
	lineOffset int
	curr       Token
}

// New returns a new Lexer with its buffer populated with scanner tokens read
// from the Scanner s.
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
		offset:     lexerOffsetStart,
		lineOffset: lexerLineOffsetStart,
		buffer:     buf,
		curr: Token{
			Buffer: nil,
			Type:   Undefined,
			Start:  0,
			End:    0,
		},
	}
	return &l, nil
}

// Token return the most recently scanned token.
func (l *Lexer) Token() Token {
	return l.curr
}

// Scan attempts to scan the next lexer token from the internal buffer.
//
// It returns true if the scan succeeds and false in case it fails.
func (l *Lexer) Scan() bool {
	if l.offset > len(l.buffer)-1 {
		return false
	}
	t := l.buffer[l.offset]
	switch r := t.Rune; {
	case isKeyChar(r):
		return l.scanKeyChar()
	case isQuote(r):
		return l.scanString()
	case isDigit(r):
		return l.scanInteger()
	case isBareChar(r):
		return l.scanBareString()
	case isWhitespace(r):
		return l.scanWhitespace()
	default:
		l.setToken(Undefined, l.offset, l.offset+1)
		l.pushErr(ErrDisallowedChar)
		return false
	}
}

// ScanAll attempts to scann all lexer tokens from the internal buffer.
//
// It has ignoreWhitespace boolean parameter that controls if white space
// tokens are to be ignored in the output slice of tokens. In case errors are
// encountered in the course of scanning over the buffer, it returns false as
// its second return value.
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

// Errored reports if Lexer accumulated errors in the course of its execution.
//
// The Lexer field Errors holds a slice of errors that occurred during its
// execution.
func (l *Lexer) Errored() bool {
	return len(l.Errors) > 0
}

func (l *Lexer) advance() {
	l.offset++
	l.lineOffset++
}

func (l *Lexer) resetLineOffsetOnLineBreak(r rune) {
	if isLineBreak(r) {
		l.lineOffset = lexerLineOffsetStart
	}
}

func (l *Lexer) setToken(tp TokenType, start, end int) {
	l.curr = Token{
		Buffer:     &l.buffer,
		Type:       tp,
		Start:      start,
		End:        end,
		LineOffset: l.lineOffset,
	}
}

func (l *Lexer) pushErr(err error) {
	e := Error{
		buffer:     &l.buffer,
		offset:     l.offset,
		lineOffset: l.lineOffset,
		err:        err,
	}
	l.Errors = append(l.Errors, &e)
}

func (l *Lexer) scanKeyChar() bool {
	t := l.buffer[l.offset]
	tp, ok := keyCharMap[t.Rune]
	if !ok {
		l.pushErr(ErrKeyCharUnsupported)
		return false
	}
	l.setToken(tp, l.offset, l.offset+1)
	l.advance()
	return true
}

func (l *Lexer) scanBareString() bool {
	t := l.buffer[l.offset]
	start := l.offset
	l.advance()
	for l.offset <= len(l.buffer)-1 {
		t = l.buffer[l.offset]
		if !isBareChar(t.Rune) {
			break
		}
		l.advance()
	}
	l.setToken(String, start, l.offset)
	return true
}

func (l *Lexer) scanString() bool {
	var prev scanner.Token
	t := l.buffer[l.offset]
	tq := t.Rune
	start := l.offset
	prev = t
	l.advance()
	for {
		if l.offset > len(l.buffer)-1 {
			// NOTE: This error is reported because the string goes past the
			// buffer without encountering the matching quote character that
			// should terminate the quoted string.
			l.setToken(Undefined, start, l.offset+1)
			l.pushErr(ErrUnterminatedString)
			return false
		}
		l.resetLineOffsetOnLineBreak(t.Rune)
		t = l.buffer[l.offset]
		if t.Rune == tq && prev.Rune != '\\' {
			prev = t
			l.advance()
			break
		}
		prev = t
		l.advance()
	}
	l.setToken(String, start, l.offset)
	return true
}

func (l *Lexer) scanInteger() bool {
	t := l.buffer[l.offset]
	start := l.offset
	l.advance()
	for l.offset <= len(l.buffer)-1 {
		t = l.buffer[l.offset]
		if !isDigit(t.Rune) {
			break
		}
		l.advance()
	}
	l.setToken(Integer, start, l.offset)
	return true
}

func (l *Lexer) scanWhitespace() bool {
	t := l.buffer[l.offset]
	start := l.offset
	l.advance()
	for l.offset <= len(l.buffer)-1 {
		l.resetLineOffsetOnLineBreak(t.Rune)
		t = l.buffer[l.offset]
		if !isWhitespace(t.Rune) {
			break
		}
		l.advance()
	}
	l.setToken(Whitespace, start, l.offset)
	return true
}
