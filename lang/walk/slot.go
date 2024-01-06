package walk

import (
	r "reflect"
	"strings"
	"unicode"
)

// transform PascalCase to under_score
// maybe store this in the slot registry instead
// *or* add it t the the if labels slot=...
// ( which would be redundant but useful )
func SlotName(slot r.Type) string {
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
