package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/internal/ast"
	"github.com/mdm-code/tq/internal/lexer"
)

// Test the failing Parser New() constructor.
func TestNewFails(t *testing.T) {
	q := ".[\"instances\"][$]" // the unsupported $ character causes the error.
	r := strings.NewReader(q)
	s, _ := scanner.New(r)
	l, _ := lexer.New(s)
	_, err := New(l)
	if err == nil {
		t.Fatal("expected the constructor to fail")
	}
}

// Verify the range of possible Parser errors returned by Parse() method.
func TestParseErrored(t *testing.T) {
	// TODO: implement the cases
	cases := []struct {
		query string
		err   error
	}{}
	for _, c := range cases {
		t.Run(c.query, func(t *testing.T) {
		})
	}
}

// Check if the AST returned by Parse() method matches its predicted output.
func TestParse(t *testing.T) {
	cases := []struct {
		query string
		want  ast.Expr
	}{
		{
			query: ".['students'][2:4][0]['grades'][:6][]",
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
								Value: &ast.Span{
									Left: nil,
									Right: &ast.Integer{
										Value: "6",
									},
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
