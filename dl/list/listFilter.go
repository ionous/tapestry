package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

// reduce the passed list in place.
// removes any elements which don't return true for the passed pattern.
// doesn't validate that vs is a list
func FilterByPattern(run rt.Runtime, vs rt.Value, name string) (err error) {
	pat := inflect.Normalize(name)
	out, cnt := 0, vs.Len()
	for i := 0; i < cnt; i++ {
		val := vs.Index(i)
		if keep, e := run.Call(pat, affine.Bool, nil, []rt.Value{val}); e != nil {
			err = e
			break
		} else if keep.Bool() {
			if i != out { // don't move if the indices are the same
				vs.SetIndex(out, val)
			}
			out++
		}
	}
	// chop down the list ( unless we didn't chop anything )
	if err == nil && out != cnt {
		err = vs.Splice(0, out, nil, nil)
	}
	return
}
