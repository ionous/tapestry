package decode

import (
	"errors"
	"fmt"
	r "reflect"
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
	switch k := out.Kind(); k {
	default:
		panic("unknown type of value")
	case r.Slice:
		switch el := out.Elem().Kind(); el {
		default:
			panic("unknown type of slice")
		case r.Float64:
			if v, ok := val.([]float64); !ok {
				err = errors.New("expected a float slice")
			} else {
				out.Set(r.ValueOf(v))
			}
		case r.String:
			if v, ok := val.([]string); !ok {
				err = errors.New("expected a string slice")
			} else {
				out.Set(r.ValueOf(v))
			}
		}
	}
	return
}
