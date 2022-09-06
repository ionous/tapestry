package core

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// backwards compat ish
var Equal = Comparison{Str: Comparison_EqualTo}
var Unequal = Comparison{Str: Comparison_OtherThan}
var AtLeast = Comparison{Str: Comparison_AtLeast}
var GreaterThan = Comparison{Str: Comparison_GreaterThan}
var LessThan = Comparison{Str: Comparison_LessThan}

func (op *CompareNum) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if src, e := safe.GetNumber(run, op.A); e != nil {
		err = cmdErrorCtx(op, "A", e)
	} else if tgt, e := safe.GetNumber(run, op.B); e != nil {
		err = cmdErrorCtx(op, "B", e)
	} else {
		// fix: atleast for numbers, optional values should be pointers
		// ( requires expansion of the gomake templates )
		tolerance := 1e-3
		if op.Tolerance > 0.0 {
			tolerance = op.Tolerance
		}
		res := op.Is.Compare().CompareFloat(src.Float()-tgt.Float(), tolerance)
		ret = g.BoolOf(res)
	}
	return
}

func (op *CompareText) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if src, e := safe.GetText(run, op.A); e != nil {
		err = cmdErrorCtx(op, "A", e)
	} else if tgt, e := safe.GetText(run, op.B); e != nil {
		err = cmdErrorCtx(op, "B", e)
	} else {
		c := strings.Compare(src.String(), tgt.String())
		res := op.Is.Compare().CompareInt(c)
		ret = g.BoolOf(res)
	}
	return
}

// return flags to help compare numbers
func (op *Comparison) Compare() (ret CompareType) {
	switch op.Str {
	case Comparison_EqualTo:
		ret = Compare_EqualTo

	case Comparison_OtherThan:
		ret = Compare_GreaterThan | Compare_LessThan

	case Comparison_GreaterThan:
		ret = Compare_GreaterThan

	case Comparison_LessThan:
		ret = Compare_LessThan

	case Comparison_AtLeast:
		ret = Compare_GreaterThan | Compare_EqualTo

	case Comparison_AtMost:
		ret = Compare_LessThan | Compare_EqualTo
	}
	return
}
