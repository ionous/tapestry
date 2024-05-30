package core

import (
	"math"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

func (op *AbsValue) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := op.abs(run); e != nil {
		err = cmdError(op, e)
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
		err = cmdError(op, e)
	} else {
		ret = rt.FloatOf(a + b)
	}
	return
}

func (op *SubtractValue) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		ret = rt.FloatOf(a - b)
	}
	return
}

func (op *MultiplyValue) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		ret = rt.FloatOf(a * b)
	}
	return
}

func (op *DivideValue) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else if math.Abs(b) <= 1e-5 {
		e := errutil.New("QuotientOf second operand is too small")
		err = cmdError(op, e)
	} else {
		ret = rt.FloatOf(a / b)
	}
	return
}

func (op *ModValue) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if a, b, e := getPair(run, op.A, op.B); e != nil {
		err = cmdError(op, e)
	} else {
		mod := math.Mod(a, b)
		ret = rt.FloatOf(mod)
	}
	return
}

func (op *Increment) Execute(run rt.Runtime) (err error) {
	if _, e := inc(run, op.Target, op.Step, 1.0); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Increment) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := inc(run, op.Target, op.Step, 1.0); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *Decrement) Execute(run rt.Runtime) (err error) {
	if _, e := inc(run, op.Target, op.Step, -1.0); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *Decrement) GetNum(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := inc(run, op.Target, op.Step, -1.0); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *IncrementAspect) Execute(run rt.Runtime) (err error) {
	if _, e := adjustAspect(run, op.Target, op.Aspect, op.Step, op.Clamp, incTrait); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *IncrementAspect) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := adjustAspect(run, op.Target, op.Aspect, op.Step, op.Clamp, incTrait); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (op *DecrementAspect) Execute(run rt.Runtime) (err error) {
	if _, e := adjustAspect(run, op.Target, op.Aspect, op.Step, op.Clamp, decTrait); e != nil {
		err = cmdError(op, e)
	}
	return
}

func (op *DecrementAspect) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := adjustAspect(run, op.Target, op.Aspect, op.Step, op.Clamp, decTrait); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func incTrait(curr, step, max int, clamp bool) (ret int) {
	if next := curr + step; next < max {
		ret = next
	} else if clamp {
		ret = max - 1 // saturate
	} else {
		ret = next % max
	}
	return
}

func decTrait(curr, step, max int, clamp bool) (ret int) {
	if next := curr - step; next >= 0 {
		ret = next
	} else if clamp {
		ret = 0 // clip
	} else {
		ret = max + (next % max) // -1 % 5= -1; 5 + (-1 % 5) = 4
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

func inc(run rt.Runtime, tgt assign.Address, val rt.NumEval, dir float64) (ret rt.Value, err error) {
	if at, e := assign.GetReference(run, tgt); e != nil {
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

// where aspect is the name of the aspect.
func adjustAspect(run rt.Runtime, target, aspect rt.TextEval, steps rt.NumEval, clamps rt.BoolEval,
	update func(curr, step, max int, wrap bool) int) (ret rt.Value, err error) {
	if tgt, e := safe.GetText(run, target); e != nil {
		err = e
	} else if field, e := safe.GetText(run, aspect); e != nil {
		err = e
	} else if step, e := safe.GetOptionalNumber(run, steps, 1); e != nil {
		err = e
	} else if clamp, e := safe.GetOptionalBool(run, clamps, false); e != nil {
		err = e
	} else if obj, e := run.GetField(meta.ObjectId, tgt.String()); e != nil {
		err = e
	} else if currTrait, e := run.GetField(obj.String(), field.String()); e != nil {
		err = e
	} else if aspect, e := run.GetKindByName(field.String()); e != nil {
		err = e
	} else if !aspect.Implements(kindsOf.Aspect.String()) {
		err = errutil.Fmt("field %q is not an aspect", field.String())
	} else {
		prev := aspect.FieldIndex(currTrait.String())
		index := update(prev, step.Int(), aspect.FieldCount(), clamp.Bool())
		newTrait := rt.StringOf(aspect.Field(index).Name)
		if e := run.SetField(obj.String(), field.String(), newTrait); e != nil {
			err = e
		} else {
			ret = newTrait
		}
	}
	return
}
