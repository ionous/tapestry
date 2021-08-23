// Code generated by "makeops"; edit at your own risk.
package rel

import (
	"encoding/json"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/export/jsonexp"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

// ReciprocalOf Returns the implied relative of a noun (ex. the source in a one-to-many relation.)
type ReciprocalOf struct {
	Via    value.RelationName `if:"label=_"`
	Object rt.TextEval        `if:"label=object"`
}

func (*ReciprocalOf) Compose() composer.Spec {
	return composer.Spec{
		Name: ReciprocalOf_Type,
		Uses: composer.Type_Flow,
		Lede: "reciprocal",
	}
}

const ReciprocalOf_Type = "reciprocal_of"
const ReciprocalOf_Lede = "reciprocal"
const ReciprocalOf_Field_Via = "$VIA"
const ReciprocalOf_Field_Object = "$OBJECT"

func (op *ReciprocalOf) MarshalCompact(n jsonexp.Context) (ret []byte, err error) {
	return ReciprocalOf_Compact_Marshal(n, op)
}
func (op *ReciprocalOf) UnmarshalCompact(n jsonexp.Context, b []byte) error {
	return ReciprocalOf_Compact_Unmarshal(n, b, op)
}
func (op *ReciprocalOf) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return ReciprocalOf_Detailed_Marshal(n, op)
}
func (op *ReciprocalOf) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return ReciprocalOf_Detailed_Unmarshal(n, b, op)
}

func ReciprocalOf_Compact_Repeats_Marshal(n jsonexp.Context, vals *[]ReciprocalOf) ([]byte, error) {
	return ReciprocalOf_Repeats_Marshal(n, vals, ReciprocalOf_Compact_Marshal)
}
func ReciprocalOf_Detailed_Repeats_Marshal(n jsonexp.Context, vals *[]ReciprocalOf) ([]byte, error) {
	return ReciprocalOf_Repeats_Marshal(n, vals, ReciprocalOf_Detailed_Marshal)
}
func ReciprocalOf_Repeats_Marshal(n jsonexp.Context, vals *[]ReciprocalOf, marshEl func(jsonexp.Context, *ReciprocalOf) ([]byte, error)) (ret []byte, err error) {
	var msgs []json.RawMessage
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		msgs = make([]json.RawMessage, cnt)
		for i, el := range *vals {
			if b, e := marshEl(n, &el); e != nil {
				err = errutil.New(ReciprocalOf_Type, "at", i, "-", e)
				break
			} else {
				msgs[i] = b
			}
		}
	}
	if err == nil {
		ret, err = json.Marshal(msgs)
	}
	return
}

func ReciprocalOf_Compact_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]ReciprocalOf) error {
	return ReciprocalOf_Repeats_Unmarshal(n, b, out, ReciprocalOf_Compact_Unmarshal)
}
func ReciprocalOf_Detailed_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]ReciprocalOf) error {
	return ReciprocalOf_Repeats_Unmarshal(n, b, out, ReciprocalOf_Detailed_Unmarshal)
}
func ReciprocalOf_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]ReciprocalOf, unmarshEl func(jsonexp.Context, []byte, *ReciprocalOf) error) (err error) {
	var vals []ReciprocalOf
	if len(b) > 0 { // generated code collapses optional and empty.
		var msgs []json.RawMessage
		if e := json.Unmarshal(b, &msgs); e != nil {
			err = errutil.New(ReciprocalOf_Type, "-", e)
		} else {
			vals = make([]ReciprocalOf, len(msgs))
			for i, msg := range msgs {
				if e := unmarshEl(n, msg, &vals[i]); e != nil {
					err = errutil.New(ReciprocalOf_Type, "at", i, "-", e)
					break
				}
			}
		}
	}
	if err == nil {
		*out = vals
	}
	return
}

func ReciprocalOf_Compact_Optional_Marshal(n jsonexp.Context, val **ReciprocalOf) (ret []byte, err error) {
	if *val != nil {
		ret, err = ReciprocalOf_Compact_Marshal(n, *val)
	}
	return
}
func ReciprocalOf_Compact_Marshal(n jsonexp.Context, val *ReciprocalOf) (ret []byte, err error) {
	var sig jsonexp.CompactFlow
	sig.WriteLede(ReciprocalOf_Lede)
	if b, e := value.RelationName_Compact_Marshal(n, &val.Via); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("", b)
	}
	if b, e := rt.TextEval_Compact_Marshal(n, &val.Object); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("object", b)
	}
	if err == nil {
		ret, err = sig.MarshalJSON()
	}
	return
}

