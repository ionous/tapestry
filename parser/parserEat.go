package parser

// Eat absorbs existing results
type Eat struct {
	Scanner Scanner
}

func (w *Eat) Scan(ctx Context, bounds Bounds, cs Cursor) (ret Result, err error) {
	if res, e := w.Scanner.Scan(ctx, bounds, cs); e != nil {
		err = e
	} else {
		ret = ResolvedWords{WordCount: res.WordsMatched()}
	}
	return
}
