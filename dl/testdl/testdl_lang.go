// Code generated by "makeops"; edit at your own risk.
package testdl

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/jsn"
	"github.com/ionous/errutil"
)

// TestBool requires a predefined string.
type TestBool struct {
	Str string
}

func (op *TestBool) String() string {
	return op.Str
}

const TestBool_True = "$TRUE"
const TestBool_False = "$FALSE"

func (*TestBool) Compose() composer.Spec {
	return composer.Spec{
		Name: TestBool_Type,
		Uses: composer.Type_Str,
		Choices: []string{
			TestBool_True, TestBool_False,
		},
		Strings: []string{
			"true", "false",
		},
	}
}

const TestBool_Type = "test_bool"

func (op *TestBool) Marshal(m jsn.Marshaler) error {
	return TestBool_Marshal(m, op)
}

func TestBool_Optional_Marshal(m jsn.Marshaler, val *TestBool) (err error) {
	var zero TestBool
	if enc := m.IsEncoding(); !enc || val.Str != zero.Str {
		err = TestBool_Marshal(m, val)
	}
	return
}

func TestBool_Marshal(m jsn.Marshaler, val *TestBool) (err error) {
	return m.MarshalValue(TestBool_Type, jsn.MakeEnum(val, &val.Str))
}

type TestBool_Slice []TestBool

func (op *TestBool_Slice) GetType() string { return TestBool_Type }

func (op *TestBool_Slice) Marshal(m jsn.Marshaler) error {
	return TestBool_Repeats_Marshal(m, (*[]TestBool)(op))
}

func (op *TestBool_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TestBool_Slice) SetSize(cnt int) {
	var els []TestBool
	if cnt >= 0 {
		els = make(TestBool_Slice, cnt)
	}
	(*op) = els
}

func (op *TestBool_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TestBool_Marshal(m, &(*op)[i])
}

func TestBool_Repeats_Marshal(m jsn.Marshaler, vals *[]TestBool) error {
	return jsn.RepeatBlock(m, (*TestBool_Slice)(vals))
}

func TestBool_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TestBool) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TestBool_Repeats_Marshal(m, pv)
	}
	return
}

// TestEmbed
type TestEmbed struct {
	TestFlow TestFlow `if:"label=test_flow"`
	Markup   map[string]any
}

// User implemented slots:
var _ TestSlot = (*TestEmbed)(nil)

func (*TestEmbed) Compose() composer.Spec {
	return composer.Spec{
		Name: TestEmbed_Type,
		Uses: composer.Type_Flow,
		Lede: "embed",
	}
}

const TestEmbed_Type = "test_embed"
const TestEmbed_Field_TestFlow = "$TEST_FLOW"

func (op *TestEmbed) Marshal(m jsn.Marshaler) error {
	return TestEmbed_Marshal(m, op)
}

type TestEmbed_Slice []TestEmbed

func (op *TestEmbed_Slice) GetType() string { return TestEmbed_Type }

func (op *TestEmbed_Slice) Marshal(m jsn.Marshaler) error {
	return TestEmbed_Repeats_Marshal(m, (*[]TestEmbed)(op))
}

func (op *TestEmbed_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TestEmbed_Slice) SetSize(cnt int) {
	var els []TestEmbed
	if cnt >= 0 {
		els = make(TestEmbed_Slice, cnt)
	}
	(*op) = els
}

func (op *TestEmbed_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TestEmbed_Marshal(m, &(*op)[i])
}

func TestEmbed_Repeats_Marshal(m jsn.Marshaler, vals *[]TestEmbed) error {
	return jsn.RepeatBlock(m, (*TestEmbed_Slice)(vals))
}

func TestEmbed_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TestEmbed) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TestEmbed_Repeats_Marshal(m, pv)
	}
	return
}

type TestEmbed_Flow struct{ ptr *TestEmbed }