func ReciprocalOf_Compact_Optional_Unmarshal(n jsonexp.Context, b []byte, out **ReciprocalOf) (err error) {
	if len(b) > 0 {
		var val ReciprocalOf
		if e := ReciprocalOf_Compact_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func ReciprocalOf_Compact_Unmarshal(n jsonexp.Context, b []byte, out *ReciprocalOf) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(ReciprocalOf_Type, "-", e)
	} else if e := value.RelationName_Compact_Unmarshal(n, msg.Fields[ReciprocalOf_Field_Via], &out.Via); e != nil {
		err = errutil.New(ReciprocalOf_Type+"."+ReciprocalOf_Field_Via, "-", e)
	} else if e := rt.TextEval_Compact_Unmarshal(n, msg.Fields[ReciprocalOf_Field_Object], &out.Object); e != nil {
		err = errutil.New(ReciprocalOf_Type+"."+ReciprocalOf_Field_Object, "-", e)
	}
	return
}

func ReciprocalOf_Detailed_Optional_Marshal(n jsonexp.Context, val **ReciprocalOf) (ret []byte, err error) {
	if *val != nil {
		ret, err = ReciprocalOf_Detailed_Marshal(n, *val)
	}
	return
}
func ReciprocalOf_Detailed_Marshal(n jsonexp.Context, val *ReciprocalOf) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := value.RelationName_Detailed_Marshal(n, &val.Via); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[ReciprocalOf_Field_Via] = b
	}

	if b, e := rt.TextEval_Detailed_Marshal(n, &val.Object); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[ReciprocalOf_Field_Object] = b
	}

	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   ReciprocalOf_Type,
			Fields: fields,
		})
	}
	return
}

func ReciprocalOf_Detailed_Optional_Unmarshal(n jsonexp.Context, b []byte, out **ReciprocalOf) (err error) {
	if len(b) > 0 {
		var val ReciprocalOf
		if e := ReciprocalOf_Detailed_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func ReciprocalOf_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *ReciprocalOf) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(ReciprocalOf_Type, "-", e)
	} else if e := value.RelationName_Detailed_Unmarshal(n, msg.Fields[ReciprocalOf_Field_Via], &out.Via); e != nil {
		err = errutil.New(ReciprocalOf_Type+"."+ReciprocalOf_Field_Via, "-", e)
	} else if e := rt.TextEval_Detailed_Unmarshal(n, msg.Fields[ReciprocalOf_Field_Object], &out.Object); e != nil {
		err = errutil.New(ReciprocalOf_Type+"."+ReciprocalOf_Field_Object, "-", e)
	}
	return
}

// ReciprocalsOf Returns the implied relative of a noun (ex. the sources of a many-to-many relation.)
type ReciprocalsOf struct {
	Via    value.RelationName `if:"label=_"`
	Object rt.TextEval        `if:"label=object"`
}

func (*ReciprocalsOf) Compose() composer.Spec {
	return composer.Spec{
		Name: ReciprocalsOf_Type,
		Uses: composer.Type_Flow,
		Lede: "reciprocals",
	}
}

const ReciprocalsOf_Type = "reciprocals_of"
const ReciprocalsOf_Lede = "reciprocals"
const ReciprocalsOf_Field_Via = "$VIA"
const ReciprocalsOf_Field_Object = "$OBJECT"

func (op *ReciprocalsOf) MarshalCompact(n jsonexp.Context) (ret []byte, err error) {
	return ReciprocalsOf_Compact_Marshal(n, op)
}
func (op *ReciprocalsOf) UnmarshalCompact(n jsonexp.Context, b []byte) error {
	return ReciprocalsOf_Compact_Unmarshal(n, b, op)
}
func (op *ReciprocalsOf) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return ReciprocalsOf_Detailed_Marshal(n, op)
}
func (op *ReciprocalsOf) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return ReciprocalsOf_Detailed_Unmarshal(n, b, op)
}

