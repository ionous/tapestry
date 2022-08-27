// Code generated by "makeops"; edit at your own risk.
package rel

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// ReciprocalOf Returns the implied relative of a noun (ex. the source in a one-to-many relation.).
type ReciprocalOf struct {
	Via    RelationName `if:"label=_"`
	Object rt.TextEval  `if:"label=object"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.TextEval = (*ReciprocalOf)(nil)

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

func (op *ReciprocalOf) Marshal(m jsn.Marshaler) error {
	return ReciprocalOf_Marshal(m, op)
}

type ReciprocalOf_Slice []ReciprocalOf

func (op *ReciprocalOf_Slice) GetType() string { return ReciprocalOf_Type }

func (op *ReciprocalOf_Slice) Marshal(m jsn.Marshaler) error {
	return ReciprocalOf_Repeats_Marshal(m, (*[]ReciprocalOf)(op))
}

func (op *ReciprocalOf_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *ReciprocalOf_Slice) SetSize(cnt int) {
	var els []ReciprocalOf
	if cnt >= 0 {
		els = make(ReciprocalOf_Slice, cnt)
	}
	(*op) = els
}

func (op *ReciprocalOf_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return ReciprocalOf_Marshal(m, &(*op)[i])
}

func ReciprocalOf_Repeats_Marshal(m jsn.Marshaler, vals *[]ReciprocalOf) error {
	return jsn.RepeatBlock(m, (*ReciprocalOf_Slice)(vals))
}

func ReciprocalOf_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]ReciprocalOf) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = ReciprocalOf_Repeats_Marshal(m, pv)
	}
	return
}

type ReciprocalOf_Flow struct{ ptr *ReciprocalOf }

func (n ReciprocalOf_Flow) GetType() string      { return ReciprocalOf_Type }
func (n ReciprocalOf_Flow) GetLede() string      { return "reciprocal" }
func (n ReciprocalOf_Flow) GetFlow() interface{} { return n.ptr }
func (n ReciprocalOf_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*ReciprocalOf); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func ReciprocalOf_Optional_Marshal(m jsn.Marshaler, pv **ReciprocalOf) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = ReciprocalOf_Marshal(m, *pv)
	} else if !enc {
		var v ReciprocalOf
		if err = ReciprocalOf_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func ReciprocalOf_Marshal(m jsn.Marshaler, val *ReciprocalOf) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(ReciprocalOf_Flow{val}); err == nil {
		e0 := m.MarshalKey("", ReciprocalOf_Field_Via)
		if e0 == nil {
			e0 = RelationName_Marshal(m, &val.Via)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", ReciprocalOf_Field_Via))
		}
		e1 := m.MarshalKey("object", ReciprocalOf_Field_Object)
		if e1 == nil {
			e1 = rt.TextEval_Marshal(m, &val.Object)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", ReciprocalOf_Field_Object))
		}
		m.EndBlock()
	}
	return
}

// ReciprocalsOf Returns the implied relative of a noun (ex. the sources of a many-to-many relation.).
type ReciprocalsOf struct {
	Via    RelationName `if:"label=_"`
	Object rt.TextEval  `if:"label=object"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.TextListEval = (*ReciprocalsOf)(nil)

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

func (op *ReciprocalsOf) Marshal(m jsn.Marshaler) error {
	return ReciprocalsOf_Marshal(m, op)
}

type ReciprocalsOf_Slice []ReciprocalsOf

func (op *ReciprocalsOf_Slice) GetType() string { return ReciprocalsOf_Type }

func (op *ReciprocalsOf_Slice) Marshal(m jsn.Marshaler) error {
	return ReciprocalsOf_Repeats_Marshal(m, (*[]ReciprocalsOf)(op))
}

func (op *ReciprocalsOf_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *ReciprocalsOf_Slice) SetSize(cnt int) {
	var els []ReciprocalsOf
	if cnt >= 0 {
		els = make(ReciprocalsOf_Slice, cnt)
	}
	(*op) = els
}

func (op *ReciprocalsOf_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return ReciprocalsOf_Marshal(m, &(*op)[i])
}

func ReciprocalsOf_Repeats_Marshal(m jsn.Marshaler, vals *[]ReciprocalsOf) error {
	return jsn.RepeatBlock(m, (*ReciprocalsOf_Slice)(vals))
}

func ReciprocalsOf_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]ReciprocalsOf) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = ReciprocalsOf_Repeats_Marshal(m, pv)
	}
	return
}

