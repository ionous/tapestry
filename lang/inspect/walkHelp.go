package inspect

import (
	r "reflect"
)

// assuming that the passed value implements Stringer, return its string.
func ReflectStringer(v r.Value) string {
	// is there a better way?
	return v.Interface().(interface{ String() string }).String()
}
