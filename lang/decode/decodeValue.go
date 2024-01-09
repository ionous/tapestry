package decode

import (
	"errors"
	"fmt"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/lang/compact"
)

// write a value into the target of an iterator.
func SetValue(out r.Value, val any) (err error) {
	switch k := out.Kind(); k {
	default:
		err = fmt.Errorf("invalid kind of %q", k)

	case r.Bool:
		if v, ok := val.(bool); !ok {
			panic("---- are these stored as strings???")
		} else {
			out.Set(r.ValueOf(v))
		}

	case r.Float64:
		if v, ok := val.(float64); !ok {
			err = errors.New("expected a float")
		} else {
			out.Set(r.ValueOf(v))
		}

	case r.String:
		if v, ok := val.(string); !ok {
			err = errors.New("expected a string")
			return
		} else {
			out.Set(r.ValueOf(v))
		}
	}
	return
}

// write a value into the target of an iterator.
func SetValues(out r.Value, val any) (err error) {
	t := out.Type()
	switch k := t.Kind(); k {
	default:
		panic("unknown type of value")

	case r.Slice:
		switch el := t.Elem().Kind(); el {
		default:
			panic("unknown type of slice")
		case r.Float64:
			switch v := val.(type) {
			case float64:
				out.Set(r.ValueOf([]float64{v}))
			case []any:
				if vs, ok := compact.SliceFloats(v); !ok {
					err = errors.New("couldnt convert to a slice of floats")
				} else {
					out.Set(r.ValueOf(vs))
				}
			default:
				err = errors.New("expected a float or float slice")
			}
		case r.String:
			switch v := val.(type) {
			case string:
				out.Set(r.ValueOf([]string{v}))
			case []any:
				if vs, ok := compact.SliceStrings(v); !ok {
					err = errors.New("couldnt convert to a slice of floats")
				} else {
					out.Set(r.ValueOf(vs))
				}
			default:
				err = errors.New("expected a string or string slice")
			}
		}
	}
	return
}
