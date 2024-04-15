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
func ReadCommaAnd(ws []match.TokenValue) (ret Separator, err error) {
	return countCommaAnd(ws, false)
}

func ReadCommaAndOr(ws []match.TokenValue) (ret Separator, err error) {
	return countCommaAnd(ws, true)
}

func countCommaAnd(ws []match.TokenValue, alsoOr bool) (retFlag Separator, err error) {
	var flag Separator
	var skip int
Loop:
	for _, w := range ws {
		switch w.Token {
		default:
			// anything other than a comma or and?
			break Loop

		case match.Comma:
			if flag != 0 {
				err = makeWordError(w, "unexpected comma")
				break Loop
			}
			flag |= CommaSep
			skip = skip + 1

		case match.String:
			h := w.Hash()
			m := h == keywords.And || (alsoOr && h == keywords.Or)
			if !m {
				break Loop
			} else {
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
	word   match.TokenValue
	reason string
}

func makeWordError(w match.TokenValue, reason string) error {
	return &wordError{w, reason}
}

func (w *wordError) Error() string {
	return fmt.Sprintf("%s in %v at %s", w.reason, w.word.Value, w.word.Pos)
}
