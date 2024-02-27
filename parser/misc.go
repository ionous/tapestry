package parser

import (
	"errors"
	"strings"
)

// Words - match one of the specified words.
// if the last word is blank and nothing matches, it will defer the decision
// re: a/an/--
type Words []string

func (ws Words) Scan(ctx Context, bounds Bounds, cs Cursor) (ret Result, err error) {
	if cnt := len(ws); cnt == 0 {
		err = errors.New("nothing to match")
	} else {
		word := cs.CurrentWord()
		if len(word) > 0 {
			for _, w := range ws {
				if strings.EqualFold(word, w) {
					ret = ResolvedWords{word, 1}
					break
				}
			}
		}

		if ret == nil {
			// if the last word to match is blank;
			// then it means we're okay not matching after all
			if last := ws[cnt-1]; len(last) == 0 {
				ret = ResolvedWords{WordCount: 0}
			} else if len(word) > 0 {
				// if the last word was required,
				// and we failed to match valid input
				// that's a mismatch
				want := strings.Join(ws, "|")
				err = MismatchedWord{want, word, Depth(cs.Pos)}
			} else {
				// if the last word was required,
				// and the input was empty;
				// mention the missing input
				err = Underflow{Depth(cs.Pos)}
			}
		}
	}
	return
}
