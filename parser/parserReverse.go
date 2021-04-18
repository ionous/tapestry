package parser

// Reverse swaps the first and last matches after scanning as per "AllOf".
// Generally, the first and last scanners are nouns.
type Reverse struct {
	Match []Scanner
}

func (m *Reverse) Scan(ctx Context, bounds Bounds, cs Cursor) (ret Result, err error) {
	if rev, e := scanAllOf(ctx, bounds, cs, m.Match); e != nil {
		err = e
	} else if cnt := len(rev.list); cnt > 1 {
		rev.list[0], rev.list[cnt-1] = rev.list[cnt-1], rev.list[0]
		ret = rev
	}
	return
}
