package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/support/match"
)

// flags indicating the presence of a comma followed by an optional and.
type Separator int

const (
	CommaSep Separator = 1 << iota
	AndSep
)

func (s Separator) Len() int {
	return int(s&CommaSep) + int((s&AndSep)>>(AndSep-1))
}

// tests for an optional `, and` or `and`
// and returns the number of words necessary to skip them.
// errors if it detects an unexpected sequence of commas or ands.
// note: when words were hashed, commas became their own Word.
func ReadCommaAnd(ws []match.Word) (ret Separator, err error) {
	return countCommaAnd(ws, false)
}

func ReadCommaAndOr(ws []match.Word) (ret Separator, err error) {
	return countCommaAnd(ws, true)
}

func countCommaAnd(ws []match.Word, alsoOr bool) (retFlag Separator, err error) {
	var flag Separator
	var skip int
Loop:
	for _, w := range ws {
		switch h := w.Hash(); {
		default:
			// anything other than a comma or and?
			break Loop

		case h == keywords.Comma:
			if flag != 0 {
				err = makeWordError(w, "unexpected comma")
				break Loop
			}
			flag |= CommaSep
			skip = skip + 1

		case h == keywords.And || (alsoOr && h == keywords.Or):
			if flag == 0 {
				// start = i
			} else if flag&AndSep != 0 {
				err = makeWordError(w, "unexpected and")
				break Loop
			}
			flag |= AndSep
			skip++
		}
	}
	// nothingness is okay, but not nothingness after a comma or and.
	if err == nil && flag != 0 {
		if skip == len(ws) {
			err = makeWordError(ws[0], "unexpected ending")
		} else {
			retFlag = flag
		}
	}
	return
}

type wordError struct {
	word   match.Word
	reason string
}

func makeWordError(w match.Word, reason string) error {
	return &wordError{w, reason}
}

func (w *wordError) Error() string {
	// i suppose if you wanted to be evil, you would unsafe pointer this string
	// back it up by start to get the actual position
	return fmt.Sprintf("%s in %q", w.reason, w.word)
}
