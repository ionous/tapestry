// Code generated by "makeops"; edit at your own risk.
package debug

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/jsn"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// Comment Add a note.
// Information about the story for you and other authors.
type Comment struct {
	Lines  prim.Lines `if:"label=_"`
	Markup map[string]any
}

// User implemented slots:
var _ story.StoryStatement = (*Comment)(nil)
var _ rt.Execute = (*Comment)(nil)

func (*Comment) Compose() composer.Spec {
	return composer.Spec{
		Name: Comment_Type,
		Uses: composer.Type_Flow,
	}
}

const Comment_Type = "comment"
const Comment_Field_Lines = "$LINES"

func (op *Comment) Marshal(m jsn.Marshaler) error {
	return Comment_Marshal(m, op)
}

type Comment_Slice []Comment

func (op *Comment_Slice) GetType() string { return Comment_Type }

func (op *Comment_Slice) Marshal(m jsn.Marshaler) error {
	return Comment_Repeats_Marshal(m, (*[]Comment)(op))
}

func (op *Comment_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Comment_Slice) SetSize(cnt int) {
	var els []Comment
	if cnt >= 0 {
		els = make(Comment_Slice, cnt)
	}
	(*op) = els
}

func (op *Comment_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Comment_Marshal(m, &(*op)[i])
}

func Comment_Repeats_Marshal(m jsn.Marshaler, vals *[]Comment) error {
	return jsn.RepeatBlock(m, (*Comment_Slice)(vals))
}

func Comment_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Comment) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = Comment_Repeats_Marshal(m, pv)
	}
	return
}

type Comment_Flow struct{ ptr *Comment }

func (n Comment_Flow) GetType() string      { return Comment_Type }
func (n Comment_Flow) GetLede() string      { return Comment_Type }
func (n Comment_Flow) GetFlow() interface{} { return n.ptr }
func (n Comment_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Comment); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Comment_Optional_Marshal(m jsn.Marshaler, pv **Comment) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Comment_Marshal(m, *pv)
	} else if !enc {
		var v Comment
		if err = Comment_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Comment_Marshal(m jsn.Marshaler, val *Comment) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(Comment_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Comment_Field_Lines)
		if e0 == nil {
			e0 = prim.Lines_Marshal(m, &val.Lines)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Comment_Field_Lines))
		}
		m.EndBlock()
	}
	return
}

// DebugLog Debug log.
type DebugLog struct {
	Value    rt.Assignment `if:"label=_"`
	LogLevel LoggingLevel  `if:"label=as,optional"`
	Markup   map[string]any
}

// User implemented slots:
var _ rt.Execute = (*DebugLog)(nil)

func (*DebugLog) Compose() composer.Spec {
	return composer.Spec{
		Name: DebugLog_Type,
		Uses: composer.Type_Flow,
		Lede: "log",
	}
}

const DebugLog_Type = "debug_log"
const DebugLog_Field_Value = "$VALUE"
const DebugLog_Field_LogLevel = "$LOG_LEVEL"

func (op *DebugLog) Marshal(m jsn.Marshaler) error {
	return DebugLog_Marshal(m, op)
}

type DebugLog_Slice []DebugLog

func (op *DebugLog_Slice) GetType() string { return DebugLog_Type }

func (op *DebugLog_Slice) Marshal(m jsn.Marshaler) error {
	return DebugLog_Repeats_Marshal(m, (*[]DebugLog)(op))
}

func (op *DebugLog_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *DebugLog_Slice) SetSize(cnt int) {
	var els []DebugLog
	if cnt >= 0 {
		els = make(DebugLog_Slice, cnt)
	}
	(*op) = els
}

func (op *DebugLog_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return DebugLog_Marshal(m, &(*op)[i])
}

func DebugLog_Repeats_Marshal(m jsn.Marshaler, vals *[]DebugLog) error {
	return jsn.RepeatBlock(m, (*DebugLog_Slice)(vals))
}

func DebugLog_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]DebugLog) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = DebugLog_Repeats_Marshal(m, pv)
	}
	return
}

type DebugLog_Flow struct{ ptr *DebugLog }