func (n TestEmbed_Flow) GetType() string      { return TestEmbed_Type }
func (n TestEmbed_Flow) GetLede() string      { return "embed" }
func (n TestEmbed_Flow) GetFlow() interface{} { return n.ptr }
func (n TestEmbed_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*TestEmbed); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func TestEmbed_Optional_Marshal(m jsn.Marshaler, pv **TestEmbed) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = TestEmbed_Marshal(m, *pv)
	} else if !enc {
		var v TestEmbed
		if err = TestEmbed_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func TestEmbed_Marshal(m jsn.Marshaler, val *TestEmbed) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(TestEmbed_Flow{val}); err == nil {
		e0 := m.MarshalKey("test_flow", TestEmbed_Field_TestFlow)
		if e0 == nil {
			e0 = TestFlow_Marshal(m, &val.TestFlow)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", TestEmbed_Field_TestFlow))
		}
		m.EndBlock()
	}
	return
}

// TestFlow
type TestFlow struct {
	Slot   TestSlot   `if:"label=slot,optional"`
	Txt    TestTxt    `if:"label=txt,optional"`
	Num    float64    `if:"label=num,optional,type=test_num"`
	Bool   TestBool   `if:"label=bool,optional"`
	Swap   TestSwap   `if:"label=swap,optional"`
	Slots  []TestSlot `if:"label=slots,optional"`
	Markup map[string]any
}

// User implemented slots:
var _ TestSlot = (*TestFlow)(nil)

func (*TestFlow) Compose() composer.Spec {
	return composer.Spec{
		Name: TestFlow_Type,
		Uses: composer.Type_Flow,
		Lede: "flow",
	}
}

const TestFlow_Type = "test_flow"
const TestFlow_Field_Slot = "$SLOT"
const TestFlow_Field_Txt = "$TXT"
const TestFlow_Field_Num = "$NUM"
const TestFlow_Field_Bool = "$BOOL"
const TestFlow_Field_Swap = "$SWAP"
const TestFlow_Field_Slots = "$SLOTS"

func (op *TestFlow) Marshal(m jsn.Marshaler) error {
	return TestFlow_Marshal(m, op)
}

type TestFlow_Slice []TestFlow

func (op *TestFlow_Slice) GetType() string { return TestFlow_Type }

func (op *TestFlow_Slice) Marshal(m jsn.Marshaler) error {
	return TestFlow_Repeats_Marshal(m, (*[]TestFlow)(op))
}

func (op *TestFlow_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TestFlow_Slice) SetSize(cnt int) {
	var els []TestFlow
	if cnt >= 0 {
		els = make(TestFlow_Slice, cnt)
	}
	(*op) = els
}

func (op *TestFlow_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TestFlow_Marshal(m, &(*op)[i])
}

func TestFlow_Repeats_Marshal(m jsn.Marshaler, vals *[]TestFlow) error {
	return jsn.RepeatBlock(m, (*TestFlow_Slice)(vals))
}

func TestFlow_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TestFlow) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TestFlow_Repeats_Marshal(m, pv)
	}
	return
}

type TestFlow_Flow struct{ ptr *TestFlow }

