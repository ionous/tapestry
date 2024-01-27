package decode

import (
	"errors"
	"fmt"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

func decodeString(out r.Value, spec *typeinfo.Str, val any) (err error) {
	switch k := out.Kind(); k {
	default:
		panic("unexpected enum target")
	case r.Bool:
		if a, ok := val.(bool); !ok {
			err = conversionError(k, val)
		} else {
			out.SetBool(a)
		}
	case r.Int:
		if a, ok := val.(string); !ok {
			err = conversionError(k, val)
		} else if i := spec.FindOption(a); i < 0 {
			err = fmt.Errorf("invalid option %q for %s", val, spec.Name)
		} else {
			out.SetInt(int64(i))
		}
	case r.String:
		// with no constraints, the generated code is a string.
		switch val := val.(type) {
		default:
			err = conversionError(r.String, val)
		case string:
			out.SetString(val)
		case []any:
			// join multiple lines
			// tbd: remove, and just use multiline tell strings;
			// although possibly they'd be carried to here the same way... so.
			if a, ok := compact.JoinLines(val); !ok {
				err = errors.New("expected lines of plain value strings")
			} else {
				out.SetString(a)
			}
		}
	}
	return
}

// expects val is a string or a plain value list of strings
func decodeStrings(out r.Value, spec *typeinfo.Str, val any) (err error) {
	// two cases, the generated code might be a bool or an int enum
	switch k := out.Type().Elem().Kind(); k {
	default:
		panic("unexpected enum target")
	case r.Bool:
		if vals, ok := normalizeBools(val); !ok {
			err = errors.New("expected a bool or slice")
		} else {
			out.Set(r.ValueOf(vals))
		}

	case r.Int:
		if strs, ok := normalizeStrings(val); !ok {
			err = errors.New("expected a string or slice")
		} else if at := verifyStrings(spec, strs); at < 0 {
			out.Set(r.ValueOf(strs))
		} else {
			err = fmt.Errorf("invalid option %q for %s", strs[at], spec.Name)
		}

	case r.String:
		if strs, ok := normalizeStrings(val); !ok {
			err = errors.New("expected a string or slice")
		} else {
			out.Set(r.ValueOf(strs))
		}
	}
	return
}

// expects val is a float or a plain value list of floats
// tbd: this assumes floats because that's what json uses;
// tell can optionally use ints (etc); expand this?
func decodeNumbers(out r.Value, val any) (err error) {
	if vals, ok := normalizeFloats(val); !ok {
		err = errors.New("expected a number or slice")
	} else {
		out.Set(r.ValueOf(vals))
	}
	return
}

func normalizeBools(val any) (ret []bool, okay bool) {
	switch v := val.(type) {
	case bool:
		ret, okay = []bool{v}, true
	case []any:
		ret, okay = compact.SliceBools(v)
	}
	return
}

func normalizeFloats(val any) (ret []float64, okay bool) {
	switch v := val.(type) {
	case float64:
		ret, okay = []float64{v}, true
	case []any:
		ret, okay = compact.SliceFloats(v)
	}
	return
}

func normalizeStrings(val any) (ret []string, okay bool) {
	switch v := val.(type) {
	case string:
		ret, okay = []string{v}, true
	case []any:
		ret, okay = compact.SliceStrings(v)
	}
	return
}

// index of first invalid string
func verifyStrings(spec *typeinfo.Str, strs []string) (ret int) {
	ret = -1 // provisionally
	for at, str := range strs {
		if spec.FindOption(str) < 0 {
			ret = at
			break
		}
	}
	return
}

func conversionError(k r.Kind, have any) (err error) {
	return fmt.Errorf("can't set a %s using a %T", k, have)
}
