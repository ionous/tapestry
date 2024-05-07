package safe

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
)

// FillRecord fill the passed record with the arguments named by keys and values
// keys are optional, but if they exist must match the order of fields in the record;
// there doesnt have to be be a key for every field; they can be sparse.
// if there are more values than keys, this assumes the first few values are indexed
// ( keys are right justified )
// returns the passed record if there's no error
func FillRecord(run rt.Runtime, rec *rt.Record, keys []string, vals []rt.Value) (ret *rt.Record, err error) {
	if nk, nv := len(keys), len(vals); nv < nk {
		err = errors.New("too many keys")
	} else if nv == 0 {
		ret = rec
	} else {
		if lf, e := NewLabelFinder(run, rec.Kind); e != nil {
			err = e
		} else {
			indexedArgs := nv - nk
			for i, val := range vals {
				var key string
				if ofs := i - indexedArgs; ofs >= 0 {
					key = keys[ofs]
				}
				if at, e := lf.FindNext(key); e != nil {
					err = fmt.Errorf("%w while reading arg %d(%s)", e, i, key)
					break
				} else if at < 0 {
					break
				} else {
					ft := rec.Field(at)
					if convertedVal, e := RectifyText(run, val, ft.Affinity, ft.Type); e != nil {
						err = e
						break
					} else if e := rec.SetIndexedField(at, convertedVal); e != nil {
						// note: set indexed field assigns without copying
						// but get value copies out, so this should be okay.
						err = fmt.Errorf("%w while setting arg %d(%s)", e, i, key)
						break
					}
				}
			} // end for
			if err == nil {
				ret = rec
			}
		}
	}
	return
}