func (n TestFlow_Flow) GetType() string      { return TestFlow_Type }
func (n TestFlow_Flow) GetLede() string      { return "flow" }
func (n TestFlow_Flow) GetFlow() interface{} { return n.ptr }
func (n TestFlow_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*TestFlow); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func TestFlow_Optional_Marshal(m jsn.Marshaler, pv **TestFlow) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = TestFlow_Marshal(m, *pv)
	} else if !enc {
		var v TestFlow
		if err = TestFlow_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func TestFlow_Marshal(m jsn.Marshaler, val *TestFlow) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(TestFlow_Flow{val}); err == nil {
		e0 := m.MarshalKey("slot", TestFlow_Field_Slot)
		if e0 == nil {
			e0 = TestSlot_Optional_Marshal(m, &val.Slot)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", TestFlow_Field_Slot))
		}
		e1 := m.MarshalKey("txt", TestFlow_Field_Txt)
		if e1 == nil {
			e1 = TestTxt_Optional_Marshal(m, &val.Txt)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", TestFlow_Field_Txt))
		}
		e2 := m.MarshalKey("num", TestFlow_Field_Num)
		if e2 == nil {
			e2 = TestNum_Unboxed_Optional_Marshal(m, &val.Num)
		}
		if e2 != nil && e2 != jsn.Missing {
			m.Error(errutil.New(e2, "in flow at", TestFlow_Field_Num))
		}
		e3 := m.MarshalKey("bool", TestFlow_Field_Bool)
		if e3 == nil {
			e3 = TestBool_Optional_Marshal(m, &val.Bool)
		}
		if e3 != nil && e3 != jsn.Missing {
			m.Error(errutil.New(e3, "in flow at", TestFlow_Field_Bool))
		}
		e4 := m.MarshalKey("swap", TestFlow_Field_Swap)
		if e4 == nil {
			e4 = TestSwap_Optional_Marshal(m, &val.Swap)
		}
		if e4 != nil && e4 != jsn.Missing {
			m.Error(errutil.New(e4, "in flow at", TestFlow_Field_Swap))
		}
		e5 := m.MarshalKey("slots", TestFlow_Field_Slots)
		if e5 == nil {
			e5 = TestSlot_Optional_Repeats_Marshal(m, &val.Slots)
		}
		if e5 != nil && e5 != jsn.Missing {
			m.Error(errutil.New(e5, "in flow at", TestFlow_Field_Slots))
		}
		m.EndBlock()
	}
	return
}

// TestNum requires a user-specified number.
type TestNum struct {
	Num float64
}

func (*TestNum) Compose() composer.Spec {
	return composer.Spec{
		Name: TestNum_Type,
		Uses: composer.Type_Num,
	}
}

const TestNum_Type = "test_num"

func (op *TestNum) Marshal(m jsn.Marshaler) error {
	return TestNum_Marshal(m, op)
}

type TestNum_Unboxed_Slice []float64

func (op *TestNum_Unboxed_Slice) GetType() string { return TestNum_Type }

func (op *TestNum_Unboxed_Slice) Marshal(m jsn.Marshaler) error {
	return TestNum_Unboxed_Repeats_Marshal(m, (*[]float64)(op))
}

func (op *TestNum_Unboxed_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TestNum_Unboxed_Slice) SetSize(cnt int) {
	var els []float64
	if cnt >= 0 {
		els = make(TestNum_Unboxed_Slice, cnt)
	}
	(*op) = els
}

func (op *TestNum_Unboxed_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TestNum_Unboxed_Marshal(m, &(*op)[i])
}

func TestNum_Unboxed_Repeats_Marshal(m jsn.Marshaler, vals *[]float64) error {
	return jsn.RepeatBlock(m, (*TestNum_Unboxed_Slice)(vals))
}

func TestNum_Unboxed_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]float64) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TestNum_Unboxed_Repeats_Marshal(m, pv)
	}
	return
}

func TestNum_Unboxed_Optional_Marshal(m jsn.Marshaler, val *float64) (err error) {
	var zero float64
	if enc := m.IsEncoding(); !enc || *val != zero {
		err = TestNum_Unboxed_Marshal(m, val)
	}
	return
}

func TestNum_Unboxed_Marshal(m jsn.Marshaler, val *float64) error {
	return m.MarshalValue(TestNum_Type, jsn.BoxFloat64(val))
}

func TestNum_Optional_Marshal(m jsn.Marshaler, val *TestNum) (err error) {
	var zero TestNum
	if enc := m.IsEncoding(); !enc || val.Num != zero.Num {
		err = TestNum_Marshal(m, val)
	}
	return
}

func TestNum_Marshal(m jsn.Marshaler, val *TestNum) (err error) {
	return m.MarshalValue(TestNum_Type, &val.Num)
}

type TestNum_Slice []TestNum

func (op *TestNum_Slice) GetType() string { return TestNum_Type }

func (op *TestNum_Slice) Marshal(m jsn.Marshaler) error {
	return TestNum_Repeats_Marshal(m, (*[]TestNum)(op))
}

