package rules

import (
	"strings"
)

// words that can follow a rule name
// ex. "taking then continue"
type Suffix int

//go:generate stringer -type=Suffix -linecomment
const (
	UnspecfiedSuffix Suffix = iota
	Continues               // then continue
	Stops                   // then stop
	Skips                   // then skip phase
	//
	NumSuffixes = iota
)

// return name sans any suffix, and any suffix the name had.
// ( i believe names have been normalized by this point )
func findSuffix(name string) (short string, suffix Suffix) {
	short = name // provisional
	for i := 1; i < NumSuffixes; i++ {
		n := Suffix(i)
		if str := n.String(); strings.HasSuffix(name, str) {
			// the suffix string doesnt have the padding space so check that manually
			if end := len(name) - len(str) - 1; end > 0 && name[end] == ' ' {
				short, suffix = name[:end], n
				break
			}
		}
	}
	return
}
