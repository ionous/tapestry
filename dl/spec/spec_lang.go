// Code generated by "makeops"; edit at your own risk.
package spec

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/jsn"
	"github.com/ionous/errutil"
)

// ChoiceSpec for swap choices
// if either label or type are not specified, they are derived from the name.
type ChoiceSpec struct {
	Name   string `if:"label=_,type=text"`
	Label  string `if:"label=label,optional,type=text"`
	Type   string `if:"label=type,optional,type=text"`
	Markup map[string]any
}

func (*ChoiceSpec) Compose() composer.Spec {
	return composer.Spec{
		Name: ChoiceSpec_Type,
		Uses: composer.Type_Flow,
		Lede: "pick",
	}
}

const ChoiceSpec_Type = "choice_spec"
const ChoiceSpec_Field_Name = "$NAME"
const ChoiceSpec_Field_Label = "$LABEL"
const ChoiceSpec_Field_Type = "$TYPE"

func (op *ChoiceSpec) Marshal(m jsn.Marshaler) error {
	return ChoiceSpec_Marshal(m, op)
}

type ChoiceSpec_Slice []ChoiceSpec

func (op *ChoiceSpec_Slice) GetType() string { return ChoiceSpec_Type }

func (op *ChoiceSpec_Slice) Marshal(m jsn.Marshaler) error {
	return ChoiceSpec_Repeats_Marshal(m, (*[]ChoiceSpec)(op))
}

func (op *ChoiceSpec_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *ChoiceSpec_Slice) SetSize(cnt int) {
	var els []ChoiceSpec
	if cnt >= 0 {
		els = make(ChoiceSpec_Slice, cnt)
	}
	(*op) = els
}

func (op *ChoiceSpec_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return ChoiceSpec_Marshal(m, &(*op)[i])
}

func ChoiceSpec_Repeats_Marshal(m jsn.Marshaler, vals *[]ChoiceSpec) error {
	return jsn.RepeatBlock(m, (*ChoiceSpec_Slice)(vals))
}

func ChoiceSpec_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]ChoiceSpec) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = ChoiceSpec_Repeats_Marshal(m, pv)
	}
	return
}

type ChoiceSpec_Flow struct{ ptr *ChoiceSpec }

func (n ChoiceSpec_Flow) GetType() string      { return ChoiceSpec_Type }
func (n ChoiceSpec_Flow) GetLede() string      { return "pick" }
func (n ChoiceSpec_Flow) GetFlow() interface{} { return n.ptr }
func (n ChoiceSpec_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*ChoiceSpec); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func ChoiceSpec_Optional_Marshal(m jsn.Marshaler, pv **ChoiceSpec) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = ChoiceSpec_Marshal(m, *pv)
	} else if !enc {
		var v ChoiceSpec
		if err = ChoiceSpec_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func ChoiceSpec_Marshal(m jsn.Marshaler, val *ChoiceSpec) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(ChoiceSpec_Flow{val}); err == nil {
		e0 := m.MarshalKey("", ChoiceSpec_Field_Name)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", ChoiceSpec_Field_Name))
		}
		e1 := m.MarshalKey("label", ChoiceSpec_Field_Label)
		if e1 == nil {
			e1 = prim.Text_Unboxed_Optional_Marshal(m, &val.Label)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", ChoiceSpec_Field_Label))
		}
		e2 := m.MarshalKey("type", ChoiceSpec_Field_Type)
		if e2 == nil {
			e2 = prim.Text_Unboxed_Optional_Marshal(m, &val.Type)
		}
		if e2 != nil && e2 != jsn.Missing {
			m.Error(errutil.New(e2, "in flow at", ChoiceSpec_Field_Type))
		}
		m.EndBlock()
	}
	return
}

// FlowSpec name: overrides the name of the operation used in compact story files.
// phrase: english text with embedded tokens referring to existing terms.
type FlowSpec struct {
	Name   string     `if:"label=_,optional,type=text"`
	Phrase string     `if:"label=phrase,optional,type=text"`
	Terms  []TermSpec `if:"label=uses"`
	Markup map[string]any
}

func (*FlowSpec) Compose() composer.Spec {
	return composer.Spec{
		Name: FlowSpec_Type,
		Uses: composer.Type_Flow,
		Lede: "flow",
	}
}

const FlowSpec_Type = "flow_spec"
const FlowSpec_Field_Name = "$NAME"
const FlowSpec_Field_Phrase = "$PHRASE"
const FlowSpec_Field_Terms = "$TERMS"

func (op *FlowSpec) Marshal(m jsn.Marshaler) error {
	return FlowSpec_Marshal(m, op)
}

type FlowSpec_Slice []FlowSpec

func (op *FlowSpec_Slice) GetType() string { return FlowSpec_Type }

func (op *FlowSpec_Slice) Marshal(m jsn.Marshaler) error {
	return FlowSpec_Repeats_Marshal(m, (*[]FlowSpec)(op))
}

func (op *FlowSpec_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *FlowSpec_Slice) SetSize(cnt int) {
	var els []FlowSpec
	if cnt >= 0 {
		els = make(FlowSpec_Slice, cnt)
	}
	(*op) = els
}

func (op *FlowSpec_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return FlowSpec_Marshal(m, &(*op)[i])
}

func FlowSpec_Repeats_Marshal(m jsn.Marshaler, vals *[]FlowSpec) error {
	return jsn.RepeatBlock(m, (*FlowSpec_Slice)(vals))
}

func FlowSpec_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]FlowSpec) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = FlowSpec_Repeats_Marshal(m, pv)
	}
	return
}