func ReciprocalsOf_Compact_Repeats_Marshal(n jsonexp.Context, vals *[]ReciprocalsOf) ([]byte, error) {
	return ReciprocalsOf_Repeats_Marshal(n, vals, ReciprocalsOf_Compact_Marshal)
}
func ReciprocalsOf_Detailed_Repeats_Marshal(n jsonexp.Context, vals *[]ReciprocalsOf) ([]byte, error) {
	return ReciprocalsOf_Repeats_Marshal(n, vals, ReciprocalsOf_Detailed_Marshal)
}
func ReciprocalsOf_Repeats_Marshal(n jsonexp.Context, vals *[]ReciprocalsOf, marshEl func(jsonexp.Context, *ReciprocalsOf) ([]byte, error)) (ret []byte, err error) {
	var msgs []json.RawMessage
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		msgs = make([]json.RawMessage, cnt)
		for i, el := range *vals {
			if b, e := marshEl(n, &el); e != nil {
				err = errutil.New(ReciprocalsOf_Type, "at", i, "-", e)
				break
			} else {
				msgs[i] = b
			}
		}
	}
	if err == nil {
		ret, err = json.Marshal(msgs)
	}
	return
}

func ReciprocalsOf_Compact_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]ReciprocalsOf) error {
	return ReciprocalsOf_Repeats_Unmarshal(n, b, out, ReciprocalsOf_Compact_Unmarshal)
}
func ReciprocalsOf_Detailed_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]ReciprocalsOf) error {
	return ReciprocalsOf_Repeats_Unmarshal(n, b, out, ReciprocalsOf_Detailed_Unmarshal)
}
func ReciprocalsOf_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]ReciprocalsOf, unmarshEl func(jsonexp.Context, []byte, *ReciprocalsOf) error) (err error) {
	var vals []ReciprocalsOf
	if len(b) > 0 { // generated code collapses optional and empty.
		var msgs []json.RawMessage
		if e := json.Unmarshal(b, &msgs); e != nil {
			err = errutil.New(ReciprocalsOf_Type, "-", e)
		} else {
			vals = make([]ReciprocalsOf, len(msgs))
			for i, msg := range msgs {
				if e := unmarshEl(n, msg, &vals[i]); e != nil {
					err = errutil.New(ReciprocalsOf_Type, "at", i, "-", e)
					break
				}
			}
		}
	}
	if err == nil {
		*out = vals
	}
	return
}

func ReciprocalsOf_Compact_Optional_Marshal(n jsonexp.Context, val **ReciprocalsOf) (ret []byte, err error) {
	if *val != nil {
		ret, err = ReciprocalsOf_Compact_Marshal(n, *val)
	}
	return
}
func ReciprocalsOf_Compact_Marshal(n jsonexp.Context, val *ReciprocalsOf) (ret []byte, err error) {
	var sig jsonexp.CompactFlow
	sig.WriteLede(ReciprocalsOf_Lede)
	if b, e := value.RelationName_Compact_Marshal(n, &val.Via); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("", b)
	}
	if b, e := rt.TextEval_Compact_Marshal(n, &val.Object); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("object", b)
	}
	if err == nil {
		ret, err = sig.MarshalJSON()
	}
	return
}

func ReciprocalsOf_Compact_Optional_Unmarshal(n jsonexp.Context, b []byte, out **ReciprocalsOf) (err error) {
	if len(b) > 0 {
		var val ReciprocalsOf
		if e := ReciprocalsOf_Compact_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func ReciprocalsOf_Compact_Unmarshal(n jsonexp.Context, b []byte, out *ReciprocalsOf) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(ReciprocalsOf_Type, "-", e)
	} else if e := value.RelationName_Compact_Unmarshal(n, msg.Fields[ReciprocalsOf_Field_Via], &out.Via); e != nil {
		err = errutil.New(ReciprocalsOf_Type+"."+ReciprocalsOf_Field_Via, "-", e)
	} else if e := rt.TextEval_Compact_Unmarshal(n, msg.Fields[ReciprocalsOf_Field_Object], &out.Object); e != nil {
		err = errutil.New(ReciprocalsOf_Type+"."+ReciprocalsOf_Field_Object, "-", e)
	}
	return
}

