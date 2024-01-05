package decode

import (
	"errors"
	"fmt"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/dl/composer"
)

// fix? it's a bit of a cheat that there's no "SetCompactValue"
// we peek at the value to see if its a key and switch on processing
// we currently store the $KEY in the in-memory values.
func SetString(dst r.Value, strtype r.Type, kv any) (err error) {
	c := r.New(strtype).Interface() // ugh. fix.
	if n, ok := c.(composer.Composer); !ok {
		err = fmt.Errorf("is %s not a generated type?", dst.Type())
	} else if src, ok := kv.(string); !ok {
		err = errors.New("not string data")
	} else if str, ok := xformString(src, n.Compose()); !ok {
		err = errors.New("invalid string ")
	} else {
		dst.Set(r.ValueOf(str))
	}
	return
}

func xformString(str string, spec composer.Spec) (ret string, okay bool) {
	if len(str) > 0 && str[0] == '$' {
		// checks open strings in case we're using a $ in a keyd value
		// which we should probably handle in a more clever way
		// fix: i think OpenStrings should probably never use keys
		// and the values should just be hints ....
		if _, i := spec.IndexOfChoice(str); i >= 0 || spec.OpenStrings {
			ret, okay = str, true
		}
	} else {
		if k, i := spec.IndexOfValue(str); i >= 0 {
			ret, okay = k, true
		} else if spec.OpenStrings {
			ret, okay = str, true
		}
	}
	return
}