type ReciprocalsOf_Flow struct{ ptr *ReciprocalsOf }

func (n ReciprocalsOf_Flow) GetType() string      { return ReciprocalsOf_Type }
func (n ReciprocalsOf_Flow) GetLede() string      { return "reciprocals" }
func (n ReciprocalsOf_Flow) GetFlow() interface{} { return n.ptr }
func (n ReciprocalsOf_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*ReciprocalsOf); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func ReciprocalsOf_Optional_Marshal(m jsn.Marshaler, pv **ReciprocalsOf) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = ReciprocalsOf_Marshal(m, *pv)
	} else if !enc {
		var v ReciprocalsOf
		if err = ReciprocalsOf_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func ReciprocalsOf_Marshal(m jsn.Marshaler, val *ReciprocalsOf) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(ReciprocalsOf_Flow{val}); err == nil {
		e0 := m.MarshalKey("", ReciprocalsOf_Field_Via)
		if e0 == nil {
			e0 = RelationName_Marshal(m, &val.Via)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", ReciprocalsOf_Field_Via))
		}
		e1 := m.MarshalKey("object", ReciprocalsOf_Field_Object)
		if e1 == nil {
			e1 = rt.TextEval_Marshal(m, &val.Object)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", ReciprocalsOf_Field_Object))
		}
		m.EndBlock()
	}
	return
}

// Relate Relate two nouns.
type Relate struct {
	Object   rt.TextEval  `if:"label=_"`
	ToObject rt.TextEval  `if:"label=to"`
	Via      RelationName `if:"label=via"`
	Markup   map[string]any
}

// User implemented slots:
var _ rt.Execute = (*Relate)(nil)

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

func (op *Relate) Marshal(m jsn.Marshaler) error {
	return Relate_Marshal(m, op)
}

type Relate_Slice []Relate

func (op *Relate_Slice) GetType() string { return Relate_Type }

func (op *Relate_Slice) Marshal(m jsn.Marshaler) error {
	return Relate_Repeats_Marshal(m, (*[]Relate)(op))
}

func (op *Relate_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Relate_Slice) SetSize(cnt int) {
	var els []Relate
	if cnt >= 0 {
		els = make(Relate_Slice, cnt)
	}
	(*op) = els
}

func (op *Relate_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Relate_Marshal(m, &(*op)[i])
}

func Relate_Repeats_Marshal(m jsn.Marshaler, vals *[]Relate) error {
	return jsn.RepeatBlock(m, (*Relate_Slice)(vals))
}

func Relate_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Relate) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = Relate_Repeats_Marshal(m, pv)
	}
	return
}

type Relate_Flow struct{ ptr *Relate }

func (n Relate_Flow) GetType() string      { return Relate_Type }
func (n Relate_Flow) GetLede() string      { return Relate_Type }
func (n Relate_Flow) GetFlow() interface{} { return n.ptr }
func (n Relate_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Relate); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Relate_Optional_Marshal(m jsn.Marshaler, pv **Relate) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Relate_Marshal(m, *pv)
	} else if !enc {
		var v Relate
		if err = Relate_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Relate_Marshal(m jsn.Marshaler, val *Relate) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(Relate_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Relate_Field_Object)
		if e0 == nil {
			e0 = rt.TextEval_Marshal(m, &val.Object)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Relate_Field_Object))
		}
		e1 := m.MarshalKey("to", Relate_Field_ToObject)
		if e1 == nil {
			e1 = rt.TextEval_Marshal(m, &val.ToObject)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", Relate_Field_ToObject))
		}
		e2 := m.MarshalKey("via", Relate_Field_Via)
		if e2 == nil {
			e2 = RelationName_Marshal(m, &val.Via)
		}
		if e2 != nil && e2 != jsn.Missing {
			m.Error(errutil.New(e2, "in flow at", Relate_Field_Via))
		}
		m.EndBlock()
	}
	return
}

