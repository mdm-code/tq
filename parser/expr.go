package parser

import (
	"strconv"
	"strings"
)

// Expr defines the expression interface for the visitor to operate on the
// contents of the expression node.
type Expr interface {
	accept(v ASTVisitor)
}

// Root stands for the top-level root node of the tq query. This version of
// the tq parser allows a single query, but extending the root to span
// multiple queries in always an option.
type Root struct {
	query Expr
}

// Query represents a single tq query that can be run against a deserialized
// TOML data object. It potentially comprises of zero or more filters used to
// filter the TOML data. Although filters are stored in a slice implying a
// sequence, the order in not enforced neither by the expression nor the
// parser. It is the responsibility of the visiting interpreter run against the
// AST to provide the filering mechanism.
type Query struct {
	filters []Expr
}

// Filter stands for a single tq filter. It the fundamental building block of
// the tq query.
type Filter struct {
	kind Expr
}

// Identity specifies the identity data transformation that returns the
// filtered data argument unchanged.
type Identity struct{}

// Selector represents a select-driven data filter.
type Selector struct {
	value Expr
}

// Span represents a filter that takes a slice of a list-like sequence.
type Span struct {
	left, right *Integer
}

// Iterator represents a sequeced iterator. The implementation of the iterator
// for TOML data types is to be provided by the visiting interpreter.
type Iterator struct{}

// String represents the key selector that can be used, for instance, in a form
// of dictionary lookup.
type String struct {
	value string
}

// Integer represents your everyday integer. It can be used, for example, as an
// index of a data point in a sequence or a start/stop index of a span.
type Integer struct {
	value string
}

func (r *Root) accept(v ASTVisitor) {
	v.visitRoot(r)
}

func (q *Query) accept(v ASTVisitor) {
	v.visitQuery(q)
}

func (f *Filter) accept(v ASTVisitor) {
	v.visitFilter(f)
}

func (i *Identity) accept(v ASTVisitor) {
	v.visitIdentity(i)
}

func (s *Selector) accept(v ASTVisitor) {
	v.visitSelector(s)
}

func (s *Span) accept(v ASTVisitor) {
	v.visitSpan(s)
}

// Left returns the value of the left-hand side expression node of the Span.
func (s *Span) Left(def int) int {
	return s.asInt(s.left, def)
}

// Right returns the value of the right-hand side expression node of the Span.
func (s *Span) Right(def int) int {
	return s.asInt(s.right, def)
}

func (s *Span) asInt(i *Integer, def int) int {
	var result = def
	if i != nil {
		integer, err := strconv.Atoi(i.value)
		if err != nil {
			return result
		}
		result = integer
	}
	return result
}

func (i *Iterator) accept(v ASTVisitor) {
	v.visitIterator(i)
}

func (s *String) accept(v ASTVisitor) {
	v.visitString(s)
}

func (s *String) trimmed() string {
	result := s.value
	for _, c := range `'"` {
		result = strings.Trim(result, string(c))
	}
	return result
}

func (i *Integer) accept(v ASTVisitor) {
	v.visitInteger(i)
}

func (i *Integer) vtoi() (int, error) {
	return strconv.Atoi(i.value)
}
