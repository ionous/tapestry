// Package decode unpacks stored programs and values from byte slices
package decode

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

type Decoder struct {
	signatures cin.Signatures
}

func NewDecoder(signatures cin.Signatures) *Decoder {
	return &Decoder{
		signatures: signatures,
	}
}

func (d *Decoder) DecodeField(b []byte, a affine.Affinity, fieldType string) (ret literal.LiteralValue, err error) {
	return literal.ReadLiteral(a, fieldType, b)
}

func AssignLiteral(v literal.LiteralValue) (ret rt.Assignment) {
	switch v := v.(type) {
	case *literal.BoolValue:
		ret = &assign.FromBool{Value: v}
	case *literal.NumValue:
		ret = &assign.FromNumber{Value: v}
	case *literal.TextValue:
		ret = &assign.FromText{Value: v}
	case *literal.RecordValue:
		ret = &assign.FromRecord{Value: v}
	case *literal.NumValues:
		ret = &assign.FromNumList{Value: v}
	case *literal.TextValues:
		ret = &assign.FromTextList{Value: v}
	case *literal.RecordList:
		ret = &assign.FromRecordList{Value: v}
	default:
		panic(errutil.Fmt("unknown literal %T", v))
	}
	return
}

func (d *Decoder) DecodeAssignment(b []byte) (ret rt.Assignment, err error) {
	// fix? mdl_default technically has redundant info --
	// because it implicitly stores the assignment type ( FromBool, etc. ) even though the field wont allow anything else
	// storing the evals ( re: readEval below ) would eliminate that.
	var out assign.Assignment
	if e := core.Decode(assign.Assignment_Slot{Value: &out}, b, d.signatures); e != nil {
		err = e
	} else {
		ret = out
	}
	return
}

func (d *Decoder) DecodeFilter(b []byte) (ret rt.BoolEval, err error) {
	if e := core.Decode(rt.BoolEval_Slot{Value: &ret}, b, d.signatures); e != nil {
		err = e
	}
	return
}

func (d *Decoder) DecodeProg(b []byte) (ret rt.Execute_Slice, err error) {
	if e := core.Decode(&ret, b, d.signatures); e != nil {
		err = e
	}
	return
}

// the expected eval depends on the affinity (a) of the destination field.
// fix? merge somehow with express.newAssignment? with compact decoding.
func parseEval(a affine.Affinity, rawValue []byte, signatures cin.Signatures) (ret rt.Assignment, err error) {
	switch a {
	case affine.Bool:
		var v rt.BoolEval
		if e := core.Decode(rt.BoolEval_Slot{Value: &v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromBool{Value: v}
		}
	case affine.Number:
		var v rt.NumberEval
		if e := core.Decode(rt.NumberEval_Slot{Value: &v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromNumber{Value: v}
		}
	case affine.Text:
		var v rt.TextEval
		if e := core.Decode(rt.TextEval_Slot{Value: &v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromText{Value: v}
		}
	case affine.NumList:
		var v rt.NumListEval
		if e := core.Decode(rt.NumListEval_Slot{Value: &v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromNumList{Value: v}
		}
	case affine.TextList:
		var v rt.TextListEval
		if e := core.Decode(rt.TextListEval_Slot{Value: &v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromTextList{Value: v}
		}
	case affine.Record:
		var v rt.RecordEval
		if e := core.Decode(rt.RecordEval_Slot{Value: &v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromRecord{Value: v}
		}
	case affine.RecordList:
		var v rt.RecordListEval
		if e := core.Decode(rt.RecordListEval_Slot{Value: &v}, rawValue, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromRecordList{Value: v}
		}
	default:
		err = errutil.New("unhandled affinity", a.String())
	}
	return
}
