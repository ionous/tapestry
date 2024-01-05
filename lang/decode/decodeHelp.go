package decode

import (
	r "reflect"
	"strings"
	"unicode"

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

// transform PascalCase to under_score
// maybe store this in the slot registry instead
// *or* add it t the the if labels slot=...
// ( which would be redundant but useful )
func slotName(slot r.Type) string {
	var out strings.Builder
	var prev bool
	str := slot.Name()
	for _, r := range str {
		l := unicode.ToLower(r)
		cap := l != r
		if !prev && cap && out.Len() > 0 {
			out.WriteRune('_')
		}
		out.WriteRune(l)
		prev = cap
	}
	return out.String()
}
