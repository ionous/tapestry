package core

import (
	"math"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *AbsValue) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := op.abs(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *AbsValue) abs(run rt.Runtime) (ret g.Value, err error) {
	if v, e := safe.GetNumber(run, op.Value); e != nil {
		err = errutil.New("couldnt get value, because", e)
	} else {
		abs := math.Abs(v.Float())
		ret = g.FloatOf(abs)
	}
	return
}

func (op *AddValue) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.FloatOf(a + b)
	}
	return
}

func (op *SubtractValue) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.FloatOf(a - b)
	}
	return
}

func (op *MultiplyValue) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.FloatOf(a * b)
	}
	return
}

func (op *DivideValue) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else if math.Abs(b) <= 1e-5 {
		e := errutil.New("QuotientOf second operand is too small")
		err = cmdError(op, e)
	} else {
		ret = g.FloatOf(a / b)
	}
	return
}

func (op *ModValue) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		mod := math.Mod(a, b)
		ret = g.FloatOf(mod)
	}
	return
}

func (op *Increment) Execute(run rt.Runtime) (err error) {
	if _, e := inc(run, op.Target, op.Value, 1.0); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Increment) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := inc(run, op.Target, op.Value, 1.0); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *Decrement) Execute(run rt.Runtime) (err error) {
	if _, e := inc(run, op.Target, op.Value, -1.0); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Decrement) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if v, e := inc(run, op.Target, op.Value, -1.0); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func getPair(run rt.Runtime, a, b rt.NumberEval) (reta, retb float64, err error) {
	if a, e := safe.GetNumber(run, a); e != nil {
		err = errutil.New("couldnt get first operand, because", e)
	} else if b, e := safe.GetNumber(run, b); e != nil {
		err = errutil.New("couldnt get second operand, because", e)
	} else {
		reta, retb = a.Float(), b.Float()
	}
	return
}

func inc(run rt.Runtime, tgt assign.Address, val rt.NumberEval, dir float64) (ret g.Value, err error) {
	if root, e := assign.GetRootValue(run, tgt); e != nil {
		err = e
	} else if b, e := safe.GetOptionalNumber(run, val, 1); e != nil {
		err = e
	} else if a, e := root.GetCheckedValue(run, affine.Number); e != nil {
		err = e
	} else {
		v := g.FloatOf(a.Float() + (dir * b.Float()))
		if e := root.SetValue(run, v); e != nil {
			err = e
		} else {
			ret = v
		}
	}
	return
}