func (op *TestNum_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TestNum_Slice) SetSize(cnt int) {
	var els []TestNum
	if cnt >= 0 {
		els = make(TestNum_Slice, cnt)
	}
	(*op) = els
}

func (op *TestNum_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TestNum_Marshal(m, &(*op)[i])
}

func TestNum_Repeats_Marshal(m jsn.Marshaler, vals *[]TestNum) error {
	return jsn.RepeatBlock(m, (*TestNum_Slice)(vals))
}

func TestNum_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TestNum) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TestNum_Repeats_Marshal(m, pv)
	}
	return
}

const TestSlot_Type = "test_slot"

var TestSlot_Optional_Marshal = TestSlot_Marshal

type TestSlot_Slot struct{ Value *TestSlot }

func (at TestSlot_Slot) Marshal(m jsn.Marshaler) (err error) {
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
func (at TestSlot_Slot) GetType() string              { return TestSlot_Type }
func (at TestSlot_Slot) GetSlot() (interface{}, bool) { return *at.Value, *at.Value != nil }
func (at TestSlot_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.Value), okay = v.(TestSlot)
	return
}

func TestSlot_Marshal(m jsn.Marshaler, ptr *TestSlot) (err error) {
	slot := TestSlot_Slot{ptr}
	return slot.Marshal(m)
}

type TestSlot_Slice []TestSlot

func (op *TestSlot_Slice) GetType() string { return TestSlot_Type }

func (op *TestSlot_Slice) Marshal(m jsn.Marshaler) error {
	return TestSlot_Repeats_Marshal(m, (*[]TestSlot)(op))
}

func (op *TestSlot_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TestSlot_Slice) SetSize(cnt int) {
	var els []TestSlot
	if cnt >= 0 {
		els = make(TestSlot_Slice, cnt)
	}
	(*op) = els
}

func (op *TestSlot_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TestSlot_Marshal(m, &(*op)[i])
}

func TestSlot_Repeats_Marshal(m jsn.Marshaler, vals *[]TestSlot) error {
	return jsn.RepeatBlock(m, (*TestSlot_Slice)(vals))
}

func TestSlot_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TestSlot) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TestSlot_Repeats_Marshal(m, pv)
	}
	return
}

// TestStr requires a predefined string.
type TestStr struct {
	Str string
}

func (op *TestStr) String() string {
	return op.Str
}

const TestStr_One = "$ONE"
const TestStr_Other = "$OTHER"
const TestStr_Option = "$OPTION"

func (*TestStr) Compose() composer.Spec {
	return composer.Spec{
		Name: TestStr_Type,
		Uses: composer.Type_Str,
		Choices: []string{
			TestStr_One, TestStr_Other, TestStr_Option,
		},
		Strings: []string{
			"one", "other", "option",
		},
	}
}

const TestStr_Type = "test_str"

func (op *TestStr) Marshal(m jsn.Marshaler) error {
	return TestStr_Marshal(m, op)
}

func TestStr_Optional_Marshal(m jsn.Marshaler, val *TestStr) (err error) {
	var zero TestStr
	if enc := m.IsEncoding(); !enc || val.Str != zero.Str {
		err = TestStr_Marshal(m, val)
	}
	return
}

func TestStr_Marshal(m jsn.Marshaler, val *TestStr) (err error) {
	return m.MarshalValue(TestStr_Type, jsn.MakeEnum(val, &val.Str))
}

type TestStr_Slice []TestStr

func (op *TestStr_Slice) GetType() string { return TestStr_Type }

func (op *TestStr_Slice) Marshal(m jsn.Marshaler) error {
	return TestStr_Repeats_Marshal(m, (*[]TestStr)(op))
}

func (op *TestStr_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TestStr_Slice) SetSize(cnt int) {
	var els []TestStr
	if cnt >= 0 {
		els = make(TestStr_Slice, cnt)
	}
	(*op) = els
}

func (op *TestStr_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TestStr_Marshal(m, &(*op)[i])
}

func TestStr_Repeats_Marshal(m jsn.Marshaler, vals *[]TestStr) error {
	return jsn.RepeatBlock(m, (*TestStr_Slice)(vals))
}

