package parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// FilterFunc ...
type FilterFunc func(data ...interface{}) ([]interface{}, error)

// QueryConstructor ...
type QueryConstructor struct {
	Filters []FilterFunc
}

func (q *QueryConstructor) interpret(es ...Expr) {
	for _, e := range es {
		e.accept(q)
	}
}

// Run ...
func (q *QueryConstructor) Run(e Expr) {
	e.accept(q)
}

func (q *QueryConstructor) visitRoot(e Expr) {
	switch v := e.(type) {
	case *Root:
		q.interpret(v.query)
	default:
		// error out
	}
}

func (q *QueryConstructor) visitQuery(e Expr) {
	switch v := e.(type) {
	case *Query:
		q.interpret(v.filters...)
	default:
		// error out
	}
}

func (q *QueryConstructor) visitFilter(e Expr) {
	switch v := e.(type) {
	case *Filter:
		q.interpret(v.kind)
	default:
		// error out
	}
}
func (q *QueryConstructor) visitIdentity(e Expr) {
	switch v := e.(type) {
	case *Identity:
		fmt.Fprintf(io.Discard, "%v", *v)
		q.Filters = append(q.Filters, identityFn)
	default:
		// error out
	}
}

func identityFn(data ...interface{}) ([]interface{}, error) {
	return data, nil
}

func (q *QueryConstructor) visitSelector(e Expr) {
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
		q.Filters = append(q.Filters, fn)
	default:
		// error out
	}
}