type FlowSpec_Flow struct{ ptr *FlowSpec }

func (n FlowSpec_Flow) GetType() string      { return FlowSpec_Type }
func (n FlowSpec_Flow) GetLede() string      { return "flow" }
func (n FlowSpec_Flow) GetFlow() interface{} { return n.ptr }
func (n FlowSpec_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*FlowSpec); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func FlowSpec_Optional_Marshal(m jsn.Marshaler, pv **FlowSpec) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = FlowSpec_Marshal(m, *pv)
	} else if !enc {
		var v FlowSpec
		if err = FlowSpec_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func FlowSpec_Marshal(m jsn.Marshaler, val *FlowSpec) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(FlowSpec_Flow{val}); err == nil {
		e0 := m.MarshalKey("", FlowSpec_Field_Name)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Optional_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", FlowSpec_Field_Name))
		}
		e1 := m.MarshalKey("phrase", FlowSpec_Field_Phrase)
		if e1 == nil {
			e1 = prim.Text_Unboxed_Optional_Marshal(m, &val.Phrase)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", FlowSpec_Field_Phrase))
		}
		e2 := m.MarshalKey("uses", FlowSpec_Field_Terms)
		if e2 == nil {
			e2 = TermSpec_Repeats_Marshal(m, &val.Terms)
		}
		if e2 != nil && e2 != jsn.Missing {
			m.Error(errutil.New(e2, "in flow at", FlowSpec_Field_Terms))
		}
		m.EndBlock()
	}
	return
}

// GroupSpec a collection of one or more other specs.
type GroupSpec struct {
	Specs  []TypeSpec `if:"label=contains"`
	Markup map[string]any
}

func (*GroupSpec) Compose() composer.Spec {
	return composer.Spec{
		Name: GroupSpec_Type,
		Uses: composer.Type_Flow,
		Lede: "group",
	}
}

const GroupSpec_Type = "group_spec"
const GroupSpec_Field_Specs = "$SPECS"

func (op *GroupSpec) Marshal(m jsn.Marshaler) error {
	return GroupSpec_Marshal(m, op)
}

type GroupSpec_Slice []GroupSpec

func (op *GroupSpec_Slice) GetType() string { return GroupSpec_Type }

func (op *GroupSpec_Slice) Marshal(m jsn.Marshaler) error {
	return GroupSpec_Repeats_Marshal(m, (*[]GroupSpec)(op))
}

func (op *GroupSpec_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *GroupSpec_Slice) SetSize(cnt int) {
	var els []GroupSpec
	if cnt >= 0 {
		els = make(GroupSpec_Slice, cnt)
	}
	(*op) = els
}

func (op *GroupSpec_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return GroupSpec_Marshal(m, &(*op)[i])
}

func GroupSpec_Repeats_Marshal(m jsn.Marshaler, vals *[]GroupSpec) error {
	return jsn.RepeatBlock(m, (*GroupSpec_Slice)(vals))
}

func GroupSpec_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]GroupSpec) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = GroupSpec_Repeats_Marshal(m, pv)
	}
	return
}

type GroupSpec_Flow struct{ ptr *GroupSpec }

func (n GroupSpec_Flow) GetType() string      { return GroupSpec_Type }
func (n GroupSpec_Flow) GetLede() string      { return "group" }
func (n GroupSpec_Flow) GetFlow() interface{} { return n.ptr }
func (n GroupSpec_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*GroupSpec); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func GroupSpec_Optional_Marshal(m jsn.Marshaler, pv **GroupSpec) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = GroupSpec_Marshal(m, *pv)
	} else if !enc {
		var v GroupSpec
		if err = GroupSpec_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func GroupSpec_Marshal(m jsn.Marshaler, val *GroupSpec) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(GroupSpec_Flow{val}); err == nil {
		e0 := m.MarshalKey("contains", GroupSpec_Field_Specs)
		if e0 == nil {
			e0 = TypeSpec_Repeats_Marshal(m, &val.Specs)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", GroupSpec_Field_Specs))
		}
		m.EndBlock()
	}
	return
}

// NumSpec when exclusive is true means the user can only specify one of the options
// otherwise the options are treated as hints.
type NumSpec struct {
	Name        string    `if:"label=_,optional,type=text"`
	Exclusively bool      `if:"label=exclusively,optional,type=bool"`
	Uses        []float64 `if:"label=uses,type=number"`
	Markup      map[string]any
}

func (*NumSpec) Compose() composer.Spec {
	return composer.Spec{
		Name: NumSpec_Type,
		Uses: composer.Type_Flow,
		Lede: "num",
	}
}

const NumSpec_Type = "num_spec"
const NumSpec_Field_Name = "$NAME"
const NumSpec_Field_Exclusively = "$EXCLUSIVELY"
const NumSpec_Field_Uses = "$USES"

func (op *NumSpec) Marshal(m jsn.Marshaler) error {
	return NumSpec_Marshal(m, op)
}

type NumSpec_Slice []NumSpec

func (op *NumSpec_Slice) GetType() string { return NumSpec_Type }

func (op *NumSpec_Slice) Marshal(m jsn.Marshaler) error {
	return NumSpec_Repeats_Marshal(m, (*[]NumSpec)(op))
}

func (op *NumSpec_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *NumSpec_Slice) SetSize(cnt int) {
	var els []NumSpec
	if cnt >= 0 {
		els = make(NumSpec_Slice, cnt)
	}
	(*op) = els
}

