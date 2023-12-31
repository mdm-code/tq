package interpreter

import (
	"github.com/mdm-code/tq/internal/ast"
)

// FilterFn specifies the data transformation function type.
type FilterFn func(data ...interface{}) ([]interface{}, error)

// Interpreter interprets the tq query AST into a pipe-like sequence of
// filtering functions processing TOML input data as specified in the query.
type Interpreter struct {
	filters []FilterFn
}

func (i *Interpreter) eval(es ...ast.Expr) {
	for _, e := range es {
		e.Accept(i)
	}
}

// Interpret extracts a sequence of filtering functions by traversing the AST.
// It returns an entry function that takes in deserialized TOML data and
// applies filtering functions in the sequence provided by the Interpreter.
func (i *Interpreter) Interpret(root ast.Expr) FilterFn {
	i.filters = nil // clear out previously accumulated filtering functions
	i.eval(root)
	return func(data ...interface{}) ([]interface{}, error) {
		var err error
		for _, fn := range i.filters {
			data, err = fn(data...)
			if err != nil {
				return data, err
			}
		}
		return data, nil
	}
}

// VisitRoot interprets the Root AST node.
func (i *Interpreter) VisitRoot(e ast.Expr) {
	r := e.(*ast.Root)
	i.eval(r.Query)
}

// VisitQuery interprets the Query AST node.
func (i *Interpreter) VisitQuery(e ast.Expr) {
	q := e.(*ast.Query)
	i.eval(q.Filters...)
}

// VisitFilter interprets the Filter AST node.
func (i *Interpreter) VisitFilter(e ast.Expr) {
	f := e.(*ast.Filter)
	i.eval(f.Kind)
}

// VisitIdentity interprets the Identity AST node.
func (i *Interpreter) VisitIdentity(e ast.Expr) {
	identityFn := func(data ...interface{}) ([]interface{}, error) {
		return data, nil
	}
	i.filters = append(i.filters, identityFn)
}

// VisitSelector interprets the Selector AST node.
func (i *Interpreter) VisitSelector(e ast.Expr) {
	s := e.(*ast.Selector)
	i.eval(s.Value)
}

// VisitSpan interprets the Span AST node.
func (i *Interpreter) VisitSpan(e ast.Expr) {
	s := e.(*ast.Span)
	spanFn := func(data ...interface{}) ([]interface{}, error) {
		result := make([]interface{}, 0, len(data))
		var err error
		for _, d := range data {
			switch v := d.(type) {
			case []interface{}:
				l, r := s.GetLeft(0), s.GetRight(len(v))
				if r > len(v) {
					r = len(v)
				}
				if l > r || l >= len(v) {
					continue
				}
				result = append(result, v[l:r])
			default:
				err = ErrTOMLDataType
			}
		}
		return result, err
	}
	i.filters = append(i.filters, spanFn)
}

// VisitIterator interprets the Iterator AST node.
func (i *Interpreter) VisitIterator(e ast.Expr) {
	iterFn := func(data ...interface{}) ([]interface{}, error) {
		result := make([]interface{}, 0, len(data))
		var err error
		for _, d := range data {
			switch v := d.(type) {
			case map[string]interface{}:
				for _, val := range v {
					result = append(result, val)
				}
			case []interface{}:
				for _, val := range v {
					result = append(result, val)
				}
			default:
				err = ErrTOMLDataType
			}
		}
		return result, err
	}
	i.filters = append(i.filters, iterFn)
}

// VisitString interprets the String AST node.
func (i *Interpreter) VisitString(e ast.Expr) {
	str := e.(*ast.String)
	strFn := func(data ...interface{}) ([]interface{}, error) {
		result := make([]interface{}, 0, len(data))
		var err error
		for _, d := range data {
			switch v := d.(type) {
			case map[string]interface{}:
				key := str.Trim()
				res, ok := v[key]
				if ok {
					result = append(result, res)
				}
			default:
				err = ErrTOMLDataType
			}
		}
		return result, err
	}
	i.filters = append(i.filters, strFn)
}

// VisitInteger interprets the Integer AST node.
func (i *Interpreter) VisitInteger(e ast.Expr) {
	integer := e.(*ast.Integer)
	intFn := func(data ...interface{}) ([]interface{}, error) {
		result := make([]interface{}, 0, len(data))
		var err error
		for _, d := range data {
			switch v := d.(type) {
			case []interface{}:
				idx, _ := integer.Vtoi()
				if idx >= 0 && idx < len(v) {
					result = append(result, v[idx])
				}
			default:
				err = ErrTOMLDataType
			}
		}
		return result, err
	}
	i.filters = append(i.filters, intFn)
}
