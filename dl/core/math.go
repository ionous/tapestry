package core

import (
	"math"

	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func getPair(run rt.Runtime, a, b rt.NumberEval) (reta, retb float64, err error) {
	if a, e := safe.GetNumber(run, a); e != nil {
		err = errutil.New("couldnt get first operand, because", e)
	} else if b, e := safe.GetOptionalNumber(run, b, 0); e != nil {
		err = errutil.New("couldnt get second operand, because", e)
	} else {
		reta, retb = a.Float(), b.Float()
	}
	return
}

func (op *SumOf) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.FloatOf(a + b)
	}
	return
}

func (op *DiffOf) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.FloatOf(a - b)
	}
	return
}

func (op *ProductOf) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.FloatOf(a * b)
	}
	return
}

func (op *QuotientOf) GetNumber(run rt.Runtime) (ret g.Value, err error) {
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

func (op *RemainderOf) GetNumber(run rt.Runtime) (ret g.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		mod := math.Mod(a, b)
		ret = g.FloatOf(mod)
	}
	return
}
