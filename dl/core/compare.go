package core

import (
	"git.sr.ht/~ionous/tapestry/affine"
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

const defaultTolerance = 1e-3

func (op *CompareNum) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.GetNumber(run, op.A); e != nil {
		err = cmdErrorCtx(op, "A", e)
	} else if b, e := safe.GetNumber(run, op.B); e != nil {
		err = cmdErrorCtx(op, "B", e)
	} else {
		// fix: should tolerance be a literal? should optional literals be pointers?
		tolerance := defaultTolerance
		if op.Tolerance > 0.0 {
			tolerance = op.Tolerance
		}
		cmp := op.Is.Compare()
		res := cmp.CompareFloat(a.Float(), b.Float(), tolerance)
		ret = g.BoolOf(res)
	}
	return
}

func (op *CompareText) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.GetText(run, op.A); e != nil {
		err = cmdErrorCtx(op, "A", e)
	} else if b, e := safe.GetText(run, op.B); e != nil {
		err = cmdErrorCtx(op, "B", e)
	} else {
		cmp := op.Is.Compare()
		res := cmp.CompareString(a.String(), b.String())
		ret = g.BoolOf(res)
	}
	return
}

func (op *CompareValue) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.GetAssignment(run, op.A); e != nil {
		err = cmdErrorCtx(op, "A", e)
	} else if bn, e := safe.GetAssignment(run, op.B); e != nil {
		err = cmdErrorCtx(op, "B", e)
	} else if bv, e := safe.ConvertValue(run, bn, a.Affinity()); e != nil {
		err = cmdErrorCtx(op, "convert", e)
	} else if d, e := compareValues(a, bv, defaultTolerance); e != nil {
		err = cmdErrorCtx(op, "compare", e)
	} else {
		cmp := op.Is.Compare()
		res := cmp.diff(d)
		ret = g.BoolOf(res)
	}
	return
}

// fix: look around at other languages a bit to see what they do....
func compareValues(a, b g.Value, tolerance float64) (ret int, err error) {
	switch a.Affinity() {
	case affine.Bool:
		ret = compareBool(a.Bool(), b.Bool())
	case affine.Number:
		ret = compareFloats(a.Float(), b.Float(), tolerance)
	case affine.Text:
		ret = compareStrings(a.String(), b.String())
	case affine.Record:
		a, b := a.Record(), b.Record()
		an, bn := safeRecordName(a), safeRecordName(b)
		if d := compareStrings(an, bn); d != 0 {
			ret = d
		} else {
			// fix: need to report on the mismatch
			// an optional log statement?
			for i, cnt := 0, a.Kind().NumField(); i < cnt; i++ {
				if d := compareBool(a.HasValue(i), b.HasValue(i)); d != 0 {
					ret = d
					break
				} else if av, e := a.GetIndexedField(i); e != nil {
					err = e
					break
				} else if bv, e := b.GetIndexedField(i); e != nil {
					err = e
					break
				} else if d, e := compareValues(av, bv, tolerance); e != nil {
					err = e
					break
				} else if d != 0 {
					ret = d
					break
				}
			}
		}
	case affine.NumList, affine.TextList, affine.RecordList:
		an, bn := a.Len(), b.Len()
		if d := compareInt(an, bn); d != 0 {
			ret = d
		} else {
			for i := 0; i < an; i++ {
				if d, e := compareValues(a.Index(i), b.Index(i), defaultTolerance); e != nil {
					err = e
				} else if d != 0 {
					ret = d
					break
				}
			}
		}
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
