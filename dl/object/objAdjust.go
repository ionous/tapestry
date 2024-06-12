package object

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/dot"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func (op *IncrementAspect) Execute(run rt.Runtime) (err error) {
	if _, e := adjustAspect(run, op.Target, op.AspectName, op.Step, op.Clamp, incTrait); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *IncrementAspect) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := adjustAspect(run, op.Target, op.AspectName, op.Step, op.Clamp, incTrait); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}

func (op *DecrementAspect) Execute(run rt.Runtime) (err error) {
	if _, e := adjustAspect(run, op.Target, op.AspectName, op.Step, op.Clamp, decTrait); e != nil {
		err = cmd.Error(op, e)
	}
	return
}

func (op *DecrementAspect) GetText(run rt.Runtime) (ret rt.Value, err error) {
	if v, e := adjustAspect(run, op.Target, op.AspectName, op.Step, op.Clamp, decTrait); e != nil {
		err = cmd.Error(op, e)
	} else {
		ret = v
	}
	return
}

// where aspect is the name of the aspect.
func adjustAspect(run rt.Runtime, target rt.Address, aspect rt.TextEval, steps rt.NumEval, clamps rt.BoolEval,
	update func(curr, step, max int, wrap bool) int) (ret rt.Value, err error) {
	if tgt, e := safe.GetReference(run, target); e != nil {
		err = e
	} else if field, e := safe.GetText(run, aspect); e != nil {
		err = e
	} else if step, e := safe.GetOptionalNumber(run, steps, 1); e != nil {
		err = e
	} else if clamp, e := safe.GetOptionalBool(run, clamps, false); e != nil {
		err = e
	} else if traitRef, e := tgt.Dot(dot.Field(field.String())); e != nil {
		err = e
	} else if aspect, e := run.GetKindByName(field.String()); e != nil {
		err = e
	} else if !aspect.Implements(kindsOf.Aspect.String()) {
		err = fmt.Errorf("field %q is not an aspect", field.String())
	} else if currTrait, e := traitRef.GetValue(); e != nil {
		err = e
	} else {
		prev := aspect.FieldIndex(currTrait.String())
		index := update(prev, step.Int(), aspect.FieldCount(), clamp.Bool())
		newTrait := rt.StringOf(aspect.Field(index).Name)
		if e := traitRef.SetValue(newTrait); e != nil {
			err = e
		} else {
			ret = newTrait
		}
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
