package decode

import (
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
