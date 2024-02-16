package grok

import "fmt"

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
func CommaAnd(ws []Word) (ret Separator, err error) {
	_, ret, err = countCommaAnd(ws, true /*atStart*/)
	return
}

// scan for an and,  the span of the and, the number of words to skip after
// returns 0,0,nil if no ands found.
func anyAnd(ws []Word) (retStart, retEnd int, err error) {
	if start, sep, e := countCommaAnd(ws, false /*atStart*/); e != nil {
		err = e
	} else if skip := sep.Len(); skip > 0 {
		retStart, retEnd = start, start+skip
	}
	return
}

// when at start is false, keeps scanning forward to find an and.
func countCommaAnd(ws []Word, atStart bool) (retStart int, retFlag Separator, err error) {
	var flag Separator
	var start, skip int
Loop:
	for i, w := range ws {
		switch w.hash {
		default:
			// anything other than a comma or and?
			if atStart || flag != 0 {
				break Loop
			}

		case Keyword.Comma:
			if flag != 0 {
				err = makeWordError(w, "unexpected comma")
				break Loop
			}
			flag |= CommaSep
			start, skip = i, skip+1

		case Keyword.And:
			if flag == 0 {
				start = i
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
			retStart, retFlag = start, flag
		}
	}
	return
}

type wordError struct {
	word   Word
	reason string
}

func makeWordError(w Word, reason string) error {
	return &wordError{w, reason}
}

func (w *wordError) Error() string {
	// i suppose if you wanted to be evil, you would unsafe pointer this string
	// back it up by start to get the actual position
	return fmt.Sprintf("%s in %q", w.reason, w.word.slice)
}