func TestStr_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TestStr) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TestStr_Repeats_Marshal(m, pv)
	}
	return
}

// TestSwap swaps between various options
type TestSwap struct {
	Choice string
	Value  interface{}
}

var TestSwap_Optional_Marshal = TestSwap_Marshal

const TestSwap_A_Opt = "$A"
const TestSwap_B_Opt = "$B"
const TestSwap_C_Opt = "$C"

func (*TestSwap) Compose() composer.Spec {
	return composer.Spec{
		Name: TestSwap_Type,
		Uses: composer.Type_Swap,
		Choices: []string{
			TestSwap_A_Opt, TestSwap_B_Opt, TestSwap_C_Opt,
		},
		Swaps: []interface{}{
			(*TestFlow)(nil),
			(*TestSlot)(nil),
			(*TestTxt)(nil),
		},
	}
}

const TestSwap_Type = "test_swap"

func (op *TestSwap) GetType() string { return TestSwap_Type }

func (op *TestSwap) GetSwap() (string, interface{}) {
	return op.Choice, op.Value
}

func (op *TestSwap) SetSwap(c string) (okay bool) {
	switch c {
	case "":
		op.Choice, op.Value = c, nil
		okay = true
	case TestSwap_A_Opt:
		op.Choice, op.Value = c, new(TestFlow)
		okay = true
	case TestSwap_B_Opt:
		op.Choice, op.Value = c, new(TestSlot)
		okay = true
	case TestSwap_C_Opt:
		op.Choice, op.Value = c, new(TestTxt)
		okay = true
	}
	return
}

func (op *TestSwap) Marshal(m jsn.Marshaler) error {
	return TestSwap_Marshal(m, op)
}
func TestSwap_Marshal(m jsn.Marshaler, val *TestSwap) (err error) {
	if err = m.MarshalBlock(val); err == nil {
		if _, ptr := val.GetSwap(); ptr != nil {
			if e := ptr.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}

type TestSwap_Slice []TestSwap

func (op *TestSwap_Slice) GetType() string { return TestSwap_Type }

func (op *TestSwap_Slice) Marshal(m jsn.Marshaler) error {
	return TestSwap_Repeats_Marshal(m, (*[]TestSwap)(op))
}

func (op *TestSwap_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TestSwap_Slice) SetSize(cnt int) {
	var els []TestSwap
	if cnt >= 0 {
		els = make(TestSwap_Slice, cnt)
	}
	(*op) = els
}

func (op *TestSwap_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TestSwap_Marshal(m, &(*op)[i])
}

func TestSwap_Repeats_Marshal(m jsn.Marshaler, vals *[]TestSwap) error {
	return jsn.RepeatBlock(m, (*TestSwap_Slice)(vals))
}

func TestSwap_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TestSwap) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TestSwap_Repeats_Marshal(m, pv)
	}
	return
}

// TestTxt requires a user-specified string.
type TestTxt struct {
	Str string
}

func (op *TestTxt) String() string {
	return op.Str
}

func (*TestTxt) Compose() composer.Spec {
	return composer.Spec{
		Name:        TestTxt_Type,
		Uses:        composer.Type_Str,
		OpenStrings: true,
	}
}

const TestTxt_Type = "test_txt"

func (op *TestTxt) Marshal(m jsn.Marshaler) error {
	return TestTxt_Marshal(m, op)
}

func TestTxt_Optional_Marshal(m jsn.Marshaler, val *TestTxt) (err error) {
	var zero TestTxt
	if enc := m.IsEncoding(); !enc || val.Str != zero.Str {
		err = TestTxt_Marshal(m, val)
	}
	return
}

func TestTxt_Marshal(m jsn.Marshaler, val *TestTxt) (err error) {
	return m.MarshalValue(TestTxt_Type, &val.Str)
}

type TestTxt_Slice []TestTxt

func (op *TestTxt_Slice) GetType() string { return TestTxt_Type }