// RelationName requires a user-specified string.
type RelationName struct {
	Str string
}

func (op *RelationName) String() string {
	return op.Str
}

func (*RelationName) Compose() composer.Spec {
	return composer.Spec{
		Name:        RelationName_Type,
		Uses:        composer.Type_Str,
		OpenStrings: true,
	}
}

const RelationName_Type = "relation_name"

func (op *RelationName) Marshal(m jsn.Marshaler) error {
	return RelationName_Marshal(m, op)
}

func RelationName_Optional_Marshal(m jsn.Marshaler, val *RelationName) (err error) {
	var zero RelationName
	if enc := m.IsEncoding(); !enc || val.Str != zero.Str {
		err = RelationName_Marshal(m, val)
	}
	return
}

func RelationName_Marshal(m jsn.Marshaler, val *RelationName) (err error) {
	return m.MarshalValue(RelationName_Type, &val.Str)
}

type RelationName_Slice []RelationName

func (op *RelationName_Slice) GetType() string { return RelationName_Type }

func (op *RelationName_Slice) Marshal(m jsn.Marshaler) error {
	return RelationName_Repeats_Marshal(m, (*[]RelationName)(op))
}

func (op *RelationName_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RelationName_Slice) SetSize(cnt int) {
	var els []RelationName
	if cnt >= 0 {
		els = make(RelationName_Slice, cnt)
	}
	(*op) = els
}

func (op *RelationName_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RelationName_Marshal(m, &(*op)[i])
}

func RelationName_Repeats_Marshal(m jsn.Marshaler, vals *[]RelationName) error {
	return jsn.RepeatBlock(m, (*RelationName_Slice)(vals))
}

func RelationName_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RelationName) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RelationName_Repeats_Marshal(m, pv)
	}
	return
}

// RelativeOf Returns the relative of a noun (ex. the target of a one-to-one relation.).
type RelativeOf struct {
	Via    RelationName `if:"label=_"`
	Object rt.TextEval  `if:"label=object"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.TextEval = (*RelativeOf)(nil)

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

func (op *RelativeOf) Marshal(m jsn.Marshaler) error {
	return RelativeOf_Marshal(m, op)
}

type RelativeOf_Slice []RelativeOf

func (op *RelativeOf_Slice) GetType() string { return RelativeOf_Type }

func (op *RelativeOf_Slice) Marshal(m jsn.Marshaler) error {
	return RelativeOf_Repeats_Marshal(m, (*[]RelativeOf)(op))
}

func (op *RelativeOf_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RelativeOf_Slice) SetSize(cnt int) {
	var els []RelativeOf
	if cnt >= 0 {
		els = make(RelativeOf_Slice, cnt)
	}
	(*op) = els
}

func (op *RelativeOf_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RelativeOf_Marshal(m, &(*op)[i])
}

func RelativeOf_Repeats_Marshal(m jsn.Marshaler, vals *[]RelativeOf) error {
	return jsn.RepeatBlock(m, (*RelativeOf_Slice)(vals))
}

func RelativeOf_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RelativeOf) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RelativeOf_Repeats_Marshal(m, pv)
	}
	return
}

type RelativeOf_Flow struct{ ptr *RelativeOf }

func (n RelativeOf_Flow) GetType() string      { return RelativeOf_Type }
func (n RelativeOf_Flow) GetLede() string      { return "relative" }
func (n RelativeOf_Flow) GetFlow() interface{} { return n.ptr }
func (n RelativeOf_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*RelativeOf); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func RelativeOf_Optional_Marshal(m jsn.Marshaler, pv **RelativeOf) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RelativeOf_Marshal(m, *pv)
	} else if !enc {
		var v RelativeOf
		if err = RelativeOf_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RelativeOf_Marshal(m jsn.Marshaler, val *RelativeOf) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(RelativeOf_Flow{val}); err == nil {
		e0 := m.MarshalKey("", RelativeOf_Field_Via)
		if e0 == nil {
			e0 = RelationName_Marshal(m, &val.Via)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RelativeOf_Field_Via))
		}
		e1 := m.MarshalKey("object", RelativeOf_Field_Object)
		if e1 == nil {
			e1 = rt.TextEval_Marshal(m, &val.Object)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", RelativeOf_Field_Object))
		}
		m.EndBlock()
	}
	return
}

// RelativesOf Returns the relatives of a noun as a list of names (ex. the targets of one-to-many relation).
type RelativesOf struct {
	Via    RelationName `if:"label=_"`
	Object rt.TextEval  `if:"label=object"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.TextListEval = (*RelativesOf)(nil)

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

func (op *RelativesOf) Marshal(m jsn.Marshaler) error {
	return RelativesOf_Marshal(m, op)
}

type RelativesOf_Slice []RelativesOf

func (op *RelativesOf_Slice) GetType() string { return RelativesOf_Type }

func (op *RelativesOf_Slice) Marshal(m jsn.Marshaler) error {
	return RelativesOf_Repeats_Marshal(m, (*[]RelativesOf)(op))
}

func (op *RelativesOf_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RelativesOf_Slice) SetSize(cnt int) {
	var els []RelativesOf
	if cnt >= 0 {
		els = make(RelativesOf_Slice, cnt)
	}
	(*op) = els
}

func (op *RelativesOf_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RelativesOf_Marshal(m, &(*op)[i])
}

func RelativesOf_Repeats_Marshal(m jsn.Marshaler, vals *[]RelativesOf) error {
	return jsn.RepeatBlock(m, (*RelativesOf_Slice)(vals))
}

func RelativesOf_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RelativesOf) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RelativesOf_Repeats_Marshal(m, pv)
	}
	return
}

