package qna

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
)

// fix? mdl_assign technically has redundant info --
// because it implicitly stores the assignment type ( FromBool, etc. ) even though the field wont allow anything else
// storing the evals ( re: readEval below ) would eliminate that.
func decodeAssignment(a affine.Affinity, prog []byte, signatures []map[uint64]interface{}) (ret rt.Assignment, err error) {
	if e := core.Decode(rt.Assignment_Slot{&ret}, prog, signatures); e != nil {
		err = e
	}
	return
}

// the expected value depends on the affinity (a) of the destination field.
// fix? if literals implemented GetAssignedValue, then we could use the literal decoder directly
func readLiteralValue(a affine.Affinity, t string, msg []byte) (ret g.Value, err error) {
	// fix? for text, technically this might be an @variable
	// but perhaps we could patch that to a full eval write in the assembler.
	if x, e := literal.ReadLiteral(a, t, msg); e != nil {
		err = e
	} else {
		ret, err = x.GetAssignedValue(nil)
	}
	return
}

// the expected eval depends on the affinity (a) of the destination field.
// fix? merge somehow with express.newAssignment?
func readEvalValue(a affine.Affinity, rawValue []byte, signatures []map[uint64]interface{}) (ret rt.Assignment, err error) {
	switch a {
	case affine.Bool:
		var v rt.BoolEval
		if e := core.Decode(rt.BoolEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromBool{Val: v}
		}
	case affine.Number:
		var v rt.NumberEval
		if e := core.Decode(rt.NumberEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromNum{Val: v}
		}
	case affine.Text:
		var v rt.TextEval
		if e := core.Decode(rt.TextEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromText{Val: v}
		}
	case affine.NumList:
		var v rt.NumListEval
		if e := core.Decode(rt.NumListEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromNumbers{Vals: v}
		}
	case affine.TextList:
		var v rt.TextListEval
		if e := core.Decode(rt.TextListEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromTexts{Vals: v}
		}
	case affine.Record:
		var v rt.RecordEval
		if e := core.Decode(rt.RecordEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromRecord{Val: v}
		}
	case affine.RecordList:
		var v rt.RecordListEval
		if e := core.Decode(rt.RecordListEval_Slot{&v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &core.FromRecords{Vals: v}
		}
	default:
		err = errutil.New("unhandled affinity", a.String())
	}
	return
}
