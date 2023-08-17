// Code generated by "makeops"; edit at your own risk.
package render

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

const RenderEval_Type = "render_eval"

var RenderEval_Optional_Marshal = RenderEval_Marshal

type RenderEval_Slot struct{ Value *RenderEval }

func (at RenderEval_Slot) Marshal(m jsn.Marshaler) (err error) {
	if err = m.MarshalBlock(at); err == nil {
		if a, ok := at.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}
func (at RenderEval_Slot) GetType() string              { return RenderEval_Type }
func (at RenderEval_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at RenderEval_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(RenderEval)
	return
}

func RenderEval_Marshal(m jsn.Marshaler, ptr *RenderEval) (err error) {
	slot := RenderEval_Slot{ptr}
	return slot.Marshal(m)
}

type RenderEval_Slice []RenderEval

func (op *RenderEval_Slice) GetType() string { return RenderEval_Type }

func (op *RenderEval_Slice) Marshal(m jsn.Marshaler) error {
	return RenderEval_Repeats_Marshal(m, (*[]RenderEval)(op))
}

func (op *RenderEval_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RenderEval_Slice) SetSize(cnt int) {
	var els []RenderEval
	if cnt >= 0 {
		els = make(RenderEval_Slice, cnt)
	}
	(*op) = els
}

func (op *RenderEval_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RenderEval_Marshal(m, &(*op)[i])
}

func RenderEval_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderEval) error {
	return jsn.RepeatBlock(m, (*RenderEval_Slice)(vals))
}

func RenderEval_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RenderEval) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RenderEval_Repeats_Marshal(m, pv)
	}
	return
}

