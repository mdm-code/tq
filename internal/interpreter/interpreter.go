package interpreter

import (
	"github.com/mdm-code/tq/v2/internal/ast"
)

// FilterFunc specifies the data transformation function type.
type FilterFunc func(data ...any) ([]any, error)

type filter struct {
	name     string
	inner    FilterFunc
	optional bool
}

func (f *filter) call(data ...any) ([]any, error) {
	return f.inner(data...)
}

// Interpreter interprets the tq query AST into a pipe-like sequence of
// filtering functions processing TOML input data as specified in the query.
type Interpreter struct {
	filters []filter
}

// New returns a new instance of Interpreter.
func New() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) eval(exprs ...ast.Expr) {
	for _, e := range exprs {
		e.Accept(i)
	}
}

// Interpret extracts a sequence of filtering functions by traversing the AST.
// It returns an entry function that takes in deserialized TOML data and
// applies filtering functions in the sequence provided by the Interpreter.
func (i *Interpreter) Interpret(expr ast.Expr) FilterFunc {
	i.filters = nil // clear out previously accumulated filtering functions
	i.eval(expr)
	return func(data ...any) ([]any, error) {
		var err error
		for _, f := range i.filters {
			data, err = f.call(data...)
			if err != nil && !f.optional {
				return data, err
			}
		}
		return data, nil
	}
}

// VisitRoot interprets the Root AST node.
func (i *Interpreter) VisitRoot(root *ast.Root) {
	i.eval(root.Query)
}

// VisitQuery interprets the Query AST node.
func (i *Interpreter) VisitQuery(query *ast.Query) {
	i.eval(query.Filters...)
}

// VisitFilter interprets the Filter AST node.
func (i *Interpreter) VisitFilter(filter *ast.Filter) {
	i.eval(filter.Kind)
	i.filters[len(i.filters)-1].optional = filter.Optional
}

// VisitIdentity interprets the Identity AST node.
func (i *Interpreter) VisitIdentity(identity *ast.Identity) {
	f := filter{
		name: "identity",
		inner: func(data ...any) ([]any, error) {
			return data, nil
		},
	}
	i.filters = append(i.filters, f)
}

// VisitSelector interprets the Selector AST node.
func (i *Interpreter) VisitSelector(selector *ast.Selector) {
	i.eval(selector.Value)
}

// VisitSpan interprets the Span AST node.
func (i *Interpreter) VisitSpan(span *ast.Span) {
	f := filter{
		name: "span",
		inner: func(data ...any) ([]any, error) {
			result := make([]any, 0, len(data))
			var err error
			for _, d := range data {
				switch v := d.(type) {
				case []any:
					l, r := span.GetLeft(0), span.GetRight(len(v))
					if r > len(v) {
						r = len(v)
					}
					if l > r || l >= len(v) {
						continue
					}
					result = append(result, v[l:r])
				default:
					err = &Error{
						data:   d,
						filter: span.String(),
						err:    ErrTOMLDataType,
					}
				}
			}
			return result, err
		},
	}
	i.filters = append(i.filters, f)
}

// VisitIterator interprets the Iterator AST node.
func (i *Interpreter) VisitIterator(iterator *ast.Iterator) {
	f := filter{
		name: "iterator",
		inner: func(data ...any) ([]any, error) {
			result := make([]any, 0, len(data))
			var err error
			for _, d := range data {
				switch v := d.(type) {
				case map[string]any:
					for _, val := range v {
						result = append(result, val)
					}
				case []any:
					result = append(result, v...)
				default:
					err = &Error{
						data:   d,
						filter: iterator.String(),
						err:    ErrTOMLDataType,
					}
				}
			}
			return result, err
		},
	}
	i.filters = append(i.filters, f)
}

// VisitString interprets the String AST node.
func (i *Interpreter) VisitString(str *ast.String) {
	f := filter{
		name: "string",
		inner: func(data ...any) ([]any, error) {
			result := make([]any, 0, len(data))
			var err error
			for _, d := range data {
				switch v := d.(type) {
				case map[string]any:
					key := str.Value
					res, ok := v[key]
					if ok {
						result = append(result, res)
					}
				default:
					err = &Error{
						data:   d,
						filter: str.String(),
						err:    ErrTOMLDataType,
					}
				}
			}
			return result, err
		},
	}
	i.filters = append(i.filters, f)
}

// VisitInteger interprets the Integer AST node.
func (i *Interpreter) VisitInteger(integer *ast.Integer) {
	f := filter{
		name: "integer",
		inner: func(data ...any) ([]any, error) {
			result := make([]any, 0, len(data))
			var err error
			for _, d := range data {
				switch v := d.(type) {
				case []any:
					idx, _ := integer.Vtoi()
					if idx >= 0 && idx < len(v) {
						result = append(result, v[idx])
					}
				default:
					err = &Error{
						data:   d,
						filter: integer.String(),
						err:    ErrTOMLDataType,
					}
				}
			}
			return result, err
		},
	}
	i.filters = append(i.filters, f)
}
