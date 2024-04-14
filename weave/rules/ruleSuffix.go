package rules

import (
	"strings"
)

// words that can follow a rule name
// ex. "taking then continue"
type Suffix int

//go:generate stringer -type=Suffix -linecomment
const (
	Continues Suffix = iota // then continue
	Stops                   // then stop
	Jumps                   // then jump
	Begins                  // begins
	Ends                    // end
	//
	NumSuffixes = iota
)

// return name sans any suffix, and any suffix the name had.
// ( i believe names have been normalized by this point )
func findSuffix(name string) (short string, suffix Suffix) {
	short, suffix = name, NumSuffixes // provisional
	for i := 0; i < NumSuffixes; i++ {
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
