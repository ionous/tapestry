package list

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

func (op *Find) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if i, e := op.findIndex(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.BoolOf(i >= 0)
	}
	return
}

func (op *Find) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if i, e := op.findIndex(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.IntOf(i + 1)
	}
	return
}

// zero based
func (op *Find) findIndex(run rt.Runtime) (ret int, err error) {
	if vs, e := safe.GetAssignedValue(run, op.List); e != nil {
		err = e
	} else if el := affine.Element(vs.Affinity()); len(el) == 0 {
		err = errutil.New("not a list")
	} else if v, e := safe.GetAssignedValue(run, op.Value); e != nil {
		err = e
	} else if aff := v.Affinity(); aff != el {
		err = errutil.New("expected", el, "have", aff)
	} else {
		switch el {
		case affine.Number:
			ret = -1
			match := v.Float()
			for i, n := range vs.Floats() {
				if n == match { //epsilon?
					ret = i
					break
				}
			}
		case affine.Text:
			ret = -1
			match := v.String()
			for i, n := range vs.Strings() {
				if n == match {
					ret = i
					break
				}
			}
		default:
			// fix?
			err = errutil.New("cant search list of", el)
		}
	}
	return
}
