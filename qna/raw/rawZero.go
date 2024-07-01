package raw

import "git.sr.ht/~ionous/tapestry/rt"

// copied from qdb.
// note: this doesnt properly determine the default trait for an aspect
// weave works around this by providing the correct default value in the db
func zeroAssignment(ft rt.Field, idx int) (ret rt.Assignment, err error) {
	if init := ft.Init; init != nil {
		ret = init
	} else if v, e := rt.ZeroField(ft.Affinity, ft.Type, idx); e != nil {
		err = e
	} else {
		ret = rt.AssignValue(v)
	}
	return
}