func ReciprocalsOf_Detailed_Optional_Marshal(n jsonexp.Context, val **ReciprocalsOf) (ret []byte, err error) {
	if *val != nil {
		ret, err = ReciprocalsOf_Detailed_Marshal(n, *val)
	}
	return
}
func ReciprocalsOf_Detailed_Marshal(n jsonexp.Context, val *ReciprocalsOf) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := value.RelationName_Detailed_Marshal(n, &val.Via); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[ReciprocalsOf_Field_Via] = b
	}

	if b, e := rt.TextEval_Detailed_Marshal(n, &val.Object); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[ReciprocalsOf_Field_Object] = b
	}

	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   ReciprocalsOf_Type,
			Fields: fields,
		})
	}
	return
}

func ReciprocalsOf_Detailed_Optional_Unmarshal(n jsonexp.Context, b []byte, out **ReciprocalsOf) (err error) {
	if len(b) > 0 {
		var val ReciprocalsOf
		if e := ReciprocalsOf_Detailed_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func ReciprocalsOf_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *ReciprocalsOf) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(ReciprocalsOf_Type, "-", e)
	} else if e := value.RelationName_Detailed_Unmarshal(n, msg.Fields[ReciprocalsOf_Field_Via], &out.Via); e != nil {
		err = errutil.New(ReciprocalsOf_Type+"."+ReciprocalsOf_Field_Via, "-", e)
	} else if e := rt.TextEval_Detailed_Unmarshal(n, msg.Fields[ReciprocalsOf_Field_Object], &out.Object); e != nil {
		err = errutil.New(ReciprocalsOf_Type+"."+ReciprocalsOf_Field_Object, "-", e)
	}
	return
}

// Relate Relate two nouns.
type Relate struct {
	Object   rt.TextEval        `if:"label=_"`
	ToObject rt.TextEval        `if:"label=to"`
	Via      value.RelationName `if:"label=via"`
}

func (*Relate) Compose() composer.Spec {
	return composer.Spec{
		Name: Relate_Type,
		Uses: composer.Type_Flow,
	}
}

const Relate_Type = "relate"
const Relate_Lede = Relate_Type
const Relate_Field_Object = "$OBJECT"
const Relate_Field_ToObject = "$TO_OBJECT"
const Relate_Field_Via = "$VIA"

func (op *Relate) MarshalCompact(n jsonexp.Context) (ret []byte, err error) {
	return Relate_Compact_Marshal(n, op)
}
func (op *Relate) UnmarshalCompact(n jsonexp.Context, b []byte) error {
	return Relate_Compact_Unmarshal(n, b, op)
}
func (op *Relate) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return Relate_Detailed_Marshal(n, op)
}
func (op *Relate) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return Relate_Detailed_Unmarshal(n, b, op)
}

func Relate_Compact_Repeats_Marshal(n jsonexp.Context, vals *[]Relate) ([]byte, error) {
	return Relate_Repeats_Marshal(n, vals, Relate_Compact_Marshal)
}
func Relate_Detailed_Repeats_Marshal(n jsonexp.Context, vals *[]Relate) ([]byte, error) {
	return Relate_Repeats_Marshal(n, vals, Relate_Detailed_Marshal)
}
func Relate_Repeats_Marshal(n jsonexp.Context, vals *[]Relate, marshEl func(jsonexp.Context, *Relate) ([]byte, error)) (ret []byte, err error) {
	var msgs []json.RawMessage
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		msgs = make([]json.RawMessage, cnt)
		for i, el := range *vals {
			if b, e := marshEl(n, &el); e != nil {
				err = errutil.New(Relate_Type, "at", i, "-", e)
				break
			} else {
				msgs[i] = b
			}
		}
	}
	if err == nil {
		ret, err = json.Marshal(msgs)
	}
	return
}

func Relate_Compact_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]Relate) error {
	return Relate_Repeats_Unmarshal(n, b, out, Relate_Compact_Unmarshal)
}
func Relate_Detailed_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]Relate) error {
	return Relate_Repeats_Unmarshal(n, b, out, Relate_Detailed_Unmarshal)
}
func Relate_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]Relate, unmarshEl func(jsonexp.Context, []byte, *Relate) error) (err error) {
	var vals []Relate
	if len(b) > 0 { // generated code collapses optional and empty.
		var msgs []json.RawMessage
		if e := json.Unmarshal(b, &msgs); e != nil {
			err = errutil.New(Relate_Type, "-", e)
		} else {
			vals = make([]Relate, len(msgs))
			for i, msg := range msgs {
				if e := unmarshEl(n, msg, &vals[i]); e != nil {
					err = errutil.New(Relate_Type, "at", i, "-", e)
					break
				}
			}
		}
	}
	if err == nil {
		*out = vals
	}
	return
}

