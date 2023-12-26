package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// FilterFn specifies the data transformation function type.
type FilterFn func(data ...interface{}) ([]interface{}, error)

// Interpreter ...
type Interpreter struct {
	filters []FilterFn
}

func (i *Interpreter) eval(es ...Expr) {
	for _, e := range es {
		e.accept(i)
	}
}

// Interpret extracts a sequence of filtering functions by traversing the AST.
// It returns an entry function that takes in deserialized TOML data and
// applies filtering functions in the sequence established by the Interpreter.
func (i *Interpreter) Interpret(root Expr) FilterFn {
	i.filters = nil // clear out previously accumulated filtering functions
	i.eval(root)
	return func(data ...interface{}) ([]interface{}, error) {
		for _, fn := range i.filters {
			data, err := fn(data...)
			if err != nil {
				return data, err
			}
		}
		return data, nil
	}
}

func (i *Interpreter) visitRoot(e Expr) {
	r := e.(*Root)
	i.eval(r.query)
}

func (i *Interpreter) visitQuery(e Expr) {
	q := e.(*Query)
	i.eval(q.filters...)
}

func (i *Interpreter) visitFilter(e Expr) {
	f := e.(*Filter)
	i.eval(f.kind)
}
func (i *Interpreter) visitIdentity(e Expr) {
	i.filters = append(i.filters, func(data ...interface{}) ([]interface{}, error) {
		return data, nil
	})
}

// TODO: this receiver function has to be undergo a major rework. This includes
// things like (1) simplified type evaluation (see methods above), (2) separate
// methods for different selector values 1. Span, 2. String, 3. Integer, 4.
// Iterator. The Span should allow to infer left/right just fine with i.eval
// the same way.
func (i *Interpreter) visitSelector(e Expr) {
	switch v := e.(type) {
	case *Selector:
		fn := func(data ...interface{}) ([]interface{}, error) {
			var err error
			result := []interface{}{}
			switch vv := v.value.(type) {
			case *String:
				for _, d := range data {
					switch vvv := d.(type) {
					case map[string]interface{}:
						val := vv.value
						val = strings.Trim(val, "'") // might want trim bytes instead
						val = strings.Trim(val, "\"")
						result = append(result, vvv[val])
					default:
						err = fmt.Errorf("type error")
					}
				}
			case *Integer:
				for _, d := range data {
					switch vvv := d.(type) {
					case []interface{}:
						i, _ := strconv.Atoi(vv.value)
						result = append(result, vvv[i])
					default:
						err = fmt.Errorf("type error")
					}
				}
			case *Span:
				var l int
				if vv.left != nil {
					l, _ = strconv.Atoi(vv.left.value)
				} else {
					l = 0
				}
				for _, d := range data {
					switch vvv := d.(type) {
					case []interface{}:
						var r int
						if vv.right != nil {
							r, _ = strconv.Atoi(vv.right.value)
							if r > len(vvv) {
								r = len(vvv)
							}
						} else {
							r = len(vvv)
						}
						result = append(result, vvv[l:r])
					default:
						err = fmt.Errorf("type error")
					}
				}
			case *Iterator:
				for _, d := range data {
					switch v := d.(type) {
					case []interface{}:
						for _, v := range v {
							result = append(result, v)
						}
					case map[string]interface{}:
						for _, v := range v {
							result = append(result, v)
						}
					default:
						err = fmt.Errorf("type error")
					}
				}
			}
			return result, err
		}
		i.filters = append(i.filters, fn)
	default:
		// error out
	}
}