type RelativesOf_Flow struct{ ptr *RelativesOf }

func (n RelativesOf_Flow) GetType() string      { return RelativesOf_Type }
func (n RelativesOf_Flow) GetLede() string      { return "relatives" }
func (n RelativesOf_Flow) GetFlow() interface{} { return n.ptr }
func (n RelativesOf_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*RelativesOf); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func RelativesOf_Optional_Marshal(m jsn.Marshaler, pv **RelativesOf) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RelativesOf_Marshal(m, *pv)
	} else if !enc {
		var v RelativesOf
		if err = RelativesOf_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RelativesOf_Marshal(m jsn.Marshaler, val *RelativesOf) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(RelativesOf_Flow{val}); err == nil {
		e0 := m.MarshalKey("", RelativesOf_Field_Via)
		if e0 == nil {
			e0 = RelationName_Marshal(m, &val.Via)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RelativesOf_Field_Via))
		}
		e1 := m.MarshalKey("object", RelativesOf_Field_Object)
		if e1 == nil {
			e1 = rt.TextEval_Marshal(m, &val.Object)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", RelativesOf_Field_Object))
		}
		m.EndBlock()
	}
	return
}

var Slats = []composer.Composer{
	(*ReciprocalOf)(nil),
	(*ReciprocalsOf)(nil),
	(*Relate)(nil),
	(*RelationName)(nil),
	(*RelativeOf)(nil),
	(*RelativesOf)(nil),
}

var Signatures = map[uint64]interface{}{
	13187065565571362208: (*RelationName)(nil),  /* RelationName: */
	6987621383789599381:  (*ReciprocalOf)(nil),  /* text_eval=Reciprocal:object: */
	16170704865359856399: (*ReciprocalsOf)(nil), /* text_list_eval=Reciprocals:object: */
	15160920709871392391: (*Relate)(nil),        /* execute=Relate:to:via: */
	14535552277213572673: (*RelativeOf)(nil),    /* text_eval=Relative:object: */
	13180339401044333799: (*RelativesOf)(nil),   /* text_list_eval=Relatives:object: */
}
