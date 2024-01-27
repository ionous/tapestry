package encode

import (
	"git.sr.ht/~ionous/tapestry/lang/inspect"
)

func encodeValues(it inspect.It) []any {
	var out []any
	if cnt := it.Len(); cnt > 0 {
		out = make([]any, 0, cnt)
		for it.Next() {
			val := it.CompactValue()
			out = append(out, val)
		}
	}
	return out
}