func (op *NumSpec_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return NumSpec_Marshal(m, &(*op)[i])
}

func NumSpec_Repeats_Marshal(m jsn.Marshaler, vals *[]NumSpec) error {
	return jsn.RepeatBlock(m, (*NumSpec_Slice)(vals))
}

func NumSpec_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]NumSpec) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = NumSpec_Repeats_Marshal(m, pv)
	}
	return
}

type NumSpec_Flow struct{ ptr *NumSpec }

func (n NumSpec_Flow) GetType() string      { return NumSpec_Type }
func (n NumSpec_Flow) GetLede() string      { return "num" }
func (n NumSpec_Flow) GetFlow() interface{} { return n.ptr }
func (n NumSpec_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*NumSpec); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func NumSpec_Optional_Marshal(m jsn.Marshaler, pv **NumSpec) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = NumSpec_Marshal(m, *pv)
	} else if !enc {
		var v NumSpec
		if err = NumSpec_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func NumSpec_Marshal(m jsn.Marshaler, val *NumSpec) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(NumSpec_Flow{val}); err == nil {
		e0 := m.MarshalKey("", NumSpec_Field_Name)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Optional_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", NumSpec_Field_Name))
		}
		e1 := m.MarshalKey("exclusively", NumSpec_Field_Exclusively)
		if e1 == nil {
			e1 = prim.Bool_Unboxed_Optional_Marshal(m, &val.Exclusively)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", NumSpec_Field_Exclusively))
		}
		e2 := m.MarshalKey("uses", NumSpec_Field_Uses)
		if e2 == nil {
			e2 = prim.Number_Unboxed_Repeats_Marshal(m, &val.Uses)
		}
		if e2 != nil && e2 != jsn.Missing {
			m.Error(errutil.New(e2, "in flow at", NumSpec_Field_Uses))
		}
		m.EndBlock()
	}
	return
}

// OptionSpec for string options
// if the label isnt specified, its derived from the name.
type OptionSpec struct {
	Name   string `if:"label=_,type=text"`
	Label  string `if:"label=label,optional,type=text"`
	Markup map[string]any
}

func (*OptionSpec) Compose() composer.Spec {
	return composer.Spec{
		Name: OptionSpec_Type,
		Uses: composer.Type_Flow,
		Lede: "option",
	}
}

const OptionSpec_Type = "option_spec"
const OptionSpec_Field_Name = "$NAME"
const OptionSpec_Field_Label = "$LABEL"

func (op *OptionSpec) Marshal(m jsn.Marshaler) error {
	return OptionSpec_Marshal(m, op)
}

type OptionSpec_Slice []OptionSpec

func (op *OptionSpec_Slice) GetType() string { return OptionSpec_Type }

func (op *OptionSpec_Slice) Marshal(m jsn.Marshaler) error {
	return OptionSpec_Repeats_Marshal(m, (*[]OptionSpec)(op))
}

func (op *OptionSpec_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *OptionSpec_Slice) SetSize(cnt int) {
	var els []OptionSpec
	if cnt >= 0 {
		els = make(OptionSpec_Slice, cnt)
	}
	(*op) = els
}

func (op *OptionSpec_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return OptionSpec_Marshal(m, &(*op)[i])
}

func OptionSpec_Repeats_Marshal(m jsn.Marshaler, vals *[]OptionSpec) error {
	return jsn.RepeatBlock(m, (*OptionSpec_Slice)(vals))
}

func OptionSpec_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]OptionSpec) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = OptionSpec_Repeats_Marshal(m, pv)
	}
	return
}

type OptionSpec_Flow struct{ ptr *OptionSpec }

func (n OptionSpec_Flow) GetType() string      { return OptionSpec_Type }
func (n OptionSpec_Flow) GetLede() string      { return "option" }
func (n OptionSpec_Flow) GetFlow() interface{} { return n.ptr }
func (n OptionSpec_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*OptionSpec); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func OptionSpec_Optional_Marshal(m jsn.Marshaler, pv **OptionSpec) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = OptionSpec_Marshal(m, *pv)
	} else if !enc {
		var v OptionSpec
		if err = OptionSpec_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func OptionSpec_Marshal(m jsn.Marshaler, val *OptionSpec) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(OptionSpec_Flow{val}); err == nil {
		e0 := m.MarshalKey("", OptionSpec_Field_Name)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", OptionSpec_Field_Name))
		}
		e1 := m.MarshalKey("label", OptionSpec_Field_Label)
		if e1 == nil {
			e1 = prim.Text_Unboxed_Optional_Marshal(m, &val.Label)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", OptionSpec_Field_Label))
		}
		m.EndBlock()
	}
	return
}

// SlotSpec A member of a flow which any of the other types can opt into.
// Aka an interface.
type SlotSpec struct {
	Markup map[string]any
}

func (*SlotSpec) Compose() composer.Spec {
	return composer.Spec{
		Name: SlotSpec_Type,
		Uses: composer.Type_Flow,
		Lede: "slot",
	}
}

const SlotSpec_Type = "slot_spec"

func (op *SlotSpec) Marshal(m jsn.Marshaler) error {
	return SlotSpec_Marshal(m, op)
}

type SlotSpec_Slice []SlotSpec

func (op *SlotSpec_Slice) GetType() string { return SlotSpec_Type }

func (op *SlotSpec_Slice) Marshal(m jsn.Marshaler) error {
	return SlotSpec_Repeats_Marshal(m, (*[]SlotSpec)(op))
}

