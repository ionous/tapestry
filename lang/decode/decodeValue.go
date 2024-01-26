package decode

import (
	"errors"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/lang/compact"
)

// expects val is a string or a plain value list of strings
func decodeStrings(out r.Value, val any) (err error) {
	switch v := val.(type) {
	case string:
		out.Set(r.ValueOf([]string{v}))
	case []any:
		if vs, ok := compact.SliceStrings(v); !ok {
			err = errors.New("expected a slice of strings")
		} else {
			out.Set(r.ValueOf(vs))
		}
	default:
		err = errors.New("expected a string or slice")
	}
	return
}

// expects val is a float or a plain value list of floats
// tbd: this assumes floats because that's what json uses;
// tell can optionally use ints (etc); expand this?
func decodeNumbers(out r.Value, val any) (err error) {
	switch v := val.(type) {
	case float64:
		out.Set(r.ValueOf([]float64{v}))
	case []any:
		if vs, ok := compact.SliceFloats(v); !ok {
			err = errors.New("expected a slice of floats")
		} else {
			out.Set(r.ValueOf(vs))
		}
	default:
		err = errors.New("expected a float or slice")
	}
	return
}