func (n DebugLog_Flow) GetType() string      { return DebugLog_Type }
func (n DebugLog_Flow) GetLede() string      { return "log" }
func (n DebugLog_Flow) GetFlow() interface{} { return n.ptr }
func (n DebugLog_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*DebugLog); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func DebugLog_Optional_Marshal(m jsn.Marshaler, pv **DebugLog) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = DebugLog_Marshal(m, *pv)
	} else if !enc {
		var v DebugLog
		if err = DebugLog_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func DebugLog_Marshal(m jsn.Marshaler, val *DebugLog) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(DebugLog_Flow{val}); err == nil {
		e0 := m.MarshalKey("", DebugLog_Field_Value)
		if e0 == nil {
			e0 = rt.Assignment_Marshal(m, &val.Value)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", DebugLog_Field_Value))
		}
		e1 := m.MarshalKey("as", DebugLog_Field_LogLevel)
		if e1 == nil {
			e1 = LoggingLevel_Optional_Marshal(m, &val.LogLevel)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", DebugLog_Field_LogLevel))
		}
		m.EndBlock()
	}
	return
}

// DoNothing Statement which does nothing.
type DoNothing struct {
	Reason string `if:"label=why,optional,type=text"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.Execute = (*DoNothing)(nil)

func (*DoNothing) Compose() composer.Spec {
	return composer.Spec{
		Name: DoNothing_Type,
		Uses: composer.Type_Flow,
	}
}

const DoNothing_Type = "do_nothing"
const DoNothing_Field_Reason = "$REASON"

func (op *DoNothing) Marshal(m jsn.Marshaler) error {
	return DoNothing_Marshal(m, op)
}

type DoNothing_Slice []DoNothing

func (op *DoNothing_Slice) GetType() string { return DoNothing_Type }

func (op *DoNothing_Slice) Marshal(m jsn.Marshaler) error {
	return DoNothing_Repeats_Marshal(m, (*[]DoNothing)(op))
}

func (op *DoNothing_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *DoNothing_Slice) SetSize(cnt int) {
	var els []DoNothing
	if cnt >= 0 {
		els = make(DoNothing_Slice, cnt)
	}
	(*op) = els
}

func (op *DoNothing_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return DoNothing_Marshal(m, &(*op)[i])
}

func DoNothing_Repeats_Marshal(m jsn.Marshaler, vals *[]DoNothing) error {
	return jsn.RepeatBlock(m, (*DoNothing_Slice)(vals))
}

func DoNothing_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]DoNothing) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = DoNothing_Repeats_Marshal(m, pv)
	}
	return
}

type DoNothing_Flow struct{ ptr *DoNothing }

func (n DoNothing_Flow) GetType() string      { return DoNothing_Type }
func (n DoNothing_Flow) GetLede() string      { return DoNothing_Type }
func (n DoNothing_Flow) GetFlow() interface{} { return n.ptr }
func (n DoNothing_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*DoNothing); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func DoNothing_Optional_Marshal(m jsn.Marshaler, pv **DoNothing) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = DoNothing_Marshal(m, *pv)
	} else if !enc {
		var v DoNothing
		if err = DoNothing_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func DoNothing_Marshal(m jsn.Marshaler, val *DoNothing) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(DoNothing_Flow{val}); err == nil {
		e0 := m.MarshalKey("why", DoNothing_Field_Reason)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Optional_Marshal(m, &val.Reason)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", DoNothing_Field_Reason))
		}
		m.EndBlock()
	}
	return
}

// ExpectBool
type ExpectBool struct {
	Value  rt.BoolEval `if:"label=_"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.Execute = (*ExpectBool)(nil)

func (*ExpectBool) Compose() composer.Spec {
	return composer.Spec{
		Name: ExpectBool_Type,
		Uses: composer.Type_Flow,
		Lede: "expect",
	}
}

const ExpectBool_Type = "expect_bool"
const ExpectBool_Field_Value = "$VALUE"

func (op *ExpectBool) Marshal(m jsn.Marshaler) error {
	return ExpectBool_Marshal(m, op)
}

type ExpectBool_Slice []ExpectBool

func (op *ExpectBool_Slice) GetType() string { return ExpectBool_Type }

func (op *ExpectBool_Slice) Marshal(m jsn.Marshaler) error {
	return ExpectBool_Repeats_Marshal(m, (*[]ExpectBool)(op))
}

func (op *ExpectBool_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *ExpectBool_Slice) SetSize(cnt int) {
	var els []ExpectBool
	if cnt >= 0 {
		els = make(ExpectBool_Slice, cnt)
	}
	(*op) = els
}

func (op *ExpectBool_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return ExpectBool_Marshal(m, &(*op)[i])
}

func ExpectBool_Repeats_Marshal(m jsn.Marshaler, vals *[]ExpectBool) error {
	return jsn.RepeatBlock(m, (*ExpectBool_Slice)(vals))
}

func ExpectBool_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]ExpectBool) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = ExpectBool_Repeats_Marshal(m, pv)
	}
	return
}