func (op *SlotSpec_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *SlotSpec_Slice) SetSize(cnt int) {
	var els []SlotSpec
	if cnt >= 0 {
		els = make(SlotSpec_Slice, cnt)
	}
	(*op) = els
}

func (op *SlotSpec_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return SlotSpec_Marshal(m, &(*op)[i])
}

func SlotSpec_Repeats_Marshal(m jsn.Marshaler, vals *[]SlotSpec) error {
	return jsn.RepeatBlock(m, (*SlotSpec_Slice)(vals))
}

func SlotSpec_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]SlotSpec) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = SlotSpec_Repeats_Marshal(m, pv)
	}
	return
}

type SlotSpec_Flow struct{ ptr *SlotSpec }

func (n SlotSpec_Flow) GetType() string      { return SlotSpec_Type }
func (n SlotSpec_Flow) GetLede() string      { return "slot" }
func (n SlotSpec_Flow) GetFlow() interface{} { return n.ptr }
func (n SlotSpec_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*SlotSpec); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func SlotSpec_Optional_Marshal(m jsn.Marshaler, pv **SlotSpec) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = SlotSpec_Marshal(m, *pv)
	} else if !enc {
		var v SlotSpec
		if err = SlotSpec_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func SlotSpec_Marshal(m jsn.Marshaler, val *SlotSpec) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(SlotSpec_Flow{val}); err == nil {
		m.EndBlock()
	}
	return
}

// StrSpec when exclusive is true means the user can only specify one of the options
// otherwise the options are treated as hints.
type StrSpec struct {
	Name        string       `if:"label=_,optional,type=text"`
	Exclusively bool         `if:"label=exclusively,optional,type=bool"`
	Uses        []OptionSpec `if:"label=uses"`
	Markup      map[string]any
}

func (*StrSpec) Compose() composer.Spec {
	return composer.Spec{
		Name: StrSpec_Type,
		Uses: composer.Type_Flow,
		Lede: "str",
	}
}

const StrSpec_Type = "str_spec"
const StrSpec_Field_Name = "$NAME"
const StrSpec_Field_Exclusively = "$EXCLUSIVELY"
const StrSpec_Field_Uses = "$USES"

func (op *StrSpec) Marshal(m jsn.Marshaler) error {
	return StrSpec_Marshal(m, op)
}

type StrSpec_Slice []StrSpec

func (op *StrSpec_Slice) GetType() string { return StrSpec_Type }

func (op *StrSpec_Slice) Marshal(m jsn.Marshaler) error {
	return StrSpec_Repeats_Marshal(m, (*[]StrSpec)(op))
}

func (op *StrSpec_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *StrSpec_Slice) SetSize(cnt int) {
	var els []StrSpec
	if cnt >= 0 {
		els = make(StrSpec_Slice, cnt)
	}
	(*op) = els
}

func (op *StrSpec_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return StrSpec_Marshal(m, &(*op)[i])
}

func StrSpec_Repeats_Marshal(m jsn.Marshaler, vals *[]StrSpec) error {
	return jsn.RepeatBlock(m, (*StrSpec_Slice)(vals))
}

func StrSpec_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]StrSpec) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = StrSpec_Repeats_Marshal(m, pv)
	}
	return
}

type StrSpec_Flow struct{ ptr *StrSpec }

func (n StrSpec_Flow) GetType() string      { return StrSpec_Type }
func (n StrSpec_Flow) GetLede() string      { return "str" }
func (n StrSpec_Flow) GetFlow() interface{} { return n.ptr }
func (n StrSpec_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*StrSpec); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func StrSpec_Optional_Marshal(m jsn.Marshaler, pv **StrSpec) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = StrSpec_Marshal(m, *pv)
	} else if !enc {
		var v StrSpec
		if err = StrSpec_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func StrSpec_Marshal(m jsn.Marshaler, val *StrSpec) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(StrSpec_Flow{val}); err == nil {
		e0 := m.MarshalKey("", StrSpec_Field_Name)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Optional_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", StrSpec_Field_Name))
		}
		e1 := m.MarshalKey("exclusively", StrSpec_Field_Exclusively)
		if e1 == nil {
			e1 = prim.Bool_Unboxed_Optional_Marshal(m, &val.Exclusively)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", StrSpec_Field_Exclusively))
		}
		e2 := m.MarshalKey("uses", StrSpec_Field_Uses)
		if e2 == nil {
			e2 = OptionSpec_Repeats_Marshal(m, &val.Uses)
		}
		if e2 != nil && e2 != jsn.Missing {
			m.Error(errutil.New(e2, "in flow at", StrSpec_Field_Uses))
		}
		m.EndBlock()
	}
	return
}

// SwapSpec specifies a choice between one or more other types.
type SwapSpec struct {
	Name    string       `if:"label=_,optional,type=text"`
	Between []ChoiceSpec `if:"label=between"`
	Markup  map[string]any
}

func (*SwapSpec) Compose() composer.Spec {
	return composer.Spec{
		Name: SwapSpec_Type,
		Uses: composer.Type_Flow,
		Lede: "swap",
	}
}

const SwapSpec_Type = "swap_spec"
const SwapSpec_Field_Name = "$NAME"
const SwapSpec_Field_Between = "$BETWEEN"

func (op *SwapSpec) Marshal(m jsn.Marshaler) error {
	return SwapSpec_Marshal(m, op)
}

type SwapSpec_Slice []SwapSpec

func (op *SwapSpec_Slice) GetType() string { return SwapSpec_Type }

func (op *SwapSpec_Slice) Marshal(m jsn.Marshaler) error {
	return SwapSpec_Repeats_Marshal(m, (*[]SwapSpec)(op))
}