func Relate_Compact_Optional_Marshal(n jsonexp.Context, val **Relate) (ret []byte, err error) {
	if *val != nil {
		ret, err = Relate_Compact_Marshal(n, *val)
	}
	return
}
func Relate_Compact_Marshal(n jsonexp.Context, val *Relate) (ret []byte, err error) {
	var sig jsonexp.CompactFlow
	sig.WriteLede(Relate_Lede)
	if b, e := rt.TextEval_Compact_Marshal(n, &val.Object); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("", b)
	}
	if b, e := rt.TextEval_Compact_Marshal(n, &val.ToObject); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("to", b)
	}
	if b, e := value.RelationName_Compact_Marshal(n, &val.Via); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("via", b)
	}
	if err == nil {
		ret, err = sig.MarshalJSON()
	}
	return
}

func Relate_Compact_Optional_Unmarshal(n jsonexp.Context, b []byte, out **Relate) (err error) {
	if len(b) > 0 {
		var val Relate
		if e := Relate_Compact_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func Relate_Compact_Unmarshal(n jsonexp.Context, b []byte, out *Relate) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(Relate_Type, "-", e)
	} else if e := rt.TextEval_Compact_Unmarshal(n, msg.Fields[Relate_Field_Object], &out.Object); e != nil {
		err = errutil.New(Relate_Type+"."+Relate_Field_Object, "-", e)
	} else if e := rt.TextEval_Compact_Unmarshal(n, msg.Fields[Relate_Field_ToObject], &out.ToObject); e != nil {
		err = errutil.New(Relate_Type+"."+Relate_Field_ToObject, "-", e)
	} else if e := value.RelationName_Compact_Unmarshal(n, msg.Fields[Relate_Field_Via], &out.Via); e != nil {
		err = errutil.New(Relate_Type+"."+Relate_Field_Via, "-", e)
	}
	return
}

func Relate_Detailed_Optional_Marshal(n jsonexp.Context, val **Relate) (ret []byte, err error) {
	if *val != nil {
		ret, err = Relate_Detailed_Marshal(n, *val)
	}
	return
}
func Relate_Detailed_Marshal(n jsonexp.Context, val *Relate) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := rt.TextEval_Detailed_Marshal(n, &val.Object); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[Relate_Field_Object] = b
	}

	if b, e := rt.TextEval_Detailed_Marshal(n, &val.ToObject); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[Relate_Field_ToObject] = b
	}

	if b, e := value.RelationName_Detailed_Marshal(n, &val.Via); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[Relate_Field_Via] = b
	}

	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   Relate_Type,
			Fields: fields,
		})
	}
	return
}

func Relate_Detailed_Optional_Unmarshal(n jsonexp.Context, b []byte, out **Relate) (err error) {
	if len(b) > 0 {
		var val Relate
		if e := Relate_Detailed_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func Relate_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *Relate) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(Relate_Type, "-", e)
	} else if e := rt.TextEval_Detailed_Unmarshal(n, msg.Fields[Relate_Field_Object], &out.Object); e != nil {
		err = errutil.New(Relate_Type+"."+Relate_Field_Object, "-", e)
	} else if e := rt.TextEval_Detailed_Unmarshal(n, msg.Fields[Relate_Field_ToObject], &out.ToObject); e != nil {
		err = errutil.New(Relate_Type+"."+Relate_Field_ToObject, "-", e)
	} else if e := value.RelationName_Detailed_Unmarshal(n, msg.Fields[Relate_Field_Via], &out.Via); e != nil {
		err = errutil.New(Relate_Type+"."+Relate_Field_Via, "-", e)
	}
	return
}

// RelativeOf Returns the relative of a noun (ex. the target of a one-to-one relation.)
type RelativeOf struct {
	Via    value.RelationName `if:"label=_"`
	Object rt.TextEval        `if:"label=object"`
}

func (*RelativeOf) Compose() composer.Spec {
	return composer.Spec{
		Name: RelativeOf_Type,
		Uses: composer.Type_Flow,
		Lede: "relative",
	}
}

const RelativeOf_Type = "relative_of"
const RelativeOf_Lede = "relative"
const RelativeOf_Field_Via = "$VIA"
const RelativeOf_Field_Object = "$OBJECT"