type ExpectBool_Flow struct{ ptr *ExpectBool }

func (n ExpectBool_Flow) GetType() string      { return ExpectBool_Type }
func (n ExpectBool_Flow) GetLede() string      { return "expect" }
func (n ExpectBool_Flow) GetFlow() interface{} { return n.ptr }
func (n ExpectBool_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*ExpectBool); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func ExpectBool_Optional_Marshal(m jsn.Marshaler, pv **ExpectBool) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = ExpectBool_Marshal(m, *pv)
	} else if !enc {
		var v ExpectBool
		if err = ExpectBool_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func ExpectBool_Marshal(m jsn.Marshaler, val *ExpectBool) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(ExpectBool_Flow{val}); err == nil {
		e0 := m.MarshalKey("", ExpectBool_Field_Value)
		if e0 == nil {
			e0 = rt.BoolEval_Marshal(m, &val.Value)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", ExpectBool_Field_Value))
		}
		m.EndBlock()
	}
	return
}

// ExpectLines
type ExpectLines struct {
	Lines  prim.Lines `if:"label=lines"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.Execute = (*ExpectLines)(nil)

func (*ExpectLines) Compose() composer.Spec {
	return composer.Spec{
		Name: ExpectLines_Type,
		Uses: composer.Type_Flow,
		Lede: "expect",
	}
}

const ExpectLines_Type = "expect_lines"
const ExpectLines_Field_Lines = "$LINES"

func (op *ExpectLines) Marshal(m jsn.Marshaler) error {
	return ExpectLines_Marshal(m, op)
}

type ExpectLines_Slice []ExpectLines

func (op *ExpectLines_Slice) GetType() string { return ExpectLines_Type }

func (op *ExpectLines_Slice) Marshal(m jsn.Marshaler) error {
	return ExpectLines_Repeats_Marshal(m, (*[]ExpectLines)(op))
}

func (op *ExpectLines_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *ExpectLines_Slice) SetSize(cnt int) {
	var els []ExpectLines
	if cnt >= 0 {
		els = make(ExpectLines_Slice, cnt)
	}
	(*op) = els
}

func (op *ExpectLines_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return ExpectLines_Marshal(m, &(*op)[i])
}

func ExpectLines_Repeats_Marshal(m jsn.Marshaler, vals *[]ExpectLines) error {
	return jsn.RepeatBlock(m, (*ExpectLines_Slice)(vals))
}

func ExpectLines_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]ExpectLines) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = ExpectLines_Repeats_Marshal(m, pv)
	}
	return
}

type ExpectLines_Flow struct{ ptr *ExpectLines }

func (n ExpectLines_Flow) GetType() string      { return ExpectLines_Type }
func (n ExpectLines_Flow) GetLede() string      { return "expect" }
func (n ExpectLines_Flow) GetFlow() interface{} { return n.ptr }
func (n ExpectLines_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*ExpectLines); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func ExpectLines_Optional_Marshal(m jsn.Marshaler, pv **ExpectLines) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = ExpectLines_Marshal(m, *pv)
	} else if !enc {
		var v ExpectLines
		if err = ExpectLines_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func ExpectLines_Marshal(m jsn.Marshaler, val *ExpectLines) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(ExpectLines_Flow{val}); err == nil {
		e0 := m.MarshalKey("lines", ExpectLines_Field_Lines)
		if e0 == nil {
			e0 = prim.Lines_Marshal(m, &val.Lines)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", ExpectLines_Field_Lines))
		}
		m.EndBlock()
	}
	return
}

// ExpectNum
type ExpectNum struct {
	Result    float64         `if:"label=_,type=number"`
	Is        core.Comparison `if:"label=is"`
	Value     rt.NumberEval   `if:"label=num"`
	Tolerance float64         `if:"label=within,optional,type=number"`
	Markup    map[string]any
}

// User implemented slots:
var _ rt.Execute = (*ExpectNum)(nil)

func (*ExpectNum) Compose() composer.Spec {
	return composer.Spec{
		Name: ExpectNum_Type,
		Uses: composer.Type_Flow,
		Lede: "expect",
	}
}

const ExpectNum_Type = "expect_num"
const ExpectNum_Field_Result = "$RESULT"
const ExpectNum_Field_Is = "$IS"
const ExpectNum_Field_Value = "$VALUE"
const ExpectNum_Field_Tolerance = "$TOLERANCE"

func (op *ExpectNum) Marshal(m jsn.Marshaler) error {
	return ExpectNum_Marshal(m, op)
}

type ExpectNum_Slice []ExpectNum

func (op *ExpectNum_Slice) GetType() string { return ExpectNum_Type }

func (op *ExpectNum_Slice) Marshal(m jsn.Marshaler) error {
	return ExpectNum_Repeats_Marshal(m, (*[]ExpectNum)(op))
}

func (op *ExpectNum_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *ExpectNum_Slice) SetSize(cnt int) {
	var els []ExpectNum
	if cnt >= 0 {
		els = make(ExpectNum_Slice, cnt)
	}
	(*op) = els
}

func (op *ExpectNum_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return ExpectNum_Marshal(m, &(*op)[i])
}

func ExpectNum_Repeats_Marshal(m jsn.Marshaler, vals *[]ExpectNum) error {
	return jsn.RepeatBlock(m, (*ExpectNum_Slice)(vals))
}

func ExpectNum_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]ExpectNum) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = ExpectNum_Repeats_Marshal(m, pv)
	}
	return
}

type ExpectNum_Flow struct{ ptr *ExpectNum }

func (n ExpectNum_Flow) GetType() string      { return ExpectNum_Type }
func (n ExpectNum_Flow) GetLede() string      { return "expect" }
func (n ExpectNum_Flow) GetFlow() interface{} { return n.ptr }
func (n ExpectNum_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*ExpectNum); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func ExpectNum_Optional_Marshal(m jsn.Marshaler, pv **ExpectNum) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = ExpectNum_Marshal(m, *pv)
	} else if !enc {
		var v ExpectNum
		if err = ExpectNum_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func ExpectNum_Marshal(m jsn.Marshaler, val *ExpectNum) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(ExpectNum_Flow{val}); err == nil {
		e0 := m.MarshalKey("", ExpectNum_Field_Result)
		if e0 == nil {
			e0 = prim.Number_Unboxed_Marshal(m, &val.Result)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", ExpectNum_Field_Result))
		}
		e1 := m.MarshalKey("is", ExpectNum_Field_Is)
		if e1 == nil {
			e1 = core.Comparison_Marshal(m, &val.Is)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", ExpectNum_Field_Is))
		}
		e2 := m.MarshalKey("num", ExpectNum_Field_Value)
		if e2 == nil {
			e2 = rt.NumberEval_Marshal(m, &val.Value)
		}
		if e2 != nil && e2 != jsn.Missing {
			m.Error(errutil.New(e2, "in flow at", ExpectNum_Field_Value))
		}
		e3 := m.MarshalKey("within", ExpectNum_Field_Tolerance)
		if e3 == nil {
			e3 = prim.Number_Unboxed_Optional_Marshal(m, &val.Tolerance)
		}
		if e3 != nil && e3 != jsn.Missing {
			m.Error(errutil.New(e3, "in flow at", ExpectNum_Field_Tolerance))
		}
		m.EndBlock()
	}
	return
}

// ExpectString
type ExpectString struct {
	String string `if:"label=string,type=text"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.Execute = (*ExpectString)(nil)

func (*ExpectString) Compose() composer.Spec {
	return composer.Spec{
		Name: ExpectString_Type,
		Uses: composer.Type_Flow,
		Lede: "expect",
	}
}

const ExpectString_Type = "expect_string"
const ExpectString_Field_String = "$STRING"

func (op *ExpectString) Marshal(m jsn.Marshaler) error {
	return ExpectString_Marshal(m, op)
}

type ExpectString_Slice []ExpectString

func (op *ExpectString_Slice) GetType() string { return ExpectString_Type }

func (op *ExpectString_Slice) Marshal(m jsn.Marshaler) error {
	return ExpectString_Repeats_Marshal(m, (*[]ExpectString)(op))
}

func (op *ExpectString_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *ExpectString_Slice) SetSize(cnt int) {
	var els []ExpectString
	if cnt >= 0 {
		els = make(ExpectString_Slice, cnt)
	}
	(*op) = els
}

func (op *ExpectString_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return ExpectString_Marshal(m, &(*op)[i])
}

func ExpectString_Repeats_Marshal(m jsn.Marshaler, vals *[]ExpectString) error {
	return jsn.RepeatBlock(m, (*ExpectString_Slice)(vals))
}

func ExpectString_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]ExpectString) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = ExpectString_Repeats_Marshal(m, pv)
	}
	return
}