func (op *SwapSpec_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *SwapSpec_Slice) SetSize(cnt int) {
	var els []SwapSpec
	if cnt >= 0 {
		els = make(SwapSpec_Slice, cnt)
	}
	(*op) = els
}

func (op *SwapSpec_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return SwapSpec_Marshal(m, &(*op)[i])
}

func SwapSpec_Repeats_Marshal(m jsn.Marshaler, vals *[]SwapSpec) error {
	return jsn.RepeatBlock(m, (*SwapSpec_Slice)(vals))
}

func SwapSpec_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]SwapSpec) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = SwapSpec_Repeats_Marshal(m, pv)
	}
	return
}

type SwapSpec_Flow struct{ ptr *SwapSpec }

func (n SwapSpec_Flow) GetType() string      { return SwapSpec_Type }
func (n SwapSpec_Flow) GetLede() string      { return "swap" }
func (n SwapSpec_Flow) GetFlow() interface{} { return n.ptr }
func (n SwapSpec_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*SwapSpec); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func SwapSpec_Optional_Marshal(m jsn.Marshaler, pv **SwapSpec) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = SwapSpec_Marshal(m, *pv)
	} else if !enc {
		var v SwapSpec
		if err = SwapSpec_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func SwapSpec_Marshal(m jsn.Marshaler, val *SwapSpec) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(SwapSpec_Flow{val}); err == nil {
		e0 := m.MarshalKey("", SwapSpec_Field_Name)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Optional_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", SwapSpec_Field_Name))
		}
		e1 := m.MarshalKey("between", SwapSpec_Field_Between)
		if e1 == nil {
			e1 = ChoiceSpec_Repeats_Marshal(m, &val.Between)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", SwapSpec_Field_Between))
		}
		m.EndBlock()
	}
	return
}

// TermSpec A member of a flow.
// The label doubles as the parameter name unless an explicit name is specified.
// The type, if not specified, uses the name.
type TermSpec struct {
	Label    string `if:"label=_,type=text"`
	Name     string `if:"label=name,optional,type=text"`
	Type     string `if:"label=type,optional,type=text"`
	Private  bool   `if:"label=private,optional,type=bool"`
	Optional bool   `if:"label=optional,optional,type=bool"`
	Repeats  bool   `if:"label=repeats,optional,type=bool"`
	Markup   map[string]any
}

func (*TermSpec) Compose() composer.Spec {
	return composer.Spec{
		Name: TermSpec_Type,
		Uses: composer.Type_Flow,
		Lede: "term",
	}
}

const TermSpec_Type = "term_spec"
const TermSpec_Field_Label = "$LABEL"
const TermSpec_Field_Name = "$NAME"
const TermSpec_Field_Type = "$TYPE"
const TermSpec_Field_Private = "$PRIVATE"
const TermSpec_Field_Optional = "$OPTIONAL"
const TermSpec_Field_Repeats = "$REPEATS"

func (op *TermSpec) Marshal(m jsn.Marshaler) error {
	return TermSpec_Marshal(m, op)
}

type TermSpec_Slice []TermSpec

func (op *TermSpec_Slice) GetType() string { return TermSpec_Type }

func (op *TermSpec_Slice) Marshal(m jsn.Marshaler) error {
	return TermSpec_Repeats_Marshal(m, (*[]TermSpec)(op))
}

func (op *TermSpec_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TermSpec_Slice) SetSize(cnt int) {
	var els []TermSpec
	if cnt >= 0 {
		els = make(TermSpec_Slice, cnt)
	}
	(*op) = els
}

func (op *TermSpec_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TermSpec_Marshal(m, &(*op)[i])
}

func TermSpec_Repeats_Marshal(m jsn.Marshaler, vals *[]TermSpec) error {
	return jsn.RepeatBlock(m, (*TermSpec_Slice)(vals))
}

func TermSpec_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TermSpec) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TermSpec_Repeats_Marshal(m, pv)
	}
	return
}

type TermSpec_Flow struct{ ptr *TermSpec }

func (n TermSpec_Flow) GetType() string      { return TermSpec_Type }
func (n TermSpec_Flow) GetLede() string      { return "term" }
func (n TermSpec_Flow) GetFlow() interface{} { return n.ptr }
func (n TermSpec_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*TermSpec); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func TermSpec_Optional_Marshal(m jsn.Marshaler, pv **TermSpec) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = TermSpec_Marshal(m, *pv)
	} else if !enc {
		var v TermSpec
		if err = TermSpec_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func TermSpec_Marshal(m jsn.Marshaler, val *TermSpec) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(TermSpec_Flow{val}); err == nil {
		e0 := m.MarshalKey("", TermSpec_Field_Label)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Marshal(m, &val.Label)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", TermSpec_Field_Label))
		}
		e1 := m.MarshalKey("name", TermSpec_Field_Name)
		if e1 == nil {
			e1 = prim.Text_Unboxed_Optional_Marshal(m, &val.Name)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", TermSpec_Field_Name))
		}
		e2 := m.MarshalKey("type", TermSpec_Field_Type)
		if e2 == nil {
			e2 = prim.Text_Unboxed_Optional_Marshal(m, &val.Type)
		}
		if e2 != nil && e2 != jsn.Missing {
			m.Error(errutil.New(e2, "in flow at", TermSpec_Field_Type))
		}
		e3 := m.MarshalKey("private", TermSpec_Field_Private)
		if e3 == nil {
			e3 = prim.Bool_Unboxed_Optional_Marshal(m, &val.Private)
		}
		if e3 != nil && e3 != jsn.Missing {
			m.Error(errutil.New(e3, "in flow at", TermSpec_Field_Private))
		}
		e4 := m.MarshalKey("optional", TermSpec_Field_Optional)
		if e4 == nil {
			e4 = prim.Bool_Unboxed_Optional_Marshal(m, &val.Optional)
		}
		if e4 != nil && e4 != jsn.Missing {
			m.Error(errutil.New(e4, "in flow at", TermSpec_Field_Optional))
		}
		e5 := m.MarshalKey("repeats", TermSpec_Field_Repeats)
		if e5 == nil {
			e5 = prim.Bool_Unboxed_Optional_Marshal(m, &val.Repeats)
		}
		if e5 != nil && e5 != jsn.Missing {
			m.Error(errutil.New(e5, "in flow at", TermSpec_Field_Repeats))
		}
		m.EndBlock()
	}
	return
}

