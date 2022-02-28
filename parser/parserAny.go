package parser

import "github.com/ionous/errutil"

// AnyOf matches any one of the passed Scanners; whichever first matches.
type AnyOf struct {
	Match []Scanner
}

// Scan implements Scanner.
func (m *AnyOf) Scan(ctx Context, bounds Bounds, cs Cursor) (ret Result, err error) {
	if cnt := len(m.Match); cnt == 0 {
		err = errutil.New("no rules specified for any of")
	} else {
		i, errorDepth := 0, -1 // keep the most informative error
		for ; i < cnt; i++ {
			if res, e := m.Match[i].Scan(ctx, bounds, cs); e == nil {
				ret, err = res, nil
				break
			} else if d := DepthOf(e); d >= errorDepth {
				err, errorDepth = e, d // we track the "deepest" / latest error.
				// keep looking for success
			}
		}
	}
	return
}