type ExpectString_Flow struct{ ptr *ExpectString }

func (n ExpectString_Flow) GetType() string      { return ExpectString_Type }
func (n ExpectString_Flow) GetLede() string      { return "expect" }
func (n ExpectString_Flow) GetFlow() interface{} { return n.ptr }
func (n ExpectString_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*ExpectString); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func ExpectString_Optional_Marshal(m jsn.Marshaler, pv **ExpectString) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = ExpectString_Marshal(m, *pv)
	} else if !enc {
		var v ExpectString
		if err = ExpectString_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func ExpectString_Marshal(m jsn.Marshaler, val *ExpectString) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(ExpectString_Flow{val}); err == nil {
		e0 := m.MarshalKey("string", ExpectString_Field_String)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Marshal(m, &val.String)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", ExpectString_Field_String))
		}
		m.EndBlock()
	}
	return
}

// ExpectText
type ExpectText struct {
	Text   rt.TextEval `if:"label=text"`
	Markup map[string]any
}

// User implemented slots:
var _ rt.Execute = (*ExpectText)(nil)

func (*ExpectText) Compose() composer.Spec {
	return composer.Spec{
		Name: ExpectText_Type,
		Uses: composer.Type_Flow,
		Lede: "expect",
	}
}

const ExpectText_Type = "expect_text"
const ExpectText_Field_Text = "$TEXT"