// TypeSpec can optionally fit one or more slots, or be part of one or more groups.
type TypeSpec struct {
	Name   string   `if:"label=_,type=text"`
	Slots  []string `if:"label=slots,optional,type=text"`
	Groups []string `if:"label=groups,optional,type=text"`
	Spec   UsesSpec `if:"label=with"`
	Markup map[string]any
}

func (*TypeSpec) Compose() composer.Spec {
	return composer.Spec{
		Name: TypeSpec_Type,
		Uses: composer.Type_Flow,
		Lede: "spec",
	}
}

const TypeSpec_Type = "type_spec"
const TypeSpec_Field_Name = "$NAME"
const TypeSpec_Field_Slots = "$SLOTS"
const TypeSpec_Field_Groups = "$GROUPS"
const TypeSpec_Field_Spec = "$SPEC"

func (op *TypeSpec) Marshal(m jsn.Marshaler) error {
	return TypeSpec_Marshal(m, op)
}

type TypeSpec_Slice []TypeSpec

func (op *TypeSpec_Slice) GetType() string { return TypeSpec_Type }

func (op *TypeSpec_Slice) Marshal(m jsn.Marshaler) error {
	return TypeSpec_Repeats_Marshal(m, (*[]TypeSpec)(op))
}

func (op *TypeSpec_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *TypeSpec_Slice) SetSize(cnt int) {
	var els []TypeSpec
	if cnt >= 0 {
		els = make(TypeSpec_Slice, cnt)
	}
	(*op) = els
}

func (op *TypeSpec_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return TypeSpec_Marshal(m, &(*op)[i])
}

func TypeSpec_Repeats_Marshal(m jsn.Marshaler, vals *[]TypeSpec) error {
	return jsn.RepeatBlock(m, (*TypeSpec_Slice)(vals))
}

func TypeSpec_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]TypeSpec) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = TypeSpec_Repeats_Marshal(m, pv)
	}
	return
}

type TypeSpec_Flow struct{ ptr *TypeSpec }

