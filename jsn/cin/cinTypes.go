package cin

import (
	r "reflect"
	"strings"

	"github.com/ionous/errutil"
)

func IsValidMap(t r.Type) bool {
	return t.Kind() == r.Map &&
		t.Elem().Kind() == r.Interface &&
		t.Key().Kind() == r.String
}

func IsValidSlice(t r.Type) bool {
	return t.Kind() == r.Slice &&
		t.Elem().Kind() == r.Interface
}

func SliceFloats(slice r.Value) (ret []float64, okay bool) {
	okay = true // provisionally
	for i, cnt := 0, slice.Len(); i < cnt; i++ {
		if el := slice.Index(i).Elem(); el.Kind() != r.Float64 {
			okay = false
			break
		} else {
			if ret == nil {
				ret = make([]float64, cnt)
			}
			ret[i] = el.Float()
		}
	}
	return
}

func SliceStrings(slice r.Value) (ret []string, okay bool) {
	okay = true // provisionally
	for i, cnt := 0, slice.Len(); i < cnt; i++ {
		if el := slice.Index(i).Elem(); el.Kind() != r.String {
			okay = false
			break
		} else {
			if ret == nil {
				ret = make([]string, cnt)
			}
			ret[i] = el.String()
		}
	}
	return
}

// attempt to read a r.Value as a slice strings:
// returns those strings joined with newlines.
func SliceLines(slice r.Value) (ret string, err error) {
	if t := slice.Type(); !IsValidSlice(t) {
		err = errutil.Fmt("expected a slice of interface, got %s", t)
	} else {
		var b strings.Builder
		for i, cnt := 0, slice.Len(); i < cnt; i++ {
			if el := slice.Index(i).Elem(); el.Kind() != r.String {
				err = errutil.New("expected a string, got %s", el.Type())
				break
			} else {
				str := el.String()
				if i > 0 {
					b.WriteRune('\n')
				}
				b.WriteString(str)
			}
		}
		if err == nil {
			ret = b.String()
		}
	}
	return
}
