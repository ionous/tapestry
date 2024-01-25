// Package decode unpacks stored programs and values from byte slices
package decode

import (
	"encoding/json"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/rtti"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
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
	var act rtti.Execute_Slots
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
		var v rtti.Assignment_Slot
		if e := d.decodeValue(&v, b); e != nil {
			err = e
		} else {
			ret = v.Value
		}
	case affine.Bool:
		var v rtti.BoolEval_Slot
		if e := d.decodeValue(&v, b); e != nil {
			err = e
		} else {
			ret = &assign.FromBool{Value: v.Value}
		}
	case affine.Number:
		var v rtti.NumberEval_Slot
		if e := d.decodeValue(&v, b); e != nil {
			err = e
		} else {
			ret = &assign.FromNumber{Value: v.Value}
		}
	case affine.Text:
		var v rtti.TextEval_Slot
		if e := d.decodeValue(&v, b); e != nil {
			err = e
		} else {
			ret = &assign.FromText{Value: v.Value}
		}
	case affine.NumList:
		var v rtti.NumListEval_Slot
		if e := d.decodeValue(&v, b); e != nil {
			err = e
		} else {
			ret = &assign.FromNumList{Value: v.Value}
		}
	case affine.TextList:
		var v rtti.TextListEval_Slot
		if e := d.decodeValue(&v, b); e != nil {
			err = e
		} else {
			ret = &assign.FromTextList{Value: v.Value}
		}
	case affine.Record:
		var v rtti.RecordEval_Slot
		if e := d.decodeValue(&v, b); e != nil {
			err = e
		} else {
			ret = &assign.FromRecord{Value: v.Value}
		}
	case affine.RecordList:
		var v rtti.RecordListEval_Slot
		if e := d.decodeValue(&v, b); e != nil {
			err = e
		} else {
			ret = &assign.FromRecordList{Value: v.Value}
		}
	default:
		err = errutil.New("unhandled affinity", a.String())
	}
	return
}

func (d *Decoder) decodeValue(out typeinfo.Inspector, b []byte) (err error) {
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