// RenderName Handles changing a template like {.boombip} into text.
// If the name is a variable containing an object name: return the printed object name ( via "print name" );
// if the name is a variable with some other text: return that text;
// if the name isn't a variable but refers to some object: return that object's printed object name;
// otherwise, its an error.
type RenderName struct {
	Name   string `if:"label=_,type=text"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.TextEval = (*RenderName)(nil)

func (*RenderName) Compose() composer.Spec {
	return composer.Spec{
		Name: RenderName_Type,
		Uses: composer.Type_Flow,
	}
}

const RenderName_Type = "render_name"
const RenderName_Field_Name = "$NAME"

func (op *RenderName) Marshal(m jsn.Marshaler) error {
	return RenderName_Marshal(m, op)
}

type RenderName_Slice []RenderName

func (op *RenderName_Slice) GetType() string { return RenderName_Type }

func (op *RenderName_Slice) Marshal(m jsn.Marshaler) error {
	return RenderName_Repeats_Marshal(m, (*[]RenderName)(op))
}

func (op *RenderName_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RenderName_Slice) SetSize(cnt int) {
	var els []RenderName
	if cnt >= 0 {
		els = make(RenderName_Slice, cnt)
	}
	(*op) = els
}

func (op *RenderName_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RenderName_Marshal(m, &(*op)[i])
}

func RenderName_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderName) error {
	return jsn.RepeatBlock(m, (*RenderName_Slice)(vals))
}

func RenderName_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RenderName) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RenderName_Repeats_Marshal(m, pv)
	}
	return
}

type RenderName_Flow struct{ ptr *RenderName }

func (n RenderName_Flow) GetType() string      { return RenderName_Type }
func (n RenderName_Flow) GetLede() string      { return RenderName_Type }
func (n RenderName_Flow) GetFlow() interface{} { return n.ptr }
func (n RenderName_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*RenderName); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func RenderName_Optional_Marshal(m jsn.Marshaler, pv **RenderName) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RenderName_Marshal(m, *pv)
	} else if !enc {
		var v RenderName
		if err = RenderName_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RenderName_Marshal(m jsn.Marshaler, val *RenderName) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(RenderName_Flow{val}); err == nil {
		e0 := m.MarshalKey("", RenderName_Field_Name)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RenderName_Field_Name))
		}
		m.EndBlock()
	}
	return
}

// RenderPattern A version of core's call pattern
// that figures out how to evaluate its arguments at runtime.
type RenderPattern struct {
	PatternName string       `if:"label=_,type=text"`
	Render      []RenderEval `if:"label=render"`
	Markup      map[string]any
}

// User implemented slots:
var _ rt.BoolEval = (*RenderPattern)(nil)
var _ rt.TextEval = (*RenderPattern)(nil)
var _ RenderEval = (*RenderPattern)(nil)

func (*RenderPattern) Compose() composer.Spec {
	return composer.Spec{
		Name: RenderPattern_Type,
		Uses: composer.Type_Flow,
		Lede: "render",
	}
}

const RenderPattern_Type = "render_pattern"
const RenderPattern_Field_PatternName = "$PATTERN_NAME"
const RenderPattern_Field_Render = "$RENDER"

func (op *RenderPattern) Marshal(m jsn.Marshaler) error {
	return RenderPattern_Marshal(m, op)
}

type RenderPattern_Slice []RenderPattern

func (op *RenderPattern_Slice) GetType() string { return RenderPattern_Type }

func (op *RenderPattern_Slice) Marshal(m jsn.Marshaler) error {
	return RenderPattern_Repeats_Marshal(m, (*[]RenderPattern)(op))
}

func (op *RenderPattern_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RenderPattern_Slice) SetSize(cnt int) {
	var els []RenderPattern
	if cnt >= 0 {
		els = make(RenderPattern_Slice, cnt)
	}
	(*op) = els
}

func (op *RenderPattern_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RenderPattern_Marshal(m, &(*op)[i])
}

func RenderPattern_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderPattern) error {
	return jsn.RepeatBlock(m, (*RenderPattern_Slice)(vals))
}

func RenderPattern_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RenderPattern) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RenderPattern_Repeats_Marshal(m, pv)
	}
	return
}

type RenderPattern_Flow struct{ ptr *RenderPattern }

func (n RenderPattern_Flow) GetType() string      { return RenderPattern_Type }
func (n RenderPattern_Flow) GetLede() string      { return "render" }
func (n RenderPattern_Flow) GetFlow() interface{} { return n.ptr }
func (n RenderPattern_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*RenderPattern); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func RenderPattern_Optional_Marshal(m jsn.Marshaler, pv **RenderPattern) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RenderPattern_Marshal(m, *pv)
	} else if !enc {
		var v RenderPattern
		if err = RenderPattern_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RenderPattern_Marshal(m jsn.Marshaler, val *RenderPattern) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(RenderPattern_Flow{val}); err == nil {
		e0 := m.MarshalKey("", RenderPattern_Field_PatternName)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Marshal(m, &val.PatternName)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RenderPattern_Field_PatternName))
		}
		e1 := m.MarshalKey("render", RenderPattern_Field_Render)
		if e1 == nil {
			e1 = RenderEval_Repeats_Marshal(m, &val.Render)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", RenderPattern_Field_Render))
		}
		m.EndBlock()
	}
	return
}

// RenderRef Pull a value from name that might refer either to a variable, or to an object.
// If the name is an object, returns the object id.
type RenderRef struct {
	Name   rt.TextEval  `if:"label=_"`
	Dot    []assign.Dot `if:"label=dot,optional"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.BoolEval = (*RenderRef)(nil)
var _ rt.NumberEval = (*RenderRef)(nil)
var _ rt.TextEval = (*RenderRef)(nil)
var _ rt.RecordEval = (*RenderRef)(nil)
var _ rt.NumListEval = (*RenderRef)(nil)
var _ rt.TextListEval = (*RenderRef)(nil)
var _ rt.RecordListEval = (*RenderRef)(nil)
var _ RenderEval = (*RenderRef)(nil)

func (*RenderRef) Compose() composer.Spec {
	return composer.Spec{
		Name: RenderRef_Type,
		Uses: composer.Type_Flow,
	}
}

const RenderRef_Type = "render_ref"
const RenderRef_Field_Name = "$NAME"
const RenderRef_Field_Dot = "$DOT"

func (op *RenderRef) Marshal(m jsn.Marshaler) error {
	return RenderRef_Marshal(m, op)
}

type RenderRef_Slice []RenderRef

func (op *RenderRef_Slice) GetType() string { return RenderRef_Type }

func (op *RenderRef_Slice) Marshal(m jsn.Marshaler) error {
	return RenderRef_Repeats_Marshal(m, (*[]RenderRef)(op))
}

func (op *RenderRef_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RenderRef_Slice) SetSize(cnt int) {
	var els []RenderRef
	if cnt >= 0 {
		els = make(RenderRef_Slice, cnt)
	}
	(*op) = els
}

func (op *RenderRef_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RenderRef_Marshal(m, &(*op)[i])
}

func RenderRef_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderRef) error {
	return jsn.RepeatBlock(m, (*RenderRef_Slice)(vals))
}

func RenderRef_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RenderRef) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RenderRef_Repeats_Marshal(m, pv)
	}
	return
}

type RenderRef_Flow struct{ ptr *RenderRef }

