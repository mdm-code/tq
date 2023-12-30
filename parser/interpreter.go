package parser

// FilterFn specifies the data transformation function type.
type FilterFn func(data ...interface{}) ([]interface{}, error)

// Interpreter interprets the tq query AST into a pipe-like sequence of
// filtering functions processing TOML input data as specified in the query.
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
// applies filtering functions in the sequence provided by the Interpreter.
func (i *Interpreter) Interpret(root Expr) FilterFn {
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
	identityFn := func(data ...interface{}) ([]interface{}, error) {
		return data, nil
	}
	i.filters = append(i.filters, identityFn)
}

func (i *Interpreter) visitSelector(e Expr) {
	s := e.(*Selector)
	i.eval(s.value)
}

func (i *Interpreter) visitSpan(e Expr) {
	s := e.(*Span)
	spanFn := func(data ...interface{}) ([]interface{}, error) {
		result := make([]interface{}, 0, len(data))
		var err error
		for _, d := range data {
			switch v := d.(type) {
			case []interface{}:
				l, r := s.Left(0), s.Right(len(v))
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

func (i *Interpreter) visitIterator(e Expr) {
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

func (i *Interpreter) visitString(e Expr) {
	str := e.(*String)
	strFn := func(data ...interface{}) ([]interface{}, error) {
		result := make([]interface{}, 0, len(data))
		var err error
		for _, d := range data {
			switch v := d.(type) {
			case map[string]interface{}:
				key := str.trimmed()
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

func (i *Interpreter) visitInteger(e Expr) {
	integer := e.(*Integer)
	intFn := func(data ...interface{}) ([]interface{}, error) {
		result := make([]interface{}, 0, len(data))
		var err error
		for _, d := range data {
			switch v := d.(type) {
			case []interface{}:
				idx, _ := integer.vtoi()
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