func (op *RelativeOf) MarshalCompact(n jsonexp.Context) (ret []byte, err error) {
	return RelativeOf_Compact_Marshal(n, op)
}
func (op *RelativeOf) UnmarshalCompact(n jsonexp.Context, b []byte) error {
	return RelativeOf_Compact_Unmarshal(n, b, op)
}
func (op *RelativeOf) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return RelativeOf_Detailed_Marshal(n, op)
}
func (op *RelativeOf) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return RelativeOf_Detailed_Unmarshal(n, b, op)
}

func RelativeOf_Compact_Repeats_Marshal(n jsonexp.Context, vals *[]RelativeOf) ([]byte, error) {
	return RelativeOf_Repeats_Marshal(n, vals, RelativeOf_Compact_Marshal)
}
func RelativeOf_Detailed_Repeats_Marshal(n jsonexp.Context, vals *[]RelativeOf) ([]byte, error) {
	return RelativeOf_Repeats_Marshal(n, vals, RelativeOf_Detailed_Marshal)
}
func RelativeOf_Repeats_Marshal(n jsonexp.Context, vals *[]RelativeOf, marshEl func(jsonexp.Context, *RelativeOf) ([]byte, error)) (ret []byte, err error) {
	var msgs []json.RawMessage
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		msgs = make([]json.RawMessage, cnt)
		for i, el := range *vals {
			if b, e := marshEl(n, &el); e != nil {
				err = errutil.New(RelativeOf_Type, "at", i, "-", e)
				break
			} else {
				msgs[i] = b
			}
		}
	}
	if err == nil {
		ret, err = json.Marshal(msgs)
	}
	return
}

func RelativeOf_Compact_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]RelativeOf) error {
	return RelativeOf_Repeats_Unmarshal(n, b, out, RelativeOf_Compact_Unmarshal)
}
func RelativeOf_Detailed_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]RelativeOf) error {
	return RelativeOf_Repeats_Unmarshal(n, b, out, RelativeOf_Detailed_Unmarshal)
}
func RelativeOf_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]RelativeOf, unmarshEl func(jsonexp.Context, []byte, *RelativeOf) error) (err error) {
	var vals []RelativeOf
	if len(b) > 0 { // generated code collapses optional and empty.
		var msgs []json.RawMessage
		if e := json.Unmarshal(b, &msgs); e != nil {
			err = errutil.New(RelativeOf_Type, "-", e)
		} else {
			vals = make([]RelativeOf, len(msgs))
			for i, msg := range msgs {
				if e := unmarshEl(n, msg, &vals[i]); e != nil {
					err = errutil.New(RelativeOf_Type, "at", i, "-", e)
					break
				}
			}
		}
	}
	if err == nil {
		*out = vals
	}
	return
}

func RelativeOf_Compact_Optional_Marshal(n jsonexp.Context, val **RelativeOf) (ret []byte, err error) {
	if *val != nil {
		ret, err = RelativeOf_Compact_Marshal(n, *val)
	}
	return
}
func RelativeOf_Compact_Marshal(n jsonexp.Context, val *RelativeOf) (ret []byte, err error) {
	var sig jsonexp.CompactFlow
	sig.WriteLede(RelativeOf_Lede)
	if b, e := value.RelationName_Compact_Marshal(n, &val.Via); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("", b)
	}
	if b, e := rt.TextEval_Compact_Marshal(n, &val.Object); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("object", b)
	}
	if err == nil {
		ret, err = sig.MarshalJSON()
	}
	return
}

func RelativeOf_Compact_Optional_Unmarshal(n jsonexp.Context, b []byte, out **RelativeOf) (err error) {
	if len(b) > 0 {
		var val RelativeOf
		if e := RelativeOf_Compact_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func RelativeOf_Compact_Unmarshal(n jsonexp.Context, b []byte, out *RelativeOf) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(RelativeOf_Type, "-", e)
	} else if e := value.RelationName_Compact_Unmarshal(n, msg.Fields[RelativeOf_Field_Via], &out.Via); e != nil {
		err = errutil.New(RelativeOf_Type+"."+RelativeOf_Field_Via, "-", e)
	} else if e := rt.TextEval_Compact_Unmarshal(n, msg.Fields[RelativeOf_Field_Object], &out.Object); e != nil {
		err = errutil.New(RelativeOf_Type+"."+RelativeOf_Field_Object, "-", e)
	}
	return
}

