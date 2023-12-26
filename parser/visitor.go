package parser

// Visitor ...
type Visitor interface {
	Interpret(Expr)
	visitRoot(Expr)
	visitQuery(Expr)
	visitFilter(Expr)
	visitIdentity(Expr)
	visitSelector(Expr)
	// visitIterator(Expr)
	// visitSpan(Expr)
	// visitString(Expr) string
	// visitInteger(Expr) string
}
