// Code generated by "makeops"; edit at your own risk.
package rel

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/rt"
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

const ReciprocalOf_Field_Via = "$VIA"
const ReciprocalOf_Field_Object = "$OBJECT"

func (op *ReciprocalOf) Marshal(n jsn.Marshaler) {
	ReciprocalOf_Marshal(n, op)
}

type ReciprocalOf_Slice []ReciprocalOf

func (op *ReciprocalOf_Slice) GetSize() int    { return len(*op) }
func (op *ReciprocalOf_Slice) SetSize(cnt int) { (*op) = make(ReciprocalOf_Slice, cnt) }

func ReciprocalOf_Repeats_Marshal(n jsn.Marshaler, vals *[]ReciprocalOf) {
	if n.RepeatValues(ReciprocalOf_Type, (*ReciprocalOf_Slice)(vals)) {
		for i := range *vals {
			ReciprocalOf_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func ReciprocalOf_Optional_Marshal(n jsn.Marshaler, pv **ReciprocalOf) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		ReciprocalOf_Marshal(n, *pv)
	} else if !enc {
		var v ReciprocalOf
		if ReciprocalOf_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func ReciprocalOf_Marshal(n jsn.Marshaler, val *ReciprocalOf) (okay bool) {
	if okay = n.MapValues("reciprocal", ReciprocalOf_Type); okay {
		if n.MapKey("", ReciprocalOf_Field_Via) {
			value.RelationName_Marshal(n, &val.Via)
		}
		if n.MapKey("object", ReciprocalOf_Field_Object) {
			rt.TextEval_Marshal(n, &val.Object)
		}
		n.EndValues()
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

const ReciprocalsOf_Field_Via = "$VIA"
const ReciprocalsOf_Field_Object = "$OBJECT"

func (op *ReciprocalsOf) Marshal(n jsn.Marshaler) {
	ReciprocalsOf_Marshal(n, op)
}

type ReciprocalsOf_Slice []ReciprocalsOf

func (op *ReciprocalsOf_Slice) GetSize() int    { return len(*op) }
func (op *ReciprocalsOf_Slice) SetSize(cnt int) { (*op) = make(ReciprocalsOf_Slice, cnt) }

func ReciprocalsOf_Repeats_Marshal(n jsn.Marshaler, vals *[]ReciprocalsOf) {
	if n.RepeatValues(ReciprocalsOf_Type, (*ReciprocalsOf_Slice)(vals)) {
		for i := range *vals {
			ReciprocalsOf_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func ReciprocalsOf_Optional_Marshal(n jsn.Marshaler, pv **ReciprocalsOf) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		ReciprocalsOf_Marshal(n, *pv)
	} else if !enc {
		var v ReciprocalsOf
		if ReciprocalsOf_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func ReciprocalsOf_Marshal(n jsn.Marshaler, val *ReciprocalsOf) (okay bool) {
	if okay = n.MapValues("reciprocals", ReciprocalsOf_Type); okay {
		if n.MapKey("", ReciprocalsOf_Field_Via) {
			value.RelationName_Marshal(n, &val.Via)
		}
		if n.MapKey("object", ReciprocalsOf_Field_Object) {
			rt.TextEval_Marshal(n, &val.Object)
		}
		n.EndValues()
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

const Relate_Field_Object = "$OBJECT"
const Relate_Field_ToObject = "$TO_OBJECT"
const Relate_Field_Via = "$VIA"

func (op *Relate) Marshal(n jsn.Marshaler) {
	Relate_Marshal(n, op)
}

type Relate_Slice []Relate

func (op *Relate_Slice) GetSize() int    { return len(*op) }
func (op *Relate_Slice) SetSize(cnt int) { (*op) = make(Relate_Slice, cnt) }

func Relate_Repeats_Marshal(n jsn.Marshaler, vals *[]Relate) {
	if n.RepeatValues(Relate_Type, (*Relate_Slice)(vals)) {
		for i := range *vals {
			Relate_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func Relate_Optional_Marshal(n jsn.Marshaler, pv **Relate) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		Relate_Marshal(n, *pv)
	} else if !enc {
		var v Relate
		if Relate_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func Relate_Marshal(n jsn.Marshaler, val *Relate) (okay bool) {
	if okay = n.MapValues(Relate_Type, Relate_Type); okay {
		if n.MapKey("", Relate_Field_Object) {
			rt.TextEval_Marshal(n, &val.Object)
		}
		if n.MapKey("to", Relate_Field_ToObject) {
			rt.TextEval_Marshal(n, &val.ToObject)
		}
		if n.MapKey("via", Relate_Field_Via) {
			value.RelationName_Marshal(n, &val.Via)
		}
		n.EndValues()
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

const RelativeOf_Field_Via = "$VIA"
const RelativeOf_Field_Object = "$OBJECT"

func (op *RelativeOf) Marshal(n jsn.Marshaler) {
	RelativeOf_Marshal(n, op)
}

type RelativeOf_Slice []RelativeOf

func (op *RelativeOf_Slice) GetSize() int    { return len(*op) }
func (op *RelativeOf_Slice) SetSize(cnt int) { (*op) = make(RelativeOf_Slice, cnt) }

func RelativeOf_Repeats_Marshal(n jsn.Marshaler, vals *[]RelativeOf) {
	if n.RepeatValues(RelativeOf_Type, (*RelativeOf_Slice)(vals)) {
		for i := range *vals {
			RelativeOf_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func RelativeOf_Optional_Marshal(n jsn.Marshaler, pv **RelativeOf) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		RelativeOf_Marshal(n, *pv)
	} else if !enc {
		var v RelativeOf
		if RelativeOf_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func RelativeOf_Marshal(n jsn.Marshaler, val *RelativeOf) (okay bool) {
	if okay = n.MapValues("relative", RelativeOf_Type); okay {
		if n.MapKey("", RelativeOf_Field_Via) {
			value.RelationName_Marshal(n, &val.Via)
		}
		if n.MapKey("object", RelativeOf_Field_Object) {
			rt.TextEval_Marshal(n, &val.Object)
		}
		n.EndValues()
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

const RelativesOf_Field_Via = "$VIA"
const RelativesOf_Field_Object = "$OBJECT"

func (op *RelativesOf) Marshal(n jsn.Marshaler) {
	RelativesOf_Marshal(n, op)
}

type RelativesOf_Slice []RelativesOf

func (op *RelativesOf_Slice) GetSize() int    { return len(*op) }
func (op *RelativesOf_Slice) SetSize(cnt int) { (*op) = make(RelativesOf_Slice, cnt) }

func RelativesOf_Repeats_Marshal(n jsn.Marshaler, vals *[]RelativesOf) {
	if n.RepeatValues(RelativesOf_Type, (*RelativesOf_Slice)(vals)) {
		for i := range *vals {
			RelativesOf_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func RelativesOf_Optional_Marshal(n jsn.Marshaler, pv **RelativesOf) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		RelativesOf_Marshal(n, *pv)
	} else if !enc {
		var v RelativesOf
		if RelativesOf_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func RelativesOf_Marshal(n jsn.Marshaler, val *RelativesOf) (okay bool) {
	if okay = n.MapValues("relatives", RelativesOf_Type); okay {
		if n.MapKey("", RelativesOf_Field_Via) {
			value.RelationName_Marshal(n, &val.Via)
		}
		if n.MapKey("object", RelativesOf_Field_Object) {
			rt.TextEval_Marshal(n, &val.Object)
		}
		n.EndValues()
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
