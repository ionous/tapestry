package rules

import (
	"strings"
)

// words that can follow a rule name
// ex. "taking then continue"
type eventSuffix int

//go:generate stringer -type=eventSuffix -linecomment
const (
	continues eventSuffix = iota // continue
	stops                        // stop
	jumps                        // jump
	//
	numSuffixes = iota
)

const suffixSeparator = " then "

// return name sans any suffix, and any suffix the name had.
// ( i believe names have been normalized by this point )
func findSuffix(name string) (short string, suffix eventSuffix) {
	short, suffix = name, numSuffixes // provisional
	if i := strings.LastIndex(name, suffixSeparator); i > 0 {
		first, rest := name[:i], name[i+len(suffixSeparator):]
		for i := 0; i < numSuffixes; i++ {
			if p := eventSuffix(i); p.String() == rest {
				short, suffix = first, p
				break
			}
		}
	}
	return
}