func RelativeOf_Detailed_Optional_Marshal(n jsonexp.Context, val **RelativeOf) (ret []byte, err error) {
	if *val != nil {
		ret, err = RelativeOf_Detailed_Marshal(n, *val)
	}
	return
}
func RelativeOf_Detailed_Marshal(n jsonexp.Context, val *RelativeOf) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := value.RelationName_Detailed_Marshal(n, &val.Via); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[RelativeOf_Field_Via] = b
	}

	if b, e := rt.TextEval_Detailed_Marshal(n, &val.Object); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[RelativeOf_Field_Object] = b
	}

	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   RelativeOf_Type,
			Fields: fields,
		})
	}
	return
}

func RelativeOf_Detailed_Optional_Unmarshal(n jsonexp.Context, b []byte, out **RelativeOf) (err error) {
	if len(b) > 0 {
		var val RelativeOf
		if e := RelativeOf_Detailed_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func RelativeOf_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *RelativeOf) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(RelativeOf_Type, "-", e)
	} else if e := value.RelationName_Detailed_Unmarshal(n, msg.Fields[RelativeOf_Field_Via], &out.Via); e != nil {
		err = errutil.New(RelativeOf_Type+"."+RelativeOf_Field_Via, "-", e)
	} else if e := rt.TextEval_Detailed_Unmarshal(n, msg.Fields[RelativeOf_Field_Object], &out.Object); e != nil {
		err = errutil.New(RelativeOf_Type+"."+RelativeOf_Field_Object, "-", e)
	}
	return
}

// RelativesOf Returns the relatives of a noun as a list of names (ex. the targets of one-to-many relation).
type RelativesOf struct {
	Via    value.RelationName `if:"label=_"`
	Object rt.TextEval        `if:"label=object"`
}

func (*RelativesOf) Compose() composer.Spec {
	return composer.Spec{
		Name: RelativesOf_Type,
		Uses: composer.Type_Flow,
		Lede: "relatives",
	}
}

const RelativesOf_Type = "relatives_of"
const RelativesOf_Lede = "relatives"
const RelativesOf_Field_Via = "$VIA"
const RelativesOf_Field_Object = "$OBJECT"

func (op *RelativesOf) MarshalCompact(n jsonexp.Context) (ret []byte, err error) {
	return RelativesOf_Compact_Marshal(n, op)
}
func (op *RelativesOf) UnmarshalCompact(n jsonexp.Context, b []byte) error {
	return RelativesOf_Compact_Unmarshal(n, b, op)
}
func (op *RelativesOf) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return RelativesOf_Detailed_Marshal(n, op)
}
func (op *RelativesOf) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return RelativesOf_Detailed_Unmarshal(n, b, op)
}

func RelativesOf_Compact_Repeats_Marshal(n jsonexp.Context, vals *[]RelativesOf) ([]byte, error) {
	return RelativesOf_Repeats_Marshal(n, vals, RelativesOf_Compact_Marshal)
}
func RelativesOf_Detailed_Repeats_Marshal(n jsonexp.Context, vals *[]RelativesOf) ([]byte, error) {
	return RelativesOf_Repeats_Marshal(n, vals, RelativesOf_Detailed_Marshal)
}
func RelativesOf_Repeats_Marshal(n jsonexp.Context, vals *[]RelativesOf, marshEl func(jsonexp.Context, *RelativesOf) ([]byte, error)) (ret []byte, err error) {
	var msgs []json.RawMessage
	if cnt := len(*vals); cnt > 0 { // generated code collapses optional and empty.
		msgs = make([]json.RawMessage, cnt)
		for i, el := range *vals {
			if b, e := marshEl(n, &el); e != nil {
				err = errutil.New(RelativesOf_Type, "at", i, "-", e)
				break
			} else {
				msgs[i] = b
			}
		}
	}
	if err == nil {
		ret, err = json.Marshal(msgs)
	}
	return
}