func (op *ExpectText) Marshal(m jsn.Marshaler) error {
	return ExpectText_Marshal(m, op)
}

type ExpectText_Slice []ExpectText

func (op *ExpectText_Slice) GetType() string { return ExpectText_Type }

func (op *ExpectText_Slice) Marshal(m jsn.Marshaler) error {
	return ExpectText_Repeats_Marshal(m, (*[]ExpectText)(op))
}

func (op *ExpectText_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *ExpectText_Slice) SetSize(cnt int) {
	var els []ExpectText
	if cnt >= 0 {
		els = make(ExpectText_Slice, cnt)
	}
	(*op) = els
}

func (op *ExpectText_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return ExpectText_Marshal(m, &(*op)[i])
}

func ExpectText_Repeats_Marshal(m jsn.Marshaler, vals *[]ExpectText) error {
	return jsn.RepeatBlock(m, (*ExpectText_Slice)(vals))
}

func ExpectText_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]ExpectText) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = ExpectText_Repeats_Marshal(m, pv)
	}
	return
}

type ExpectText_Flow struct{ ptr *ExpectText }

func (n ExpectText_Flow) GetType() string      { return ExpectText_Type }
func (n ExpectText_Flow) GetLede() string      { return "expect" }
func (n ExpectText_Flow) GetFlow() interface{} { return n.ptr }
func (n ExpectText_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*ExpectText); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func ExpectText_Optional_Marshal(m jsn.Marshaler, pv **ExpectText) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = ExpectText_Marshal(m, *pv)
	} else if !enc {
		var v ExpectText
		if err = ExpectText_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func ExpectText_Marshal(m jsn.Marshaler, val *ExpectText) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(ExpectText_Flow{val}); err == nil {
		e0 := m.MarshalKey("text", ExpectText_Field_Text)
		if e0 == nil {
			e0 = rt.TextEval_Marshal(m, &val.Text)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", ExpectText_Field_Text))
		}
		m.EndBlock()
	}
	return
}

