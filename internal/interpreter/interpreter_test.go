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