func (op *TestTxt_Slice) Marshal(m jsn.Marshaler) error {
	return TestTxt_Repeats_Marshal(m, (*[]TestTxt)(op))
}

func (op *TestTxt_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TestTxt_Slice) SetSize(cnt int) {
	var els []TestTxt
	if cnt >= 0 {
		els = make(TestTxt_Slice, cnt)
	}
	(*op) = els
}

func (op *TestTxt_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TestTxt_Marshal(m, &(*op)[i])
}

func TestTxt_Repeats_Marshal(m jsn.Marshaler, vals *[]TestTxt) error {
	return jsn.RepeatBlock(m, (*TestTxt_Slice)(vals))
}

func TestTxt_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TestTxt) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TestTxt_Repeats_Marshal(m, pv)
	}
	return
}

var Slots = []interface{}{
	(*TestSlot)(nil),
}

var Slats = []composer.Composer{
	(*TestBool)(nil),
	(*TestEmbed)(nil),
	(*TestFlow)(nil),
	(*TestNum)(nil),
	(*TestStr)(nil),
	(*TestSwap)(nil),
	(*TestTxt)(nil),
}

var Signatures = map[uint64]interface{}{
	5375871781273392485:  (*TestBool)(nil),  /* TestBool: */
	2921095451308776491:  (*TestNum)(nil),   /* TestNum: */
	8940712013797765950:  (*TestStr)(nil),   /* TestStr: */
	7140800509473528945:  (*TestSwap)(nil),  /* TestSwap a: */
	7137983560682620038:  (*TestSwap)(nil),  /* TestSwap b: */
	7138970922124564291:  (*TestSwap)(nil),  /* TestSwap c: */
	9144781193212880495:  (*TestTxt)(nil),   /* TestTxt: */
	4674877661722404900:  (*TestEmbed)(nil), /* test_slot=Embed testFlow: */
	4003780789222026231:  (*TestFlow)(nil),  /* test_slot=Flow bool:swap a: */
	18277657393474104982: (*TestFlow)(nil),  /* test_slot=Flow bool:swap a:slots: */
	4004754956524431952:  (*TestFlow)(nil),  /* test_slot=Flow bool:swap b: */
	13771748253966081413: (*TestFlow)(nil),  /* test_slot=Flow bool:swap b:slots: */
	4005610376570990885:  (*TestFlow)(nil),  /* test_slot=Flow bool:swap c: */
	5674371060339167924:  (*TestFlow)(nil),  /* test_slot=Flow bool:swap c:slots: */
	1264674884072578481:  (*TestFlow)(nil),  /* test_slot=Flow num:bool:swap a: */
	9254416107583163512:  (*TestFlow)(nil),  /* test_slot=Flow num:bool:swap a:slots: */
	1261857935281669574:  (*TestFlow)(nil),  /* test_slot=Flow num:bool:swap b: */
	8705586475388872315:  (*TestFlow)(nil),  /* test_slot=Flow num:bool:swap b:slots: */
	1262845296723613827:  (*TestFlow)(nil),  /* test_slot=Flow num:bool:swap c: */
	9618091095372893930:  (*TestFlow)(nil),  /* test_slot=Flow num:bool:swap c:slots: */
	4892874257338497905:  (*TestFlow)(nil),  /* test_slot=Flow num:swap a: */
	9741845409022596280:  (*TestFlow)(nil),  /* test_slot=Flow num:swap a:slots: */
	4890057308547588998:  (*TestFlow)(nil),  /* test_slot=Flow num:swap b: */
	1338231074237234875:  (*TestFlow)(nil),  /* test_slot=Flow num:swap b:slots: */
	4891044669989533251:  (*TestFlow)(nil),  /* test_slot=Flow num:swap c: */
	15844288127066246442: (*TestFlow)(nil),  /* test_slot=Flow num:swap c:slots: */
	7159433199210581651:  (*TestFlow)(nil),  /* test_slot=Flow slot:bool:swap a: */
	181618760168688954:   (*TestFlow)(nil),  /* test_slot=Flow slot:bool:swap a:slots: */
	7160424958699038748:  (*TestFlow)(nil),  /* test_slot=Flow slot:bool:swap b: */
	9310949923107569881:  (*TestFlow)(nil),  /* test_slot=Flow slot:bool:swap b:slots: */
	7161403524047957313:  (*TestFlow)(nil),  /* test_slot=Flow slot:bool:swap c: */
	15216561364608590984: (*TestFlow)(nil),  /* test_slot=Flow slot:bool:swap c:slots: */
	13186860388103412533: (*TestFlow)(nil),  /* test_slot=Flow slot:num:bool:swap a: */
	12434953562970918532: (*TestFlow)(nil),  /* test_slot=Flow slot:num:bool:swap a:slots: */
	13183902701824092618: (*TestFlow)(nil),  /* test_slot=Flow slot:num:bool:swap b: */
	17923545401664087015: (*TestFlow)(nil),  /* test_slot=Flow slot:num:bool:swap b:slots: */
	13184890063266036871: (*TestFlow)(nil),  /* test_slot=Flow slot:num:bool:swap c: */
	11216869237839205606: (*TestFlow)(nil),  /* test_slot=Flow slot:num:bool:swap c:slots: */
	16337239216761860165: (*TestFlow)(nil),  /* test_slot=Flow slot:num:swap a: */
	1415113059541595092:  (*TestFlow)(nil),  /* test_slot=Flow slot:num:swap a:slots: */
	16334422267970951258: (*TestFlow)(nil),  /* test_slot=Flow slot:num:swap b: */
	896964199814848439:   (*TestFlow)(nil),  /* test_slot=Flow slot:num:swap b:slots: */
	16335268891924484503: (*TestFlow)(nil),  /* test_slot=Flow slot:num:swap c: */
	6041950645603414198:  (*TestFlow)(nil),  /* test_slot=Flow slot:num:swap c:slots: */
	12912596366765834959: (*TestFlow)(nil),  /* test_slot=Flow slot:swap a: */
	14729520574941471838: (*TestFlow)(nil),  /* test_slot=Flow slot:swap a:slots: */
	12913570534068240680: (*TestFlow)(nil),  /* test_slot=Flow slot:swap b: */
	18292510269186445773: (*TestFlow)(nil),  /* test_slot=Flow slot:swap b:slots: */
	12914425954114799613: (*TestFlow)(nil),  /* test_slot=Flow slot:swap c: */
	5250437926060867164:  (*TestFlow)(nil),  /* test_slot=Flow slot:swap c:slots: */
	16957473258452611549: (*TestFlow)(nil),  /* test_slot=Flow slot:txt:bool:swap a: */
	177830915583890364:   (*TestFlow)(nil),  /* test_slot=Flow slot:txt:bool:swap a:slots: */
	16954515572173291634: (*TestFlow)(nil),  /* test_slot=Flow slot:txt:bool:swap b: */
	18331269787819217503: (*TestFlow)(nil),  /* test_slot=Flow slot:txt:bool:swap b:slots: */
	16955502933615235887: (*TestFlow)(nil),  /* test_slot=Flow slot:txt:bool:swap c: */
	15524614790799547710: (*TestFlow)(nil),  /* test_slot=Flow slot:txt:bool:swap c:slots: */
	3501137194810327367:  (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:bool:swap a: */
	3214438114935911462:  (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:bool:swap a:slots: */
	3502111362112733088:  (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:bool:swap b: */
	2368226698451690453:  (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:bool:swap b:slots: */
	3503107519647703029:  (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:bool:swap c: */
	17034485340072054980: (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:bool:swap c:slots: */
	4597499780223865491:  (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:swap a: */
	12018772744253862202: (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:swap a:slots: */
	4598491539712322588:  (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:swap b: */
	2701359833483191513:  (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:swap b:slots: */
	4599470105061241153:  (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:swap c: */
	8606971274984212616:  (*TestFlow)(nil),  /* test_slot=Flow slot:txt:num:swap c:slots: */
	394929995353645277:   (*TestFlow)(nil),  /* test_slot=Flow slot:txt:swap a: */
	12954962843971259068: (*TestFlow)(nil),  /* test_slot=Flow slot:txt:swap a:slots: */
	391972309074325362:   (*TestFlow)(nil),  /* test_slot=Flow slot:txt:swap b: */
	12661657642497034591: (*TestFlow)(nil),  /* test_slot=Flow slot:txt:swap b:slots: */
	392959670516269615:   (*TestFlow)(nil),  /* test_slot=Flow slot:txt:swap c: */
	9855002645477364798:  (*TestFlow)(nil),  /* test_slot=Flow slot:txt:swap c:slots: */
	11595017984637151299: (*TestFlow)(nil),  /* test_slot=Flow swap a: */
	9513548357935492906:  (*TestFlow)(nil),  /* test_slot=Flow swap a:slots: */
	11595869006637197388: (*TestFlow)(nil),  /* test_slot=Flow swap b: */
	3739883436483045385:  (*TestFlow)(nil),  /* test_slot=Flow swap b:slots: */
	11596847571986115953: (*TestFlow)(nil),  /* test_slot=Flow swap c: */
	3411105639891842744:  (*TestFlow)(nil),  /* test_slot=Flow swap c:slots: */
	6116333809511773473:  (*TestFlow)(nil),  /* test_slot=Flow txt:bool:swap a: */
	2524931204840297064:  (*TestFlow)(nil),  /* test_slot=Flow txt:bool:swap a:slots: */
	6113516860720864566:  (*TestFlow)(nil),  /* test_slot=Flow txt:bool:swap b: */
	1976101572646005867:  (*TestFlow)(nil),  /* test_slot=Flow txt:bool:swap b:slots: */
	6114363484674397811:  (*TestFlow)(nil),  /* test_slot=Flow txt:bool:swap c: */
	13546634146914304538: (*TestFlow)(nil),  /* test_slot=Flow txt:bool:swap c:slots: */
	3042176413511196771:  (*TestFlow)(nil),  /* test_slot=Flow txt:num:bool:swap a: */
	3475475419867441738:  (*TestFlow)(nil),  /* test_slot=Flow txt:num:bool:swap a:slots: */
	3043168172999653868:  (*TestFlow)(nil),  /* test_slot=Flow txt:num:bool:swap b: */
	3417798913527393321:  (*TestFlow)(nil),  /* test_slot=Flow txt:num:bool:swap b:slots: */
	3044146738348572433:  (*TestFlow)(nil),  /* test_slot=Flow txt:num:bool:swap c: */
	7697019104776277592:  (*TestFlow)(nil),  /* test_slot=Flow txt:num:bool:swap c:slots: */
	11209698323589727519: (*TestFlow)(nil),  /* test_slot=Flow txt:num:swap a: */
	4533301753114013614:  (*TestFlow)(nil),  /* test_slot=Flow txt:num:swap a:slots: */
	11210672490892133240: (*TestFlow)(nil),  /* test_slot=Flow txt:num:swap b: */
	6158740069692128797:  (*TestFlow)(nil),  /* test_slot=Flow txt:num:swap b:slots: */
	11211668648427103181: (*TestFlow)(nil),  /* test_slot=Flow txt:num:swap c: */
	16532553174846242860: (*TestFlow)(nil),  /* test_slot=Flow txt:num:swap c:slots: */
	6753852062274232673:  (*TestFlow)(nil),  /* test_slot=Flow txt:swap a: */
	9552505745174684456:  (*TestFlow)(nil),  /* test_slot=Flow txt:swap a:slots: */
	6751035113483323766:  (*TestFlow)(nil),  /* test_slot=Flow txt:swap b: */
	14603574056898801963: (*TestFlow)(nil),  /* test_slot=Flow txt:swap b:slots: */
	6751881737436857011:  (*TestFlow)(nil),  /* test_slot=Flow txt:swap c: */
	362344314439903450:   (*TestFlow)(nil),  /* test_slot=Flow txt:swap c:slots: */
}