// LoggingLevel requires a predefined string.
type LoggingLevel struct {
	Str string
}

func (op *LoggingLevel) String() string {
	return op.Str
}

const LoggingLevel_Note = "$NOTE"
const LoggingLevel_ToDo = "$TO_DO"
const LoggingLevel_Fix = "$FIX"
const LoggingLevel_Info = "$INFO"
const LoggingLevel_Warning = "$WARNING"
const LoggingLevel_Error = "$ERROR"

func (*LoggingLevel) Compose() composer.Spec {
	return composer.Spec{
		Name: LoggingLevel_Type,
		Uses: composer.Type_Str,
		Choices: []string{
			LoggingLevel_Note, LoggingLevel_ToDo, LoggingLevel_Fix, LoggingLevel_Info, LoggingLevel_Warning, LoggingLevel_Error,
		},
		Strings: []string{
			"note", "to_do", "fix", "info", "warning", "error",
		},
	}
}

const LoggingLevel_Type = "logging_level"

func (op *LoggingLevel) Marshal(m jsn.Marshaler) error {
	return LoggingLevel_Marshal(m, op)
}

func LoggingLevel_Optional_Marshal(m jsn.Marshaler, val *LoggingLevel) (err error) {
	var zero LoggingLevel
	if enc := m.IsEncoding(); !enc || val.Str != zero.Str {
		err = LoggingLevel_Marshal(m, val)
	}
	return
}

func LoggingLevel_Marshal(m jsn.Marshaler, val *LoggingLevel) (err error) {
	return m.MarshalValue(LoggingLevel_Type, jsn.MakeEnum(val, &val.Str))
}

type LoggingLevel_Slice []LoggingLevel

func (op *LoggingLevel_Slice) GetType() string { return LoggingLevel_Type }

func (op *LoggingLevel_Slice) Marshal(m jsn.Marshaler) error {
	return LoggingLevel_Repeats_Marshal(m, (*[]LoggingLevel)(op))
}

func (op *LoggingLevel_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *LoggingLevel_Slice) SetSize(cnt int) {
	var els []LoggingLevel
	if cnt >= 0 {
		els = make(LoggingLevel_Slice, cnt)
	}
	(*op) = els
}

func (op *LoggingLevel_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return LoggingLevel_Marshal(m, &(*op)[i])
}

func LoggingLevel_Repeats_Marshal(m jsn.Marshaler, vals *[]LoggingLevel) error {
	return jsn.RepeatBlock(m, (*LoggingLevel_Slice)(vals))
}

func LoggingLevel_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]LoggingLevel) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = LoggingLevel_Repeats_Marshal(m, pv)
	}
	return
}

var Slats = []composer.Composer{
	(*Comment)(nil),
	(*DebugLog)(nil),
	(*DoNothing)(nil),
	(*ExpectBool)(nil),
	(*ExpectLines)(nil),
	(*ExpectNum)(nil),
	(*ExpectString)(nil),
	(*ExpectText)(nil),
	(*LoggingLevel)(nil),
}

var Signatures = map[uint64]interface{}{
	15823738440204397330: (*LoggingLevel)(nil), /* LoggingLevel: */
	3991849378064754806:  (*Comment)(nil),      /* execute=Comment: */
	16586092333187989882: (*Comment)(nil),      /* story_statement=Comment: */
	14645287343365598707: (*DoNothing)(nil),    /* execute=DoNothing */
	12243119421914882789: (*DoNothing)(nil),    /* execute=DoNothing why: */
	469594313115947985:   (*ExpectLines)(nil),  /* execute=Expect lines: */
	5505041336569015051:  (*ExpectString)(nil), /* execute=Expect string: */
	16489874106085927697: (*ExpectText)(nil),   /* execute=Expect text: */
	11108202414968227788: (*ExpectBool)(nil),   /* execute=Expect: */
	9770230868586544920:  (*ExpectNum)(nil),    /* execute=Expect:is:num: */
	8339796867902453679:  (*ExpectNum)(nil),    /* execute=Expect:is:num:within: */
	17230987244745403983: (*DebugLog)(nil),     /* execute=Log: */
	9146550673186999987:  (*DebugLog)(nil),     /* execute=Log:as: */
}