func RelativesOf_Compact_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]RelativesOf) error {
	return RelativesOf_Repeats_Unmarshal(n, b, out, RelativesOf_Compact_Unmarshal)
}
func RelativesOf_Detailed_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]RelativesOf) error {
	return RelativesOf_Repeats_Unmarshal(n, b, out, RelativesOf_Detailed_Unmarshal)
}
func RelativesOf_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]RelativesOf, unmarshEl func(jsonexp.Context, []byte, *RelativesOf) error) (err error) {
	var vals []RelativesOf
	if len(b) > 0 { // generated code collapses optional and empty.
		var msgs []json.RawMessage
		if e := json.Unmarshal(b, &msgs); e != nil {
			err = errutil.New(RelativesOf_Type, "-", e)
		} else {
			vals = make([]RelativesOf, len(msgs))
			for i, msg := range msgs {
				if e := unmarshEl(n, msg, &vals[i]); e != nil {
					err = errutil.New(RelativesOf_Type, "at", i, "-", e)
					break
				}
			}
		}
	}
	if err == nil {
		*out = vals
	}
	return
}

func RelativesOf_Compact_Optional_Marshal(n jsonexp.Context, val **RelativesOf) (ret []byte, err error) {
	if *val != nil {
		ret, err = RelativesOf_Compact_Marshal(n, *val)
	}
	return
}
func RelativesOf_Compact_Marshal(n jsonexp.Context, val *RelativesOf) (ret []byte, err error) {
	var sig jsonexp.CompactFlow
	sig.WriteLede(RelativesOf_Lede)
	if b, e := value.RelationName_Compact_Marshal(n, &val.Via); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("", b)
	}
	if b, e := rt.TextEval_Compact_Marshal(n, &val.Object); e != nil {
		err = errutil.Append(err, e)
	} else {
		sig.AddMsg("object", b)
	}
	if err == nil {
		ret, err = sig.MarshalJSON()
	}
	return
}

func RelativesOf_Compact_Optional_Unmarshal(n jsonexp.Context, b []byte, out **RelativesOf) (err error) {
	if len(b) > 0 {
		var val RelativesOf
		if e := RelativesOf_Compact_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func RelativesOf_Compact_Unmarshal(n jsonexp.Context, b []byte, out *RelativesOf) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(RelativesOf_Type, "-", e)
	} else if e := value.RelationName_Compact_Unmarshal(n, msg.Fields[RelativesOf_Field_Via], &out.Via); e != nil {
		err = errutil.New(RelativesOf_Type+"."+RelativesOf_Field_Via, "-", e)
	} else if e := rt.TextEval_Compact_Unmarshal(n, msg.Fields[RelativesOf_Field_Object], &out.Object); e != nil {
		err = errutil.New(RelativesOf_Type+"."+RelativesOf_Field_Object, "-", e)
	}
	return
}

func RelativesOf_Detailed_Optional_Marshal(n jsonexp.Context, val **RelativesOf) (ret []byte, err error) {
	if *val != nil {
		ret, err = RelativesOf_Detailed_Marshal(n, *val)
	}
	return
}
func RelativesOf_Detailed_Marshal(n jsonexp.Context, val *RelativesOf) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := value.RelationName_Detailed_Marshal(n, &val.Via); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[RelativesOf_Field_Via] = b
	}

	if b, e := rt.TextEval_Detailed_Marshal(n, &val.Object); e != nil {
		err = errutil.Append(err, e)
	} else {
		fields[RelativesOf_Field_Object] = b
	}

	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   RelativesOf_Type,
			Fields: fields,
		})
	}
	return
}

func RelativesOf_Detailed_Optional_Unmarshal(n jsonexp.Context, b []byte, out **RelativesOf) (err error) {
	if len(b) > 0 {
		var val RelativesOf
		if e := RelativesOf_Detailed_Unmarshal(n, b, &val); e != nil {
			err = e
		} else {
			*out = &val
		}
	}
	return
}
func RelativesOf_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *RelativesOf) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(RelativesOf_Type, "-", e)
	} else if e := value.RelationName_Detailed_Unmarshal(n, msg.Fields[RelativesOf_Field_Via], &out.Via); e != nil {
		err = errutil.New(RelativesOf_Type+"."+RelativesOf_Field_Via, "-", e)
	} else if e := rt.TextEval_Detailed_Unmarshal(n, msg.Fields[RelativesOf_Field_Object], &out.Object); e != nil {
		err = errutil.New(RelativesOf_Type+"."+RelativesOf_Field_Object, "-", e)
	}
	return
}

var Slats = []composer.Composer{
	(*ReciprocalOf)(nil),
	(*ReciprocalsOf)(nil),
	(*Relate)(nil),
	(*RelativeOf)(nil),
	(*RelativesOf)(nil),
}
