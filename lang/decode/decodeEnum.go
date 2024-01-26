package decode

import (
	"errors"
	"fmt"
	r "reflect"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// fix? it's a bit of a cheat that there's no "SetCompactValue"
// we peek at the value to see if its a key and switch on processing
// we currently store the $KEY in the in-memory values.
func SetString(dst r.Value, spec *typeinfo.Str, kv any) (err error) {
	if src, ok := toString(kv); !ok {
		err = errors.New("not string data")
	} else {
		// fix: when everything is conveted should be using spec.Options to switch
		if dst.Kind() == r.Int {
			if i := spec.FindOption(src); i < 0 {
				err = fmt.Errorf("invalid option %q for %s", src, spec.Name)
			} else {
				dst.SetInt(int64(i))
			}

		} else {
			if str, ok := xformString(src, spec); !ok {
				// hrm...
				err = errors.New("invalid string")
			} else {
				// FIX
				if dst.Kind() == r.String {
					dst.Set(r.ValueOf(str))
				} else {
					dst.Field(0).Set(r.ValueOf(str))
				}
			}
		}
	}
	return
}

// fix: limit special handling for prim.lines?
func toString(v any) (ret string, okay bool) {
	switch s := v.(type) {
	case string:
		ret, okay = s, true
	case []any:
		ret, okay = compact.SliceLines(s)
	}
	return
}

func xformString(str string, spec *typeinfo.Str) (ret string, okay bool) {
	if opt := spec.Options; len(opt) == 0 {
		ret, okay = str, true
	} else if i := spec.FindOption(str); i >= 0 {
		ret, okay = "$"+strings.ToUpper(str), true
	}
	return
}
