package parser

// ASTVisitor declares the interface for the AST visitor class. It declares
// signatures invoked by respective AST expression nodes.
type ASTVisitor interface {
	visitRoot(Expr)
	visitQuery(Expr)
	visitFilter(Expr)
	visitIdentity(Expr)
	visitSelector(Expr)
	visitIterator(Expr)
	visitSpan(Expr)
	visitString(Expr)
	visitInteger(Expr)
}
