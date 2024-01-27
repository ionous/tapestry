package decode

import (
	"fmt"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/lang/inspect"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

func nextField(it *inspect.It, param string) (ret typeinfo.Term, okay bool) {
	for it.Next() {
		info := it.Term() // internal fields dont have labels....
		if !info.Private && info.Label == param {
			ret, okay = info, true
			break
		}
	}
	return
}

func unknownType(t r.Type) error {
	return fmt.Errorf("unknown type %s(%s)", t.Kind(), t.String())
}
