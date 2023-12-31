package ast

import (
	"strconv"
	"strings"
)

// Expr defines the expression interface for the visitor to operate on the
// contents of the expression node.
type Expr interface {
	Accept(v Visitor)
}

// Root stands for the top-level root node of the tq query. This version of
// the tq parser allows a single query, but extending the root to span
// multiple queries in always an option.
type Root struct {
	Query Expr
}

// Query represents a single tq query that can be run against a deserialized
// TOML data object. It potentially comprises of zero or more filters used to
// filter the TOML data. Although filters are stored in a slice implying a
// sequence, the order in not enforced neither by the expression nor the
// parser. It is the responsibility of the visiting interpreter run against the
// AST to provide the filering mechanism.
type Query struct {
	Filters []Expr
}

// Filter stands for a single tq filter. It the fundamental building block of
// the tq query.
type Filter struct {
	Kind Expr
}

// Identity specifies the identity data transformation that returns the
// filtered data argument unchanged.
type Identity struct{}

// Selector represents a select-driven data filter.
type Selector struct {
	Value Expr
}

// Span represents a filter that takes a slice of a list-like sequence.
type Span struct {
	Left, Right *Integer
}

// Iterator represents a sequeced iterator. The implementation of the iterator
// for TOML data types is to be provided by the visiting interpreter.
type Iterator struct{}

// String represents the key selector that can be used, for instance, in a form
// of dictionary lookup.
type String struct {
	Value string
}

// Integer represents your everyday integer. It can be used, for example, as an
// index of a data point in a sequence or a start/stop index of a span.
type Integer struct {
	Value string
}

func (r *Root) Accept(v Visitor) {
	v.VisitRoot(r)
}

func (q *Query) Accept(v Visitor) {
	v.VisitQuery(q)
}

func (f *Filter) Accept(v Visitor) {
	v.VisitFilter(f)
}

func (i *Identity) Accept(v Visitor) {
	v.VisitIdentity(i)
}

func (s *Selector) Accept(v Visitor) {
	v.VisitSelector(s)
}

func (s *Span) Accept(v Visitor) {
	v.VisitSpan(s)
}

// GetLeft returns the value of the left-hand side expression node of the Span.
func (s *Span) GetLeft(def int) int {
	return s.asInt(s.Left, def)
}

// GetRight returns the value of the right-hand side expression node of the Span.
func (s *Span) GetRight(def int) int {
	return s.asInt(s.Right, def)
}

func (s *Span) asInt(i *Integer, def int) int {
	var result = def
	if i != nil {
		integer, err := strconv.Atoi(i.Value)
		if err != nil {
			return result
		}
		result = integer
	}
	return result
}

func (i *Iterator) Accept(v Visitor) {
	v.VisitIterator(i)
}

func (s *String) Accept(v Visitor) {
	v.VisitString(s)
}

func (s *String) Trim() string {
	result := s.Value
	for _, c := range `'"` {
		result = strings.Trim(result, string(c))
	}
	return result
}

func (i *Integer) Accept(v Visitor) {
	v.VisitInteger(i)
}

func (i *Integer) Vtoi() (int, error) {
	return strconv.Atoi(i.Value)
}
