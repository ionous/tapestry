package safe

import (
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// min inclusive, max exclusive
func Range(i, min, max int) (ret int, err error) {
	if i < min {
		ret, err = min, g.Underflow{i, min}
	} else if i >= max {
		ret, err = max-1, g.Overflow{i, max}
	} else {
		ret = i
	}
	return
}
