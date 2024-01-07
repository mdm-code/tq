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

// Check the string integer value conversion to the proper integer.
func TestIntegerVtoi(t *testing.T) {
	cases := []struct {
		value string
		want  int
	}{
		{"0", 0},
		{"2", 2},
		{"12", 12},
		{"67", 67},
		{"99", 99},
	}
	for _, c := range cases {
		t.Run(c.value, func(t *testing.T) {
			i := Integer{c.value}
			have, err := i.Vtoi()
			if err != nil {
				t.Fatal(err)
			}
			if have != c.want {
				t.Errorf("have: %d; want: %d", have, c.want)
			}
		})
	}
}
