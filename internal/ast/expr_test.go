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

// Verify if left int value of the span is retrieved.
func TestSpanGetLeft(t *testing.T) {
	cases := []struct {
		name     string
		intValue *Integer
		want     int
	}{
		{
			name:     "0",
			intValue: &Integer{Value: "0"},
			want:     0,
		},
		{
			name:     "99",
			intValue: &Integer{Value: "99"},
			want:     99,
		},
		{
			// nil Integer pointer results in the default def value returned
			name:     "nil-integer",
			intValue: nil,
			want:     0,
		},
		{
			// non-convertable string results in the default def value returned
			name:     "non-convertable",
			intValue: &Integer{Value: "non-convertable"},
			want:     0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := Span{Left: c.intValue}
			def := 0
			have := s.GetLeft(def)
			if have != c.want {
				t.Errorf("have: %d; want: %d", have, c.want)
			}
		})
	}
}

// Verify if right int value of the span is retrieved.
func TestSpanGetRight(t *testing.T) {
	cases := []struct {
		name     string
		intValue *Integer
		want     int
	}{
		{
			name:     "0",
			intValue: &Integer{Value: "0"},
			want:     0,
		},
		{
			name:     "99",
			intValue: &Integer{Value: "99"},
			want:     99,
		},
		{
			// nil Integer pointer results in the default def value returned
			name:     "nil-integer",
			intValue: nil,
			want:     0,
		},
		{
			// non-convertable string results in the default def value returned
			name:     "non-convertable",
			intValue: &Integer{Value: "non-convertable"},
			want:     0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := Span{Right: c.intValue}
			def := 0
			have := s.GetRight(def)
			if have != c.want {
				t.Errorf("have: %d; want: %d", have, c.want)
			}
		})
	}
}

// Test if the quoted string values get their quotes trimmed.
func TestStringTrim(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{"'students'", "students"},
		{"\"employees\"", "employees"},
		{"", ""},
	}
	for _, c := range cases {
		t.Run(c.value, func(t *testing.T) {
			s := String{c.value}
			have := s.Trim()
			if have != c.want {
				t.Errorf("have: %s; want: %s", have, c.want)
			}
		})
	}
}

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
