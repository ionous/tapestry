package decode

import (
	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/walk"
)

func parseMessage(v any) (ret compact.Message, err error) {
	if m, ok := v.(map[string]any); !ok {
		err = ValueError("not a key value map", v)
	} else {
		ret, err = DecodeMessage(m)
	}
	return
}

func findLast(it walk.Walker) walk.Walker {
	last := it
	for it.Next() {
		last = it
	}
	return last
}

func nextField(it *walk.Walker, p compact.Param) (ret walk.Field, okay bool) {
	for it.Next() {
		info := it.Field() // internal fields dont have labels....
		if label, ok := info.Label(); ok {
			if p.Matches(label) {
				ret, okay = info, true
				break
			}
		}
	}
	return
}