func (n TypeSpec_Flow) GetType() string      { return TypeSpec_Type }
func (n TypeSpec_Flow) GetLede() string      { return "spec" }
func (n TypeSpec_Flow) GetFlow() interface{} { return n.ptr }
func (n TypeSpec_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*TypeSpec); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func TypeSpec_Optional_Marshal(m jsn.Marshaler, pv **TypeSpec) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = TypeSpec_Marshal(m, *pv)
	} else if !enc {
		var v TypeSpec
		if err = TypeSpec_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func TypeSpec_Marshal(m jsn.Marshaler, val *TypeSpec) (err error) {
	m.SetMarkup(&val.Markup)
	if err = m.MarshalBlock(TypeSpec_Flow{val}); err == nil {
		e0 := m.MarshalKey("", TypeSpec_Field_Name)
		if e0 == nil {
			e0 = prim.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", TypeSpec_Field_Name))
		}
		e1 := m.MarshalKey("slots", TypeSpec_Field_Slots)
		if e1 == nil {
			e1 = prim.Text_Unboxed_Optional_Repeats_Marshal(m, &val.Slots)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", TypeSpec_Field_Slots))
		}
		e2 := m.MarshalKey("groups", TypeSpec_Field_Groups)
		if e2 == nil {
			e2 = prim.Text_Unboxed_Optional_Repeats_Marshal(m, &val.Groups)
		}
		if e2 != nil && e2 != jsn.Missing {
			m.Error(errutil.New(e2, "in flow at", TypeSpec_Field_Groups))
		}
		e3 := m.MarshalKey("with", TypeSpec_Field_Spec)
		if e3 == nil {
			e3 = UsesSpec_Marshal(m, &val.Spec)
		}
		if e3 != nil && e3 != jsn.Missing {
			m.Error(errutil.New(e3, "in flow at", TypeSpec_Field_Spec))
		}
		m.EndBlock()
	}
	return
}

// UsesSpec swaps between various options
type UsesSpec struct {
	Choice string
	Value  interface{}
}

var UsesSpec_Optional_Marshal = UsesSpec_Marshal

const UsesSpec_Flow_Opt = "$FLOW"
const UsesSpec_Slot_Opt = "$SLOT"
const UsesSpec_Swap_Opt = "$SWAP"
const UsesSpec_Num_Opt = "$NUM"
const UsesSpec_Str_Opt = "$STR"
const UsesSpec_Group_Opt = "$GROUP"

func (*UsesSpec) Compose() composer.Spec {
	return composer.Spec{
		Name: UsesSpec_Type,
		Uses: composer.Type_Swap,
		Choices: []string{
			UsesSpec_Flow_Opt, UsesSpec_Slot_Opt, UsesSpec_Swap_Opt, UsesSpec_Num_Opt, UsesSpec_Str_Opt, UsesSpec_Group_Opt,
		},
		Swaps: []interface{}{
			(*FlowSpec)(nil),
			(*SlotSpec)(nil),
			(*SwapSpec)(nil),
			(*NumSpec)(nil),
			(*StrSpec)(nil),
			(*GroupSpec)(nil),
		},
	}
}

const UsesSpec_Type = "uses_spec"

func (op *UsesSpec) GetType() string { return UsesSpec_Type }

func (op *UsesSpec) GetSwap() (string, interface{}) {
	return op.Choice, op.Value
}

func (op *UsesSpec) SetSwap(c string) (okay bool) {
	switch c {
	case "":
		op.Choice, op.Value = c, nil
		okay = true
	case UsesSpec_Flow_Opt:
		op.Choice, op.Value = c, new(FlowSpec)
		okay = true
	case UsesSpec_Slot_Opt:
		op.Choice, op.Value = c, new(SlotSpec)
		okay = true
	case UsesSpec_Swap_Opt:
		op.Choice, op.Value = c, new(SwapSpec)
		okay = true
	case UsesSpec_Num_Opt:
		op.Choice, op.Value = c, new(NumSpec)
		okay = true
	case UsesSpec_Str_Opt:
		op.Choice, op.Value = c, new(StrSpec)
		okay = true
	case UsesSpec_Group_Opt:
		op.Choice, op.Value = c, new(GroupSpec)
		okay = true
	}
	return
}

func (op *UsesSpec) Marshal(m jsn.Marshaler) error {
	return UsesSpec_Marshal(m, op)
}
func UsesSpec_Marshal(m jsn.Marshaler, val *UsesSpec) (err error) {
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

type UsesSpec_Slice []UsesSpec

func (op *UsesSpec_Slice) GetType() string { return UsesSpec_Type }

func (op *UsesSpec_Slice) Marshal(m jsn.Marshaler) error {
	return UsesSpec_Repeats_Marshal(m, (*[]UsesSpec)(op))
}

func (op *UsesSpec_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *UsesSpec_Slice) SetSize(cnt int) {
	var els []UsesSpec
	if cnt >= 0 {
		els = make(UsesSpec_Slice, cnt)
	}
	(*op) = els
}

func (op *UsesSpec_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return UsesSpec_Marshal(m, &(*op)[i])
}

func UsesSpec_Repeats_Marshal(m jsn.Marshaler, vals *[]UsesSpec) error {
	return jsn.RepeatBlock(m, (*UsesSpec_Slice)(vals))
}

func UsesSpec_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]UsesSpec) (err error) {
	if len(*pv) > 0 || !m.IsEncoding() {
		err = UsesSpec_Repeats_Marshal(m, pv)
	}
	return
}

var Slats = []composer.Composer{
	(*ChoiceSpec)(nil),
	(*FlowSpec)(nil),
	(*GroupSpec)(nil),
	(*NumSpec)(nil),
	(*OptionSpec)(nil),
	(*SlotSpec)(nil),
	(*StrSpec)(nil),
	(*SwapSpec)(nil),
	(*TermSpec)(nil),
	(*TypeSpec)(nil),
	(*UsesSpec)(nil),
}

