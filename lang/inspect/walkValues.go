package inspect

import (
	r "reflect"
)

// turn a Number value into a float; panics otherwise.
func (w *It) Float() float64 {
	v, _ := w.getFocus()
	return v.Float()
}

// turn a Str value into a string; panics otherwise.
// as a special case, knows to transform boolean values to "true" or "false"
func (w *It) String() (ret string) {
	v, _ := w.getFocus()
	switch v.Kind() {
	case r.Bool:
		if v.Bool() {
			ret = "true"
		} else {
			ret = "false"
		}
	case r.String:
		ret = v.String()
	case r.Int:
		ret = v.Interface().(interface{ String() string }).String()
	}
	return
}
