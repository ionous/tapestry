package decode

import (
	"fmt"
	r "reflect"
	"strings"
	"unicode"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/walk"
)

// transform PascalCase to under_score
// maybe store this in the slot registry instead
// *or* add it t the the if labels slot=...
// ( which would be redundant but useful )
func nameOf(t r.Type) string {
	var out strings.Builder
	var prev bool
	str := t.Name()
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

func parseMessage(v any) (ret compact.Message, err error) {
	if m, ok := v.(map[string]any); !ok {
		err = ValueError("not a key value map", v)
	} else {
		ret, err = DecodeMessage(m)
	}
	return
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

func unknownType(t r.Type) error {
	return fmt.Errorf("unknown type %s(%s)", t.Kind(), t.String())
}
