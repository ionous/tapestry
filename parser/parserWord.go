package parser

import "strings"

// Word matches one word.
type Word struct {
	Word string
}

func (w *Word) Scan(ctx Context, bounds Bounds, cs Cursor) (ret Result, err error) {
	if word := cs.CurrentWord(); len(word) == 0 {
		err = Underflow{Depth(cs.Pos)}
	} else if !strings.EqualFold(word, w.Word) {
		err = MismatchedWord{w.Word, word, Depth(cs.Pos)}
	} else {
		ret = ResolvedWords{word, 1}
	}
	return
}
