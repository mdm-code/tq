package parser

// Expr defines the expression interface for the visitor to operate on the
// contents of the expression node.
type Expr interface {
	accept(v Visitor)
}

// Root stands for the top-level root node of the tq query. This version of
// the tq parser allows a single query, but extending the root to span
// multiple queries in always an option.
type Root struct {
	query Expr
}

func (r *Root) accept(v Visitor) {
	v.visitRoot(r)
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

func (q *Query) accept(v Visitor) {
	v.visitQuery(q)
}

// Filter stands for a single tq filter. It the fundamental building block of
// the tq query.
type Filter struct {
	kind Expr
}

func (f *Filter) accept(v Visitor) {
	v.visitFilter(f)
}

// Identity specifies the identity data transformation that returns the
// filtered data argument unchanged.
type Identity struct{}

func (i *Identity) accept(v Visitor) {
	v.visitIdentity(i)
}

// Selector represents a select-driven data filter.
type Selector struct {
	value Expr
}

func (s *Selector) accept(v Visitor) {
	v.visitSelector(s)
}

// Span represents a filter that takes a slice of a list-like sequence.
type Span struct {
	left, right *Integer
}

func (s *Span) accept(v Visitor) {
	// v.visitSpan(s)
}

// Iterator represents a sequeced iterator. The implementation of the iterator
// for TOML data types is to be provided by the visiting interpreter.
type Iterator struct{}

func (i *Iterator) accept(v Visitor) {
	// v.visitIterator(i)
}

// String represents the key selector that can be used, for instance, in a form
// of dictionary lookup.
type String struct {
	value string
}

func (s *String) accept(v Visitor) {
	// v.visitString(s)
}

// Integer represents your everyday integer. It can be used, for example, as an
// index of a data point in a sequence or a start/stop index of a span.
type Integer struct {
	value string
}

func (i *Integer) accept(v Visitor) {
	// v.visitInteger(i)
}
