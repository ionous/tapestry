package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"github.com/ionous/errutil"
)

func counterKey(name string) query.Key {
	return query.MakeKey("", meta.Counter, name)
}

// returns 0 if the counter doesnt exist
// only updates the cache on setCounter.
func (run *Runner) getCounter(name string) (ret rt.Value, err error) {
	key := counterKey(name)
	if v, e := run.unpackDynamicValue(key, affine.Num, ""); e == nil {
		ret = v
	} else if e != errMissing {
		err = e
	} else {
		ret = rt.Zero
	}
	return
}

func (run *Runner) setCounter(name string, val rt.Value) (err error) {
	if aff := val.Affinity(); aff != affine.Num {
		err = errutil.Fmt("counter %q expected a number got %s", name, aff)
	} else {
		key := counterKey(name) // no need to copy: numbers are primitives
		run.dynamicVals.Store(key, UserValue{val})
	}
	return
}
