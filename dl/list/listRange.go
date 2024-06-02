package list

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *Range) GetNumList(run rt.Runtime) (ret rt.Value, err error) {
	if vs, e := op.getNumList(run); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = vs
	}
	return
}

func (op *Range) getNumList(run rt.Runtime) (ret rt.Value, err error) {
	if start, e := safe.GetOptionalNumber(run, op.From, 1); e != nil {
		err = e
	} else if stop, e := safe.GetOptionalNumber(run, op.To, start.Float()); e != nil {
		err = e
	} else if step, e := safe.GetOptionalNumber(run, op.ByStep, 1); e != nil {
		err = e
	} else if step := step.Int(); step == 0 {
		err = errutil.New("Range error, step cannot be zero")
	} else {
		ret = &ranger{start: start.Int(), stop: stop.Int(), step: step}
	}
	return
}

// ranger is a PanicValue where every method panics except type and affinity.
type ranger struct {
	rt.PanicValue
	start, stop, step int
}

// Affinity of a range is a number list.
func (n ranger) Affinity() affine.Affinity {
	return affine.NumList
}

// Type returns "range".
func (n ranger) Type() string {
	return "range"
}

// Index computes the i(th) step of the range.
func (n ranger) Index(i int) rt.Value {
	v := n.start + i*n.step
	return rt.IntOf(v)
}

// Len returns the total number of steps.
func (n ranger) Len() (ret int) {
	if diff := (n.stop - n.start + n.step); (n.step < 0) == (diff < 0) {
		ret = diff / n.step
	}
	return
}
