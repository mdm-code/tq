package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/internal/ast"
	"github.com/mdm-code/tq/internal/lexer"
)

// TODO: Check all the errors from the cover profile that are not covered
// by the current test case.

// Check if the AST returned by Parse() method matches its predicted output.
func TestParse(t *testing.T) {
	cases := []struct {
		query string
		want  ast.Expr
	}{
		{
			query: ".['students'][2:4][0]['grades'][]",
			want: &ast.Root{
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
										Value: "2",
									},
									Right: &ast.Integer{
										Value: "4",
									},
								},
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
									Value: "'grades'",
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
	}
	for _, c := range cases {
		t.Run(c.query, func(t *testing.T) {
			r := strings.NewReader(c.query)
			s, _ := scanner.New(r)
			l, _ := lexer.New(s)
			p, err := New(l)
			if err != nil {
				t.Fatal(err)
			}
			have, err := p.Parse()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(have, c.want) {
				t.Errorf("have: %v; want: %v", have, c.want)
			}
		})
	}
}
