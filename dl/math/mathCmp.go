package math

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

const defaultTolerance = 1e-3

func (op *CompareNum) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if a, e := safe.GetNum(run, op.A); e != nil {
		err = cmd.ErrorCtx(op, "A", e)
	} else if b, e := safe.GetNum(run, op.B); e != nil {
		err = cmd.ErrorCtx(op, "B", e)
	} else {
		// fix: should tolerance be a literal? should optional literals be pointers?
		tolerance := defaultTolerance
		if op.Tolerance > 0.0 {
			tolerance = op.Tolerance
		}
		cmp := op.Is.Compare()
		res := cmp.CompareFloat(a.Float(), b.Float(), tolerance)
		ret = rt.BoolOf(res)
	}
	return
}

func (op *CompareText) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if a, e := safe.GetText(run, op.A); e != nil {
		err = cmd.ErrorCtx(op, "A", e)
	} else if b, e := safe.GetText(run, op.B); e != nil {
		err = cmd.ErrorCtx(op, "B", e)
	} else {
		cmp := op.Is.Compare()
		res := cmp.CompareString(a.String(), b.String())
		ret = rt.BoolOf(res)
	}
	return
}

func (op *CompareValue) GetBool(run rt.Runtime) (ret rt.Value, err error) {
	if a, e := safe.GetAssignment(run, op.A); e != nil {
		err = cmd.ErrorCtx(op, "A", e)
	} else if bn, e := safe.GetAssignment(run, op.B); e != nil {
		err = cmd.ErrorCtx(op, "B", e)
	} else if bv, e := safe.ConvertValue(run, bn, a.Affinity()); e != nil {
		err = cmd.ErrorCtx(op, "convert", e)
	} else if d, e := compareValues(a, bv, defaultTolerance); e != nil {
		err = cmd.ErrorCtx(op, "compare", e)
	} else {
		cmp := op.Is.Compare()
		res := cmp.diff(d)
		ret = rt.BoolOf(res)
	}
	return
}

// fix: look around at other languages a bit to see what they do....
func compareValues(a, b rt.Value, tolerance float64) (ret int, err error) {
	switch a.Affinity() {
	case affine.Bool:
		ret = compareBool(a.Bool(), b.Bool())
	case affine.Num:
		ret = compareFloats(a.Float(), b.Float(), tolerance)
	case affine.Text:
		ret = compareStrings(a.String(), b.String())
	case affine.Record:
		// fix: maybe compare serialized versions ( raw bytes ) instead
		a, b := a.Record(), b.Record()
		if d := compareStrings(a.Name(), b.Name()); d != 0 {
			ret = d
		} else {
			for i, cnt := 0, a.FieldCount(); i < cnt && ret != 0; i++ {
				// eat errors ( esp. NilRecord )
				av, _ := a.GetIndexedField(i)
				bv, _ := b.GetIndexedField(i)
				if av == nil || bv == nil {
					if av == nil {
						ret = -1
					} else if bv == nil {
						ret = 1
					}
				} else {
					if d, e := compareValues(av, bv, tolerance); e != nil {
						err = e
						break
					} else {
						ret = d
					}
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
func (op Comparison) Compare() (ret CompareType) {
	switch op {
	case C_Comparison_EqualTo:
		ret = Compare_EqualTo

	case C_Comparison_OtherThan:
		ret = Compare_GreaterThan | Compare_LessThan

	case C_Comparison_GreaterThan:
		ret = Compare_GreaterThan

	case C_Comparison_LessThan:
		ret = Compare_LessThan

	case C_Comparison_AtLeast:
		ret = Compare_GreaterThan | Compare_EqualTo

	case C_Comparison_AtMost:
		ret = Compare_LessThan | Compare_EqualTo
	}
	return
}
