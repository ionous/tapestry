package pattern

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

func getResult(rec *g.Record, res int, aff affine.Affinity) (ret g.Value, err error) {
	if n := res; n >= 0 {
		if res, e := rec.GetIndexedField(n); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else if e := safe.Check(res, aff); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else if len(aff) == 0 {
			// the caller expects nothing but we have a return value.
			if res.Affinity() == affine.Text {
				core.HackTillTemplatesCanEvaluatePatternTypes = res.String()
			} else {
				err = errutil.New("the caller expects nothing but we returned", aff)
			}
		} else {
			ret = res
		}
	} else if len(aff) != 0 {
		err = errutil.New("caller expected", aff, "returned nothing")
	}
	return
}
