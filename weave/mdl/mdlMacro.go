package mdl

import "strings"

// the reversed relation name
func fmtMacro(name string, reversed bool) (ret string) {
	ps := []string{"the"}
	if reversed {
		ps = append(ps, "reversed")
	}
	ps = append(ps, "macro", name)
	return strings.Join(ps, " ")
}
