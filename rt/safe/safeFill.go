package safe

import (
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

func FillRecord(run rt.Runtime, rec *g.Record, keys []string, vals []g.Value) (ret *g.Record, err error) {
	if nk, nv := len(keys), len(vals); nv < nk {
		err = errutil.New("too many keys")
	} else if nv == 0 {
		ret = rec
	} else {
		kind := rec.Kind()
		if lf, e := NewLabelFinder(run, kind); e != nil {
			err = e
		} else {
			indexedArgs := nv - nk
			for i, val := range vals {
				var key string
				if ofs := i - indexedArgs; ofs >= 0 {
					key = keys[ofs]
				}
				if at, e := lf.FindNext(key); e != nil {
					err = errutil.Fmt("%w while reading arg %d(%s)", e, i, key)
					break
				} else if at < 0 {
					break
				} else if convertedVal, e := AutoConvert(run, kind.Field(at), val); e != nil {
					err = e
					break
				} else if e := rec.SetIndexedField(at, convertedVal); e != nil {
					// note: set indexed field assigns without copying
					// but get value copies out, so this should be okay.
					err = errutil.Fmt("%w while setting arg %d(%s)", e, i, key)
					break
				}
			}
			if err == nil {
				ret = rec
			}
		}
	}
	return
}
