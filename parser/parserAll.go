package parser

import "github.com/ionous/errutil"

// AllOf matches the passed Scanners in order.
type AllOf struct {
	Match []Scanner
}

func (m *AllOf) Scan(ctx Context, bounds Bounds, cs Cursor) (Result, error) {
	return scanAllOf(ctx, bounds, cs, m.Match)
}

func scanAllOf(ctx Context, bounds Bounds, cs Cursor, match []Scanner) (ret *ResultList, err error) {
	var rl ResultList
	if cnt := len(match); cnt == 0 {
		err = errutil.New("no rules specified for scanning all of")
	} else {
		var i int
		for ; i < cnt; i++ {
			if res, e := match[i].Scan(ctx, bounds, cs); e != nil {
				err = e
				break
			} else {
				rl.addResult(res)
				cs = cs.Skip(res.WordsMatched())
			}
		}
		if i == cnt {
			ret = &rl
		}
	}
	return
}
