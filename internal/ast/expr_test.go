package ast

import "testing"

type mockVisitor struct{}

func (mockVisitor) VisitRoot(e Expr)     {}
func (mockVisitor) VisitQuery(e Expr)    {}
func (mockVisitor) VisitFilter(e Expr)   {}
func (mockVisitor) VisitIdentity(e Expr) {}
func (mockVisitor) VisitSelector(e Expr) {}
func (mockVisitor) VisitIterator(e Expr) {}
func (mockVisitor) VisitSpan(e Expr)     {}
func (mockVisitor) VisitString(e Expr)   {}
func (mockVisitor) VisitInteger(e Expr)  {}

// Test the Accept public method required by the visitor design pattern.
func TestAccept(t *testing.T) {
	cases := []struct {
		name string
		expr Expr
	}{
		{"root", &Root{}},
		{"query", &Query{}},
		{"filter", &Filter{}},
		{"identity", &Identity{}},
		{"selector", &Selector{}},
		{"iterator", &Iterator{}},
		{"span", &Span{}},
		{"string", &String{}},
		{"integer", &Integer{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var v mockVisitor
			c.expr.Accept(v)
		})
	}
}

func TestSpanLeftRight(t *testing.T) {}