func (n RenderRef_Flow) GetType() string      { return RenderRef_Type }
func (n RenderRef_Flow) GetLede() string      { return RenderRef_Type }
func (n RenderRef_Flow) GetFlow() interface{} { return n.ptr }
func (n RenderRef_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*RenderRef); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func RenderRef_Optional_Marshal(m jsn.Marshaler, pv **RenderRef) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RenderRef_Marshal(m, *pv)
	} else if !enc {
		var v RenderRef
		if err = RenderRef_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RenderRef_Marshal(m jsn.Marshaler, val *RenderRef) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(RenderRef_Flow{val}); err == nil {
		e0 := m.MarshalKey("", RenderRef_Field_Name)
		if e0 == nil {
			e0 = rt.TextEval_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RenderRef_Field_Name))
		}
		e1 := m.MarshalKey("dot", RenderRef_Field_Dot)
		if e1 == nil {
			e1 = assign.Dot_Optional_Repeats_Marshal(m, &val.Dot)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", RenderRef_Field_Dot))
		}
		m.EndBlock()
	}
	return
}

// RenderResponse Generate text in a replaceable manner.
type RenderResponse struct {
	Name   string      `if:"label=_,type=text"`
	Text   rt.TextEval `if:"label=text,optional"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.Execute = (*RenderResponse)(nil)
var _ rt.TextEval = (*RenderResponse)(nil)

func (*RenderResponse) Compose() composer.Spec {
	return composer.Spec{
		Name: RenderResponse_Type,
		Uses: composer.Type_Flow,
	}
}

const RenderResponse_Type = "render_response"
const RenderResponse_Field_Name = "$NAME"
const RenderResponse_Field_Text = "$TEXT"

func (op *RenderResponse) Marshal(m jsn.Marshaler) error {
	return RenderResponse_Marshal(m, op)
}

type RenderResponse_Slice []RenderResponse

func (op *RenderResponse_Slice) GetType() string { return RenderResponse_Type }

func (op *RenderResponse_Slice) Marshal(m jsn.Marshaler) error {
	return RenderResponse_Repeats_Marshal(m, (*[]RenderResponse)(op))
}

func (op *RenderResponse_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RenderResponse_Slice) SetSize(cnt int) {
	var els []RenderResponse
	if cnt >= 0 {
		els = make(RenderResponse_Slice, cnt)
	}
	(*op) = els
}

func (op *RenderResponse_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RenderResponse_Marshal(m, &(*op)[i])
}

func RenderResponse_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderResponse) error {
	return jsn.RepeatBlock(m, (*RenderResponse_Slice)(vals))
}

func RenderResponse_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RenderResponse) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RenderResponse_Repeats_Marshal(m, pv)
	}
	return
}

type RenderResponse_Flow struct{ ptr *RenderResponse }

func (n RenderResponse_Flow) GetType() string      { return RenderResponse_Type }
func (n RenderResponse_Flow) GetLede() string      { return RenderResponse_Type }
func (n RenderResponse_Flow) GetFlow() interface{} { return n.ptr }
func (n RenderResponse_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*RenderResponse); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func RenderResponse_Optional_Marshal(m jsn.Marshaler, pv **RenderResponse) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RenderResponse_Marshal(m, *pv)
	} else if !enc {
		var v RenderResponse
		if err = RenderResponse_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RenderResponse_Marshal(m jsn.Marshaler, val *RenderResponse) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(RenderResponse_Flow{val}); err == nil {
		e0 := m.MarshalKey("", RenderResponse_Field_Name)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RenderResponse_Field_Name))
		}
		e1 := m.MarshalKey("text", RenderResponse_Field_Text)
		if e1 == nil {
			e1 = rt.TextEval_Optional_Marshal(m, &val.Text)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", RenderResponse_Field_Text))
		}
		m.EndBlock()
	}
	return
}

// RenderValue Pull a value from an assignment of unknown affinity.
type RenderValue struct {
	Value  assign.Assignment `if:"label=_"`
	Markup map[string]any
}

// User implemented slots:
var _ RenderEval = (*RenderValue)(nil)

func (*RenderValue) Compose() composer.Spec {
	return composer.Spec{
		Name: RenderValue_Type,
		Uses: composer.Type_Flow,
	}
}

const RenderValue_Type = "render_value"
const RenderValue_Field_Value = "$VALUE"

func (op *RenderValue) Marshal(m jsn.Marshaler) error {
	return RenderValue_Marshal(m, op)
}

type RenderValue_Slice []RenderValue

func (op *RenderValue_Slice) GetType() string { return RenderValue_Type }

func (op *RenderValue_Slice) Marshal(m jsn.Marshaler) error {
	return RenderValue_Repeats_Marshal(m, (*[]RenderValue)(op))
}

func (op *RenderValue_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *RenderValue_Slice) SetSize(cnt int) {
	var els []RenderValue
	if cnt >= 0 {
		els = make(RenderValue_Slice, cnt)
	}
	(*op) = els
}

func (op *RenderValue_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return RenderValue_Marshal(m, &(*op)[i])
}

func RenderValue_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderValue) error {
	return jsn.RepeatBlock(m, (*RenderValue_Slice)(vals))
}

func RenderValue_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]RenderValue) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = RenderValue_Repeats_Marshal(m, pv)
	}
	return
}

type RenderValue_Flow struct{ ptr *RenderValue }

func (n RenderValue_Flow) GetType() string      { return RenderValue_Type }
func (n RenderValue_Flow) GetLede() string      { return RenderValue_Type }
func (n RenderValue_Flow) GetFlow() interface{} { return n.ptr }
func (n RenderValue_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*RenderValue); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func RenderValue_Optional_Marshal(m jsn.Marshaler, pv **RenderValue) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RenderValue_Marshal(m, *pv)
	} else if !enc {
		var v RenderValue
		if err = RenderValue_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RenderValue_Marshal(m jsn.Marshaler, val *RenderValue) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(RenderValue_Flow{val}); err == nil {
		e0 := m.MarshalKey("", RenderValue_Field_Value)
		if e0 == nil {
			e0 = assign.Assignment_Marshal(m, &val.Value)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RenderValue_Field_Value))
		}
		m.EndBlock()
	}
	return
}

var Slots = []interface{}{
	(*RenderEval)(nil),
}

var Slats = []composer.Composer{
	(*RenderName)(nil),
	(*RenderPattern)(nil),
	(*RenderRef)(nil),
	(*RenderResponse)(nil),
	(*RenderValue)(nil),
}

var Signatures = map[uint64]interface{}{
	14401057669022842575: (*RenderPattern)(nil),  /* bool_eval=Render:render: */
	2910903954323771519:  (*RenderPattern)(nil),  /* render_eval=Render:render: */
	3385363614654173788:  (*RenderPattern)(nil),  /* text_eval=Render:render: */
	4328811686385928991:  (*RenderName)(nil),     /* text_eval=RenderName: */
	12372540113328333010: (*RenderRef)(nil),      /* bool_eval=RenderRef: */
	17707941731931999319: (*RenderRef)(nil),      /* num_list_eval=RenderRef: */
	586781755231363619:   (*RenderRef)(nil),      /* number_eval=RenderRef: */
	11952381947639314199: (*RenderRef)(nil),      /* record_eval=RenderRef: */
	5794615276964665178:  (*RenderRef)(nil),      /* record_list_eval=RenderRef: */
	15289959684061875714: (*RenderRef)(nil),      /* render_eval=RenderRef: */
	10542331033523904889: (*RenderRef)(nil),      /* text_eval=RenderRef: */
	4171261980310148416:  (*RenderRef)(nil),      /* text_list_eval=RenderRef: */
	18249933776929959289: (*RenderRef)(nil),      /* bool_eval=RenderRef:dot: */
	9735547470721472920:  (*RenderRef)(nil),      /* num_list_eval=RenderRef:dot: */
	13239953219501121612: (*RenderRef)(nil),      /* number_eval=RenderRef:dot: */
	8324158095841155032:  (*RenderRef)(nil),      /* record_eval=RenderRef:dot: */
	17618593433797581633: (*RenderRef)(nil),      /* record_list_eval=RenderRef:dot: */
	7883271647282708009:  (*RenderRef)(nil),      /* render_eval=RenderRef:dot: */
	239223853229152058:   (*RenderRef)(nil),      /* text_eval=RenderRef:dot: */
	3872622981826050135:  (*RenderRef)(nil),      /* text_list_eval=RenderRef:dot: */
	15658359855727638606: (*RenderResponse)(nil), /* execute=RenderResponse: */
	6351613444865908923:  (*RenderResponse)(nil), /* text_eval=RenderResponse: */
	167592851841791829:   (*RenderResponse)(nil), /* execute=RenderResponse:text: */
	10415880721138830946: (*RenderResponse)(nil), /* text_eval=RenderResponse:text: */
	7608693554121607902:  (*RenderValue)(nil),    /* render_eval=RenderValue: */
}
