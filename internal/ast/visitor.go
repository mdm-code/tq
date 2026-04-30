package ast

// Visitor declares the interface for the AST visitor class. It declares
// signatures invoked by respective AST expression nodes.
type Visitor interface {
	VisitRoot(*Root)
	VisitQuery(*Query)
	VisitFilter(*Filter)
	VisitIdentity(*Identity)
	VisitSelector(*Selector)
	VisitIterator(*Iterator)
	VisitSpan(*Span)
	VisitString(*String)
	VisitInteger(*Integer)
}