var Signatures = map[uint64]interface{}{
	14698812208590391766: (*FlowSpec)(nil),   /* Flow phrase:uses: */
	8905948956105429135:  (*FlowSpec)(nil),   /* Flow uses: */
	7595234708931350060:  (*FlowSpec)(nil),   /* Flow:phrase:uses: */
	2613986620197267333:  (*FlowSpec)(nil),   /* Flow:uses: */
	138586504350168357:   (*GroupSpec)(nil),  /* Group contains: */
	9596524370352414924:  (*NumSpec)(nil),    /* Num exclusively:uses: */
	15841691970385929929: (*NumSpec)(nil),    /* Num uses: */
	5727702813760381158:  (*NumSpec)(nil),    /* Num:exclusively:uses: */
	12482362841218587987: (*NumSpec)(nil),    /* Num:uses: */
	5060730203210562250:  (*OptionSpec)(nil), /* Option: */
	2942004952747038014:  (*OptionSpec)(nil), /* Option:label: */
	8392859005396761062:  (*ChoiceSpec)(nil), /* Pick: */
	10192590430810142114: (*ChoiceSpec)(nil), /* Pick:label: */
	6404764409320714584:  (*ChoiceSpec)(nil), /* Pick:label:type: */
	8961232031828440580:  (*ChoiceSpec)(nil), /* Pick:type: */
	10492849629325543857: (*SlotSpec)(nil),   /* Slot */
	17485902661291661948: (*TypeSpec)(nil),   /* Spec:groups:with flow: */
	17417948388962479437: (*TypeSpec)(nil),   /* Spec:groups:with group: */
	4116744036832007200:  (*TypeSpec)(nil),   /* Spec:groups:with num: */
	14451625265593984708: (*TypeSpec)(nil),   /* Spec:groups:with slot: */
	7847085851758485277:  (*TypeSpec)(nil),   /* Spec:groups:with str: */
	15096092061144629293: (*TypeSpec)(nil),   /* Spec:groups:with swap: */
	5788656145002585531:  (*TypeSpec)(nil),   /* Spec:slots:groups:with flow: */
	2507119670160696732:  (*TypeSpec)(nil),   /* Spec:slots:groups:with group: */
	6967971190068921277:  (*TypeSpec)(nil),   /* Spec:slots:groups:with num: */
	8822933540700262771:  (*TypeSpec)(nil),   /* Spec:slots:groups:with slot: */
	14255553336076960684: (*TypeSpec)(nil),   /* Spec:slots:groups:with str: */
	14958032321477141646: (*TypeSpec)(nil),   /* Spec:slots:groups:with swap: */
	12158415875727800091: (*TypeSpec)(nil),   /* Spec:slots:with flow: */
	2058056700198914812:  (*TypeSpec)(nil),   /* Spec:slots:with group: */
	12924911410748986525: (*TypeSpec)(nil),   /* Spec:slots:with num: */
	15192693271425477331: (*TypeSpec)(nil),   /* Spec:slots:with slot: */
	1520725515772350988:  (*TypeSpec)(nil),   /* Spec:slots:with str: */
	2880907241004393582:  (*TypeSpec)(nil),   /* Spec:slots:with swap: */
	12144554005550117210: (*TypeSpec)(nil),   /* Spec:with flow: */
	18279802407160252863: (*TypeSpec)(nil),   /* Spec:with group: */
	9821026128543823202:  (*TypeSpec)(nil),   /* Spec:with num: */
	7278989256844955438:  (*TypeSpec)(nil),   /* Spec:with slot: */
	17614241350246502339: (*TypeSpec)(nil),   /* Spec:with str: */
	13597584408969322871: (*TypeSpec)(nil),   /* Spec:with swap: */
	8582367713278291167:  (*StrSpec)(nil),    /* Str exclusively:uses: */
	13984444229646943790: (*StrSpec)(nil),    /* Str uses: */
	10668987574759716109: (*StrSpec)(nil),    /* Str:exclusively:uses: */
	428528925856302860:   (*StrSpec)(nil),    /* Str:uses: */
	7464120135721846248:  (*SwapSpec)(nil),   /* Swap between: */
	11989329170787493114: (*SwapSpec)(nil),   /* Swap:between: */
	15304434050446929015: (*TermSpec)(nil),   /* Term: */
	1908770916544404618:  (*TermSpec)(nil),   /* Term:name: */
	17493832597660097384: (*TermSpec)(nil),   /* Term:name:optional: */
	10725022187143183618: (*TermSpec)(nil),   /* Term:name:optional:repeats: */
	4977282734720197993:  (*TermSpec)(nil),   /* Term:name:private: */
	15797853350075264797: (*TermSpec)(nil),   /* Term:name:private:optional: */
	4115812042834722363:  (*TermSpec)(nil),   /* Term:name:private:optional:repeats: */
	879887130050172383:   (*TermSpec)(nil),   /* Term:name:private:repeats: */
	9693871627066126088:  (*TermSpec)(nil),   /* Term:name:repeats: */
	14926838499984539376: (*TermSpec)(nil),   /* Term:name:type: */
	8114448813652831122:  (*TermSpec)(nil),   /* Term:name:type:optional: */
	492122786547772928:   (*TermSpec)(nil),   /* Term:name:type:optional:repeats: */
	15623635766246171155: (*TermSpec)(nil),   /* Term:name:type:private: */
	15651998130776576987: (*TermSpec)(nil),   /* Term:name:type:private:optional: */
	18146194056106062933: (*TermSpec)(nil),   /* Term:name:type:private:optional:repeats: */
	3100397852351784925:  (*TermSpec)(nil),   /* Term:name:type:private:repeats: */
	2434020794060870746:  (*TermSpec)(nil),   /* Term:name:type:repeats: */
	388572864592520623:   (*TermSpec)(nil),   /* Term:optional: */
	15437561936737348737: (*TermSpec)(nil),   /* Term:optional:repeats: */
	2080513686410551768:  (*TermSpec)(nil),   /* Term:private: */
	9375811883841958090:  (*TermSpec)(nil),   /* Term:private:optional: */
	250186011862547400:   (*TermSpec)(nil),   /* Term:private:optional:repeats: */
	1412567880106050386:  (*TermSpec)(nil),   /* Term:private:repeats: */
	5884140170322273081:  (*TermSpec)(nil),   /* Term:repeats: */
	16790547352605647791: (*TermSpec)(nil),   /* Term:type: */
	148296051051493815:   (*TermSpec)(nil),   /* Term:type:optional: */
	11348649781638926841: (*TermSpec)(nil),   /* Term:type:optional:repeats: */
	8106935495453137808:  (*TermSpec)(nil),   /* Term:type:private: */
	11399967267716567474: (*TermSpec)(nil),   /* Term:type:private:optional: */
	7000278554921529504:  (*TermSpec)(nil),   /* Term:type:private:optional:repeats: */
	12085390846926567418: (*TermSpec)(nil),   /* Term:type:private:repeats: */
	6606983442525629057:  (*TermSpec)(nil),   /* Term:type:repeats: */
	18192555626865043140: (*UsesSpec)(nil),   /* UsesSpec flow: */
	12509906366465792917: (*UsesSpec)(nil),   /* UsesSpec group: */
	9657992029678132728:  (*UsesSpec)(nil),   /* UsesSpec num: */
	7859797713207631340:  (*UsesSpec)(nil),   /* UsesSpec slot: */
	2626842383597263285:  (*UsesSpec)(nil),   /* UsesSpec str: */
	4680024026168720597:  (*UsesSpec)(nil),   /* UsesSpec swap: */
}
