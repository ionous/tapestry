package math

import (
	"math"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *AbsValue) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.abs(run); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}

func (op *AbsValue) abs(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := safe.GetNum(run, op.Value); e != nil {
		err = errutil.New("couldnt get value, because", e)
	} else {
		abs := math.Abs(v.Float())
		ret = rt.FloatOf(abs)
	}
	return
}

func (op *AddValue) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.FloatOf(a + b)
	}
	return
}

func (op *SubtractValue) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.FloatOf(a - b)
	}
	return
}

func (op *MultiplyValue) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = rt.FloatOf(a * b)
	}
	return
}

func (op *DivideValue) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmd.Error(op, e)
	} else if math.Abs(b) <= 1e-5 {
		e := errutil.New("QuotientOf second operand is too small")
		err = cmd.Error(op, e)
	} else {
		ret = rt.FloatOf(a / b)
	}
	return
}

func (op *ModValue) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmd.Error(op, e)
	} else {
		mod := math.Mod(a, b)
		ret = rt.FloatOf(mod)
	}
	return
}

func (op *Increment) Execute(run rt.Runtime) (err error) {
	if _, e := inc(run, op.Target, op.Step, 1.0); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *Increment) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := inc(run, op.Target, op.Step, 1.0); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}

func (op *Decrement) Execute(run rt.Runtime) (err error) {
	if _, e := inc(run, op.Target, op.Step, -1.0); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *Decrement) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := inc(run, op.Target, op.Step, -1.0); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}

func getPair(run rt.Runtime, a, b rt.NumEval) (reta, retb float64, err error) {
	if a, e := safe.GetNum(run, a); e != nil {
		err = errutil.New("couldnt get first operand, because", e)
	} else if b, e := safe.GetNum(run, b); e != nil {
		err = errutil.New("couldnt get second operand, because", e)
	} else {
		reta, retb = a.Float(), b.Float()
	}
	return
}

func inc(run rt.Runtime, tgt rt.Address, val rt.NumEval, dir float64) (ret rt.Value, err error) {
	if at, e := safe.GetReference(run, tgt); e != nil {
		err = e
	} else if b, e := safe.GetOptionalNumber(run, val, 1); e != nil {
		err = e
	} else if a, e := at.GetValue(); e != nil {
		err = e
	} else if e := safe.Check(a, affine.Num); e != nil {
		err = e
	} else {
		v := rt.FloatOf(a.Float() + (dir * b.Float()))
		if e := at.SetValue(v); e != nil {
			err = e
		} else {
			ret = v
		}
	}
	return
}
