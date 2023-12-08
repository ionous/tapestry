// Package decode unpacks stored programs and values from byte slices
package decode

import (
	"encoding/json"

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

func (d *Decoder) DecodeField(a affine.Affinity, b []byte, fieldType string) (ret literal.LiteralValue, err error) {
	var val any // any json type
	if e := json.Unmarshal(b, &val); e != nil {
		err = e
	} else {
		ret, err = literal.ReadLiteral(a, fieldType, val)
	}
	return
}

func (d *Decoder) DecodeAssignment(a affine.Affinity, b []byte) (ret rt.Assignment, err error) {
	var val any
	if e := json.Unmarshal(b, &val); e != nil {
		err = e
	} else {
		ret, err = parseEval(a, val, d.signatures)
	}
	return
}

func (d *Decoder) DecodeProg(b []byte) (ret []rt.Execute, err error) {
	if len(b) > 0 {
		var prog any // usually a slice of messages
		if e := json.Unmarshal(b, &prog); e != nil {
			err = e
		} else {
			var act rt.Execute_Slice
			if e := core.DecodeValue(&act, prog, d.signatures); e != nil {
				err = e
			} else {
				ret = act
			}
		}
	}
	return
}

// matches with mdl.marshalAssignment
// the expected eval depends on the affinity (a) of the destination field.
// fix? merge somehow with express.newAssignment? with compact decoding.
func parseEval(a affine.Affinity, val any, signatures cin.Signatures) (ret rt.Assignment, err error) {
	switch a {
	case affine.None:
		// FIX: why dont they all use this simpler route?
		var v rt.Assignment
		if e := core.DecodeValue(rt.Assignment_Slot{Value: &v}, val, signatures); e != nil {
			err = e
		} else {
			ret = v
		}

	case affine.Bool:
		var v rt.BoolEval
		if e := core.DecodeValue(rt.BoolEval_Slot{Value: &v}, val, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromBool{Value: v}
		}
	case affine.Number:
		var v rt.NumberEval
		if e := core.DecodeValue(rt.NumberEval_Slot{Value: &v}, val, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromNumber{Value: v}
		}
	case affine.Text:
		var v rt.TextEval
		if e := core.DecodeValue(rt.TextEval_Slot{Value: &v}, val, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromText{Value: v}
		}
	case affine.NumList:
		var v rt.NumListEval
		if e := core.DecodeValue(rt.NumListEval_Slot{Value: &v}, val, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromNumList{Value: v}
		}
	case affine.TextList:
		var v rt.TextListEval
		if e := core.DecodeValue(rt.TextListEval_Slot{Value: &v}, val, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromTextList{Value: v}
		}
	case affine.Record:
		var v rt.RecordEval
		if e := core.DecodeValue(rt.RecordEval_Slot{Value: &v}, val, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromRecord{Value: v}
		}
	case affine.RecordList:
		var v rt.RecordListEval
		if e := core.DecodeValue(rt.RecordListEval_Slot{Value: &v}, val, signatures); e != nil {
			err = e
		} else {
			ret = &assign.FromRecordList{Value: v}
		}
	default:
		err = errutil.New("unhandled affinity", a.String())
	}
	return
}
