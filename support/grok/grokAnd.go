package grok

type sepFlag int

const (
	CommaSep sepFlag = 1 << iota
	AndSep
)

// return the span of the and, the number of words to skip after
// returns 0,0,nil if no ands found.
func anyAnd(ws []Word) (retStart, retEnd int, err error) {
	if start, skip, _, e := _findAnd(ws, false /*atStart*/); e != nil {
		err = e
	} else if skip > 0 {
		retStart, retEnd = start, start+skip
	}
	return
}

// tests for an optional `, and` or `and`
// and returns the number of words necessary to skip them.
// errors if it detects an unexpected sequence of commas or ands.
// note: when words were hashed, commas became their own Word.
func countAnd(ws []Word) (retSkip int, retFlag sepFlag, err error) {
	_, retSkip, retFlag, err = _findAnd(ws, true /*atStart*/)
	return
}

func _findAnd(ws []Word, atStart bool) (retStart, retSkip int, retFlag sepFlag, err error) {
	var flag sepFlag
	var start, skip int
Loop:
	for i, w := range ws {
		switch w.hash {
		default:
			// anything other than a comma or and?
			if atStart || flag != 0 {
				break Loop
			}

		case keywords.comma:
			if flag != 0 {
				err = makeWordError(w, "unexpected comma")
				break Loop
			}
			flag |= CommaSep
			start, skip = i, skip+1

		case keywords.and:
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
			retStart, retSkip, retFlag = start, skip, flag
		}
	}
	return
}
