package interpreter

import (
	"testing"

	"github.com/mdm-code/tq/internal/ast"
)

// Test the public API of the Interpreter.
func TestInterpret(t *testing.T) {
	cases := []struct {
		query string
		want  []filter
		root  ast.Expr
		data  any
	}{
		{
			data: map[string]any{
				"students": []any{
					map[string]any{
						"grades": []any{
							map[string]any{
								"first": 2,
							},
						},
					},
					map[string]any{
						"grades": []any{
							map[string]any{
								"first": []any{2},
							},
						},
					},
				},
			},
			query: ".['students'][0:99][1][][0]['first'][]",
			want: []filter{
				{name: "identity"},
				{name: "string"},
				{name: "span"},
				{name: "integer"},
				{name: "iterator"},
				{name: "integer"},
				{name: "string"},
				{name: "iterator"},
			},
			root: &ast.Root{
				Query: &ast.Query{
					Filters: []ast.Expr{
						&ast.Filter{
							Kind: &ast.Identity{},
						},
						&ast.Filter{
							Kind: &ast.Selector{
								Value: &ast.String{
									Value: "'students'",
								},
							},
						},
						&ast.Filter{
							Kind: &ast.Selector{
								Value: &ast.Span{
									Left: &ast.Integer{
										Value: "0",
									},
									Right: &ast.Integer{
										Value: "99",
									},
								},
							},
						},
						&ast.Filter{
							Kind: &ast.Selector{
								Value: &ast.Integer{
									Value: "1",
								},
							},
						},
						&ast.Filter{
							Kind: &ast.Selector{
								Value: &ast.Iterator{},
							},
						},
						&ast.Filter{
							Kind: &ast.Selector{
								Value: &ast.Integer{
									Value: "0",
								},
							},
						},
						&ast.Filter{
							Kind: &ast.Selector{
								Value: &ast.String{
									Value: "'first'",
								},
							},
						},
						&ast.Filter{
							Kind: &ast.Selector{
								Value: &ast.Iterator{},
							},
						},
					},
				},
			},
		},
		{
			data: map[string]any{
				"salaries": []any{
					"1000", "2_000",
				},
			},
			query: "",
			want:  []filter{},
			root: &ast.Root{
				Query: &ast.Query{},
			},
		},
		{
			data: map[string]any{
				"employees": []any{
					"Bob", "Mike",
				},
			},
			query: ".",
			want:  []filter{{name: "identity"}},
			root: &ast.Root{
				Query: &ast.Query{
					Filters: []ast.Expr{
						&ast.Filter{
							Kind: &ast.Identity{},
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.query, func(t *testing.T) {
			i := New()
			exec := i.Interpret(c.root)
			have := i.filters
			for i, f := range c.want {
				if f.name != have[i].name {
					t.Errorf("have: %s; want: %s", f.name, have[i].name)
				}
			}
			_, err := exec(c.data)
			if err != nil {
				t.Error("failed to filter data with declared filters")
			}
		})
	}
}

// Check if filter function errors out when provided unsupported data input.
func TestVisitError(t *testing.T) {
	var data interface{}
	i := New()
	cases := []struct {
		name string
		node ast.Expr
		fn   func(ast.Expr)
	}{
		{
			"integerNode",
			&ast.Integer{},
			(*i).VisitInteger,
		},
		{
			"stringNode",
			&ast.String{},
			(*i).VisitString,
		},
		{
			"iteratorNode",
			&ast.Iterator{},
			(*i).VisitIterator,
		},
		{
			"spanNode",
			&ast.Span{},
			(*i).VisitSpan,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() { i.filters = nil }()
			c.fn(c.node)
			filter := i.filters[0]
			_, err := filter.inner(data)
			if err == nil {
				t.Errorf("filter function should error with data: %v", data)
			}
		})
	}
}

// Verify if Interpret fails when provided unsupported data.
func TestInterpretError(t *testing.T) {
	var data interface{}
	root := &ast.Root{
		Query: &ast.Query{
			Filters: []ast.Expr{
				&ast.Filter{
					Kind: &ast.Span{},
				},
			},
		},
	}
	i := Interpreter{}
	exec := i.Interpret(root)
	_, err := exec(data)
	if err == nil {
		t.Errorf("Interpret should fail with data: %v", data)
	}

}
