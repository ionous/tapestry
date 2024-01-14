// Package decode unpacks stored programs and values from byte slices
package decode

import (
	"encoding/json"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

type Decoder decode.Decoder

func NewDecoder(signatures decode.SignatureTable) *Decoder {
	// fix: doesnt really have to be a pointer...
	dec := new(decode.Decoder).
		Signatures(signatures...).
		Customize(core.CustomDecoder)
	return (*Decoder)(dec)
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

func (d *Decoder) DecodeProg(b []byte) (ret []rt.Execute, err error) {
	var act rt.Execute_Slice
	if e := d.decodeValue(&act, b); e != nil {
		err = e
	} else {
		ret = act
	}
	return
}

// matches with mdl.marshalAssignment
// the expected eval depends on the affinity (a) of the destination field.
// fix? merge somehow with express.newAssignment? with compact decoding.
func (d *Decoder) DecodeAssignment(a affine.Affinity, b []byte) (ret rt.Assignment, err error) {
	switch a {
	case affine.None:
		var v rt.Assignment
		if e := d.decodeValue(rt.Assignment_Slot{Value: &v}, b); e != nil {
			err = e
		} else {
			ret = v
		}
	case affine.Bool:
		var v rt.BoolEval
		if e := d.decodeValue(rt.BoolEval_Slot{Value: &v}, b); e != nil {
			err = e
		} else {
			ret = &assign.FromBool{Value: v}
		}
	case affine.Number:
		var v rt.NumberEval
		if e := d.decodeValue(rt.NumberEval_Slot{Value: &v}, b); e != nil {
			err = e
		} else {
			ret = &assign.FromNumber{Value: v}
		}
	case affine.Text:
		var v rt.TextEval
		if e := d.decodeValue(rt.TextEval_Slot{Value: &v}, b); e != nil {
			err = e
		} else {
			ret = &assign.FromText{Value: v}
		}
	case affine.NumList:
		var v rt.NumListEval
		if e := d.decodeValue(rt.NumListEval_Slot{Value: &v}, b); e != nil {
			err = e
		} else {
			ret = &assign.FromNumList{Value: v}
		}
	case affine.TextList:
		var v rt.TextListEval
		if e := d.decodeValue(rt.TextListEval_Slot{Value: &v}, b); e != nil {
			err = e
		} else {
			ret = &assign.FromTextList{Value: v}
		}
	case affine.Record:
		var v rt.RecordEval
		if e := d.decodeValue(rt.RecordEval_Slot{Value: &v}, b); e != nil {
			err = e
		} else {
			ret = &assign.FromRecord{Value: v}
		}
	case affine.RecordList:
		var v rt.RecordListEval
		if e := d.decodeValue(rt.RecordListEval_Slot{Value: &v}, b); e != nil {
			err = e
		} else {
			ret = &assign.FromRecordList{Value: v}
		}
	default:
		err = errutil.New("unhandled affinity", a.String())
	}
	return
}

func (d *Decoder) decodeValue(out jsn.Marshalee, b []byte) (err error) {
	if len(b) > 0 {
		var val any
		if e := json.Unmarshal(b, &val); e != nil {
			err = e
		} else {
			dec := (*decode.Decoder)(d)
			err = dec.Decode(out, val)
		}
	}
	return
}
