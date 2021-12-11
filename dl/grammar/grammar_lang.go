// Code generated by "makeops"; edit at your own risk.
package grammar

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/jsn"
	"github.com/ionous/errutil"
)

// Action makes a parser scanner producing a script defined action.
// User implements: ScannerMaker.
type Action struct {
	Action string `if:"label=_,type=text"`
}

func (*Action) Compose() composer.Spec {
	return composer.Spec{
		Name: Action_Type,
		Uses: composer.Type_Flow,
		Lede: "as",
	}
}

const Action_Type = "action"

const Action_Field_Action = "$ACTION"

func (op *Action) Marshal(m jsn.Marshaler) error {
	return Action_Marshal(m, op)
}

type Action_Slice []Action

func (op *Action_Slice) GetType() string { return Action_Type }

func (op *Action_Slice) Marshal(m jsn.Marshaler) error {
	return Action_Repeats_Marshal(m, (*[]Action)(op))
}

func (op *Action_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Action_Slice) SetSize(cnt int) {
	var els []Action
	if cnt >= 0 {
		els = make(Action_Slice, cnt)
	}
	(*op) = els
}

func (op *Action_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Action_Marshal(m, &(*op)[i])
}

func Action_Repeats_Marshal(m jsn.Marshaler, vals *[]Action) error {
	return jsn.RepeatBlock(m, (*Action_Slice)(vals))
}

func Action_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Action) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Action_Repeats_Marshal(m, pv)
	}
	return
}

type Action_Flow struct{ ptr *Action }

func (n Action_Flow) GetType() string      { return Action_Type }
func (n Action_Flow) GetLede() string      { return "as" }
func (n Action_Flow) GetFlow() interface{} { return n.ptr }
func (n Action_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Action); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Action_Optional_Marshal(m jsn.Marshaler, pv **Action) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Action_Marshal(m, *pv)
	} else if !enc {
		var v Action
		if err = Action_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Action_Marshal(m jsn.Marshaler, val *Action) (err error) {
	if err = m.MarshalBlock(Action_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Action_Field_Action)
		if e0 == nil {
			e0 = literal.Text_Unboxed_Marshal(m, &val.Action)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Action_Field_Action))
		}
		m.EndBlock()
	}
	return
}

// Alias allows the user to refer to a noun by one or more other terms.
// User implements: GrammarMaker.
type Alias struct {
	Names  []string `if:"label=_,type=text"`
	AsNoun string   `if:"label=as_noun,type=text"`
}

func (*Alias) Compose() composer.Spec {
	return composer.Spec{
		Name: Alias_Type,
		Uses: composer.Type_Flow,
	}
}

const Alias_Type = "alias"

const Alias_Field_Names = "$NAMES"
const Alias_Field_AsNoun = "$AS_NOUN"

func (op *Alias) Marshal(m jsn.Marshaler) error {
	return Alias_Marshal(m, op)
}

type Alias_Slice []Alias

func (op *Alias_Slice) GetType() string { return Alias_Type }

func (op *Alias_Slice) Marshal(m jsn.Marshaler) error {
	return Alias_Repeats_Marshal(m, (*[]Alias)(op))
}

func (op *Alias_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Alias_Slice) SetSize(cnt int) {
	var els []Alias
	if cnt >= 0 {
		els = make(Alias_Slice, cnt)
	}
	(*op) = els
}

func (op *Alias_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Alias_Marshal(m, &(*op)[i])
}

func Alias_Repeats_Marshal(m jsn.Marshaler, vals *[]Alias) error {
	return jsn.RepeatBlock(m, (*Alias_Slice)(vals))
}

func Alias_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Alias) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Alias_Repeats_Marshal(m, pv)
	}
	return
}

type Alias_Flow struct{ ptr *Alias }

func (n Alias_Flow) GetType() string      { return Alias_Type }
func (n Alias_Flow) GetLede() string      { return Alias_Type }
func (n Alias_Flow) GetFlow() interface{} { return n.ptr }
func (n Alias_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Alias); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Alias_Optional_Marshal(m jsn.Marshaler, pv **Alias) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Alias_Marshal(m, *pv)
	} else if !enc {
		var v Alias
		if err = Alias_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Alias_Marshal(m jsn.Marshaler, val *Alias) (err error) {
	if err = m.MarshalBlock(Alias_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Alias_Field_Names)
		if e0 == nil {
			e0 = literal.Text_Unboxed_Repeats_Marshal(m, &val.Names)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Alias_Field_Names))
		}
		e1 := m.MarshalKey("as_noun", Alias_Field_AsNoun)
		if e1 == nil {
			e1 = literal.Text_Unboxed_Marshal(m, &val.AsNoun)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", Alias_Field_AsNoun))
		}
		m.EndBlock()
	}
	return
}

// AllOf makes a parser scanner
// User implements: ScannerMaker.
type AllOf struct {
	Series []ScannerMaker `if:"label=_"`
}

func (*AllOf) Compose() composer.Spec {
	return composer.Spec{
		Name: AllOf_Type,
		Uses: composer.Type_Flow,
	}
}

const AllOf_Type = "all_of"

const AllOf_Field_Series = "$SERIES"

func (op *AllOf) Marshal(m jsn.Marshaler) error {
	return AllOf_Marshal(m, op)
}

type AllOf_Slice []AllOf

func (op *AllOf_Slice) GetType() string { return AllOf_Type }

func (op *AllOf_Slice) Marshal(m jsn.Marshaler) error {
	return AllOf_Repeats_Marshal(m, (*[]AllOf)(op))
}

func (op *AllOf_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *AllOf_Slice) SetSize(cnt int) {
	var els []AllOf
	if cnt >= 0 {
		els = make(AllOf_Slice, cnt)
	}
	(*op) = els
}

func (op *AllOf_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return AllOf_Marshal(m, &(*op)[i])
}

func AllOf_Repeats_Marshal(m jsn.Marshaler, vals *[]AllOf) error {
	return jsn.RepeatBlock(m, (*AllOf_Slice)(vals))
}

func AllOf_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]AllOf) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = AllOf_Repeats_Marshal(m, pv)
	}
	return
}

type AllOf_Flow struct{ ptr *AllOf }

func (n AllOf_Flow) GetType() string      { return AllOf_Type }
func (n AllOf_Flow) GetLede() string      { return AllOf_Type }
func (n AllOf_Flow) GetFlow() interface{} { return n.ptr }
func (n AllOf_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*AllOf); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func AllOf_Optional_Marshal(m jsn.Marshaler, pv **AllOf) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = AllOf_Marshal(m, *pv)
	} else if !enc {
		var v AllOf
		if err = AllOf_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func AllOf_Marshal(m jsn.Marshaler, val *AllOf) (err error) {
	if err = m.MarshalBlock(AllOf_Flow{val}); err == nil {
		e0 := m.MarshalKey("", AllOf_Field_Series)
		if e0 == nil {
			e0 = ScannerMaker_Repeats_Marshal(m, &val.Series)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", AllOf_Field_Series))
		}
		m.EndBlock()
	}
	return
}

// AnyOf makes a parser scanner
// User implements: ScannerMaker.
type AnyOf struct {
	Options []ScannerMaker `if:"label=_"`
}

func (*AnyOf) Compose() composer.Spec {
	return composer.Spec{
		Name: AnyOf_Type,
		Uses: composer.Type_Flow,
	}
}

const AnyOf_Type = "any_of"

const AnyOf_Field_Options = "$OPTIONS"

func (op *AnyOf) Marshal(m jsn.Marshaler) error {
	return AnyOf_Marshal(m, op)
}

type AnyOf_Slice []AnyOf

func (op *AnyOf_Slice) GetType() string { return AnyOf_Type }

func (op *AnyOf_Slice) Marshal(m jsn.Marshaler) error {
	return AnyOf_Repeats_Marshal(m, (*[]AnyOf)(op))
}

func (op *AnyOf_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *AnyOf_Slice) SetSize(cnt int) {
	var els []AnyOf
	if cnt >= 0 {
		els = make(AnyOf_Slice, cnt)
	}
	(*op) = els
}

func (op *AnyOf_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return AnyOf_Marshal(m, &(*op)[i])
}

func AnyOf_Repeats_Marshal(m jsn.Marshaler, vals *[]AnyOf) error {
	return jsn.RepeatBlock(m, (*AnyOf_Slice)(vals))
}

func AnyOf_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]AnyOf) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = AnyOf_Repeats_Marshal(m, pv)
	}
	return
}

type AnyOf_Flow struct{ ptr *AnyOf }

func (n AnyOf_Flow) GetType() string      { return AnyOf_Type }
func (n AnyOf_Flow) GetLede() string      { return AnyOf_Type }
func (n AnyOf_Flow) GetFlow() interface{} { return n.ptr }
func (n AnyOf_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*AnyOf); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func AnyOf_Optional_Marshal(m jsn.Marshaler, pv **AnyOf) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = AnyOf_Marshal(m, *pv)
	} else if !enc {
		var v AnyOf
		if err = AnyOf_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func AnyOf_Marshal(m jsn.Marshaler, val *AnyOf) (err error) {
	if err = m.MarshalBlock(AnyOf_Flow{val}); err == nil {
		e0 := m.MarshalKey("", AnyOf_Field_Options)
		if e0 == nil {
			e0 = ScannerMaker_Repeats_Marshal(m, &val.Options)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", AnyOf_Field_Options))
		}
		m.EndBlock()
	}
	return
}

// Directive starts a parser scanner
// User implements: GrammarMaker.
type Directive struct {
	Lede  []string       `if:"label=_,type=text"`
	Scans []ScannerMaker `if:"label=scans"`
}

func (*Directive) Compose() composer.Spec {
	return composer.Spec{
		Name: Directive_Type,
		Uses: composer.Type_Flow,
	}
}

const Directive_Type = "directive"

const Directive_Field_Lede = "$LEDE"
const Directive_Field_Scans = "$SCANS"

func (op *Directive) Marshal(m jsn.Marshaler) error {
	return Directive_Marshal(m, op)
}

type Directive_Slice []Directive

func (op *Directive_Slice) GetType() string { return Directive_Type }

func (op *Directive_Slice) Marshal(m jsn.Marshaler) error {
	return Directive_Repeats_Marshal(m, (*[]Directive)(op))
}

func (op *Directive_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Directive_Slice) SetSize(cnt int) {
	var els []Directive
	if cnt >= 0 {
		els = make(Directive_Slice, cnt)
	}
	(*op) = els
}

func (op *Directive_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Directive_Marshal(m, &(*op)[i])
}

func Directive_Repeats_Marshal(m jsn.Marshaler, vals *[]Directive) error {
	return jsn.RepeatBlock(m, (*Directive_Slice)(vals))
}

func Directive_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Directive) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Directive_Repeats_Marshal(m, pv)
	}
	return
}

type Directive_Flow struct{ ptr *Directive }

func (n Directive_Flow) GetType() string      { return Directive_Type }
func (n Directive_Flow) GetLede() string      { return Directive_Type }
func (n Directive_Flow) GetFlow() interface{} { return n.ptr }
func (n Directive_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Directive); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Directive_Optional_Marshal(m jsn.Marshaler, pv **Directive) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Directive_Marshal(m, *pv)
	} else if !enc {
		var v Directive
		if err = Directive_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Directive_Marshal(m jsn.Marshaler, val *Directive) (err error) {
	if err = m.MarshalBlock(Directive_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Directive_Field_Lede)
		if e0 == nil {
			e0 = literal.Text_Unboxed_Repeats_Marshal(m, &val.Lede)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Directive_Field_Lede))
		}
		e1 := m.MarshalKey("scans", Directive_Field_Scans)
		if e1 == nil {
			e1 = ScannerMaker_Repeats_Marshal(m, &val.Scans)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", Directive_Field_Scans))
		}
		m.EndBlock()
	}
	return
}

// Grammar Read what the player types and turn it into actions.
type Grammar struct {
	Grammar GrammarMaker `if:"label=_"`
}

func (*Grammar) Compose() composer.Spec {
	return composer.Spec{
		Name: Grammar_Type,
		Uses: composer.Type_Flow,
	}
}

const Grammar_Type = "grammar"

const Grammar_Field_Grammar = "$GRAMMAR"

func (op *Grammar) Marshal(m jsn.Marshaler) error {
	return Grammar_Marshal(m, op)
}

type Grammar_Slice []Grammar

func (op *Grammar_Slice) GetType() string { return Grammar_Type }

func (op *Grammar_Slice) Marshal(m jsn.Marshaler) error {
	return Grammar_Repeats_Marshal(m, (*[]Grammar)(op))
}

func (op *Grammar_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Grammar_Slice) SetSize(cnt int) {
	var els []Grammar
	if cnt >= 0 {
		els = make(Grammar_Slice, cnt)
	}
	(*op) = els
}

func (op *Grammar_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Grammar_Marshal(m, &(*op)[i])
}

func Grammar_Repeats_Marshal(m jsn.Marshaler, vals *[]Grammar) error {
	return jsn.RepeatBlock(m, (*Grammar_Slice)(vals))
}

func Grammar_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Grammar) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Grammar_Repeats_Marshal(m, pv)
	}
	return
}

type Grammar_Flow struct{ ptr *Grammar }

func (n Grammar_Flow) GetType() string      { return Grammar_Type }
func (n Grammar_Flow) GetLede() string      { return Grammar_Type }
func (n Grammar_Flow) GetFlow() interface{} { return n.ptr }
func (n Grammar_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Grammar); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Grammar_Optional_Marshal(m jsn.Marshaler, pv **Grammar) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Grammar_Marshal(m, *pv)
	} else if !enc {
		var v Grammar
		if err = Grammar_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Grammar_Marshal(m jsn.Marshaler, val *Grammar) (err error) {
	if err = m.MarshalBlock(Grammar_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Grammar_Field_Grammar)
		if e0 == nil {
			e0 = GrammarMaker_Marshal(m, &val.Grammar)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Grammar_Field_Grammar))
		}
		m.EndBlock()
	}
	return
}

const GrammarMaker_Type = "grammar_maker"

var GrammarMaker_Optional_Marshal = GrammarMaker_Marshal

type GrammarMaker_Slot struct{ ptr *GrammarMaker }

func (at GrammarMaker_Slot) GetType() string              { return GrammarMaker_Type }
func (at GrammarMaker_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at GrammarMaker_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(GrammarMaker)
	return
}

func GrammarMaker_Marshal(m jsn.Marshaler, ptr *GrammarMaker) (err error) {
	slot := GrammarMaker_Slot{ptr}
	if err = m.MarshalBlock(slot); err == nil {
		if a, ok := slot.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}

type GrammarMaker_Slice []GrammarMaker

func (op *GrammarMaker_Slice) GetType() string { return GrammarMaker_Type }

func (op *GrammarMaker_Slice) Marshal(m jsn.Marshaler) error {
	return GrammarMaker_Repeats_Marshal(m, (*[]GrammarMaker)(op))
}

func (op *GrammarMaker_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *GrammarMaker_Slice) SetSize(cnt int) {
	var els []GrammarMaker
	if cnt >= 0 {
		els = make(GrammarMaker_Slice, cnt)
	}
	(*op) = els
}

func (op *GrammarMaker_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return GrammarMaker_Marshal(m, &(*op)[i])
}

func GrammarMaker_Repeats_Marshal(m jsn.Marshaler, vals *[]GrammarMaker) error {
	return jsn.RepeatBlock(m, (*GrammarMaker_Slice)(vals))
}

func GrammarMaker_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]GrammarMaker) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = GrammarMaker_Repeats_Marshal(m, pv)
	}
	return
}

// Noun makes a parser scanner
// User implements: ScannerMaker.
type Noun struct {
	Kind string `if:"label=_,type=text"`
}

func (*Noun) Compose() composer.Spec {
	return composer.Spec{
		Name: Noun_Type,
		Uses: composer.Type_Flow,
	}
}

const Noun_Type = "noun"

const Noun_Field_Kind = "$KIND"

func (op *Noun) Marshal(m jsn.Marshaler) error {
	return Noun_Marshal(m, op)
}

type Noun_Slice []Noun

func (op *Noun_Slice) GetType() string { return Noun_Type }

func (op *Noun_Slice) Marshal(m jsn.Marshaler) error {
	return Noun_Repeats_Marshal(m, (*[]Noun)(op))
}

func (op *Noun_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Noun_Slice) SetSize(cnt int) {
	var els []Noun
	if cnt >= 0 {
		els = make(Noun_Slice, cnt)
	}
	(*op) = els
}

func (op *Noun_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Noun_Marshal(m, &(*op)[i])
}

func Noun_Repeats_Marshal(m jsn.Marshaler, vals *[]Noun) error {
	return jsn.RepeatBlock(m, (*Noun_Slice)(vals))
}

func Noun_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Noun) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Noun_Repeats_Marshal(m, pv)
	}
	return
}

type Noun_Flow struct{ ptr *Noun }

func (n Noun_Flow) GetType() string      { return Noun_Type }
func (n Noun_Flow) GetLede() string      { return Noun_Type }
func (n Noun_Flow) GetFlow() interface{} { return n.ptr }
func (n Noun_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Noun); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Noun_Optional_Marshal(m jsn.Marshaler, pv **Noun) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Noun_Marshal(m, *pv)
	} else if !enc {
		var v Noun
		if err = Noun_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Noun_Marshal(m jsn.Marshaler, val *Noun) (err error) {
	if err = m.MarshalBlock(Noun_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Noun_Field_Kind)
		if e0 == nil {
			e0 = literal.Text_Unboxed_Marshal(m, &val.Kind)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Noun_Field_Kind))
		}
		m.EndBlock()
	}
	return
}

// Retarget makes a parser scanner
// User implements: ScannerMaker.
type Retarget struct {
	Span []ScannerMaker `if:"label=_"`
}

func (*Retarget) Compose() composer.Spec {
	return composer.Spec{
		Name: Retarget_Type,
		Uses: composer.Type_Flow,
	}
}

const Retarget_Type = "retarget"

const Retarget_Field_Span = "$SPAN"

func (op *Retarget) Marshal(m jsn.Marshaler) error {
	return Retarget_Marshal(m, op)
}

type Retarget_Slice []Retarget

func (op *Retarget_Slice) GetType() string { return Retarget_Type }

func (op *Retarget_Slice) Marshal(m jsn.Marshaler) error {
	return Retarget_Repeats_Marshal(m, (*[]Retarget)(op))
}

func (op *Retarget_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Retarget_Slice) SetSize(cnt int) {
	var els []Retarget
	if cnt >= 0 {
		els = make(Retarget_Slice, cnt)
	}
	(*op) = els
}

func (op *Retarget_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Retarget_Marshal(m, &(*op)[i])
}

func Retarget_Repeats_Marshal(m jsn.Marshaler, vals *[]Retarget) error {
	return jsn.RepeatBlock(m, (*Retarget_Slice)(vals))
}

func Retarget_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Retarget) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Retarget_Repeats_Marshal(m, pv)
	}
	return
}

type Retarget_Flow struct{ ptr *Retarget }

func (n Retarget_Flow) GetType() string      { return Retarget_Type }
func (n Retarget_Flow) GetLede() string      { return Retarget_Type }
func (n Retarget_Flow) GetFlow() interface{} { return n.ptr }
func (n Retarget_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Retarget); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Retarget_Optional_Marshal(m jsn.Marshaler, pv **Retarget) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Retarget_Marshal(m, *pv)
	} else if !enc {
		var v Retarget
		if err = Retarget_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Retarget_Marshal(m jsn.Marshaler, val *Retarget) (err error) {
	if err = m.MarshalBlock(Retarget_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Retarget_Field_Span)
		if e0 == nil {
			e0 = ScannerMaker_Repeats_Marshal(m, &val.Span)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Retarget_Field_Span))
		}
		m.EndBlock()
	}
	return
}

// Reverse makes a parser scanner
// User implements: ScannerMaker.
type Reverse struct {
	Reverses []ScannerMaker `if:"label=_"`
}

func (*Reverse) Compose() composer.Spec {
	return composer.Spec{
		Name: Reverse_Type,
		Uses: composer.Type_Flow,
	}
}

const Reverse_Type = "reverse"

const Reverse_Field_Reverses = "$REVERSES"

func (op *Reverse) Marshal(m jsn.Marshaler) error {
	return Reverse_Marshal(m, op)
}

type Reverse_Slice []Reverse

func (op *Reverse_Slice) GetType() string { return Reverse_Type }

func (op *Reverse_Slice) Marshal(m jsn.Marshaler) error {
	return Reverse_Repeats_Marshal(m, (*[]Reverse)(op))
}

func (op *Reverse_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Reverse_Slice) SetSize(cnt int) {
	var els []Reverse
	if cnt >= 0 {
		els = make(Reverse_Slice, cnt)
	}
	(*op) = els
}

func (op *Reverse_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Reverse_Marshal(m, &(*op)[i])
}

func Reverse_Repeats_Marshal(m jsn.Marshaler, vals *[]Reverse) error {
	return jsn.RepeatBlock(m, (*Reverse_Slice)(vals))
}

func Reverse_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Reverse) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Reverse_Repeats_Marshal(m, pv)
	}
	return
}

type Reverse_Flow struct{ ptr *Reverse }

func (n Reverse_Flow) GetType() string      { return Reverse_Type }
func (n Reverse_Flow) GetLede() string      { return Reverse_Type }
func (n Reverse_Flow) GetFlow() interface{} { return n.ptr }
func (n Reverse_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Reverse); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Reverse_Optional_Marshal(m jsn.Marshaler, pv **Reverse) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Reverse_Marshal(m, *pv)
	} else if !enc {
		var v Reverse
		if err = Reverse_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Reverse_Marshal(m jsn.Marshaler, val *Reverse) (err error) {
	if err = m.MarshalBlock(Reverse_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Reverse_Field_Reverses)
		if e0 == nil {
			e0 = ScannerMaker_Repeats_Marshal(m, &val.Reverses)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Reverse_Field_Reverses))
		}
		m.EndBlock()
	}
	return
}

const ScannerMaker_Type = "scanner_maker"

var ScannerMaker_Optional_Marshal = ScannerMaker_Marshal

type ScannerMaker_Slot struct{ ptr *ScannerMaker }

func (at ScannerMaker_Slot) GetType() string              { return ScannerMaker_Type }
func (at ScannerMaker_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at ScannerMaker_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(ScannerMaker)
	return
}

func ScannerMaker_Marshal(m jsn.Marshaler, ptr *ScannerMaker) (err error) {
	slot := ScannerMaker_Slot{ptr}
	if err = m.MarshalBlock(slot); err == nil {
		if a, ok := slot.GetSlot(); ok {
			if e := a.(jsn.Marshalee).Marshal(m); e != nil && e != jsn.Missing {
				m.Error(e)
			}
		}
		m.EndBlock()
	}
	return
}

type ScannerMaker_Slice []ScannerMaker

func (op *ScannerMaker_Slice) GetType() string { return ScannerMaker_Type }

func (op *ScannerMaker_Slice) Marshal(m jsn.Marshaler) error {
	return ScannerMaker_Repeats_Marshal(m, (*[]ScannerMaker)(op))
}

func (op *ScannerMaker_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *ScannerMaker_Slice) SetSize(cnt int) {
	var els []ScannerMaker
	if cnt >= 0 {
		els = make(ScannerMaker_Slice, cnt)
	}
	(*op) = els
}

func (op *ScannerMaker_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return ScannerMaker_Marshal(m, &(*op)[i])
}

func ScannerMaker_Repeats_Marshal(m jsn.Marshaler, vals *[]ScannerMaker) error {
	return jsn.RepeatBlock(m, (*ScannerMaker_Slice)(vals))
}

func ScannerMaker_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]ScannerMaker) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = ScannerMaker_Repeats_Marshal(m, pv)
	}
	return
}

// Self makes a parser scanner which matches the player. ( the player string is just to make the composer happy. )
// User implements: ScannerMaker.
type Self struct {
	Player string `if:"label=_,type=text"`
}

func (*Self) Compose() composer.Spec {
	return composer.Spec{
		Name: Self_Type,
		Uses: composer.Type_Flow,
	}
}

const Self_Type = "self"

const Self_Field_Player = "$PLAYER"

func (op *Self) Marshal(m jsn.Marshaler) error {
	return Self_Marshal(m, op)
}

type Self_Slice []Self

func (op *Self_Slice) GetType() string { return Self_Type }

func (op *Self_Slice) Marshal(m jsn.Marshaler) error {
	return Self_Repeats_Marshal(m, (*[]Self)(op))
}

func (op *Self_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Self_Slice) SetSize(cnt int) {
	var els []Self
	if cnt >= 0 {
		els = make(Self_Slice, cnt)
	}
	(*op) = els
}

func (op *Self_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Self_Marshal(m, &(*op)[i])
}

func Self_Repeats_Marshal(m jsn.Marshaler, vals *[]Self) error {
	return jsn.RepeatBlock(m, (*Self_Slice)(vals))
}

func Self_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Self) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Self_Repeats_Marshal(m, pv)
	}
	return
}

type Self_Flow struct{ ptr *Self }

func (n Self_Flow) GetType() string      { return Self_Type }
func (n Self_Flow) GetLede() string      { return Self_Type }
func (n Self_Flow) GetFlow() interface{} { return n.ptr }
func (n Self_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Self); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Self_Optional_Marshal(m jsn.Marshaler, pv **Self) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Self_Marshal(m, *pv)
	} else if !enc {
		var v Self
		if err = Self_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Self_Marshal(m jsn.Marshaler, val *Self) (err error) {
	if err = m.MarshalBlock(Self_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Self_Field_Player)
		if e0 == nil {
			e0 = literal.Text_Unboxed_Marshal(m, &val.Player)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Self_Field_Player))
		}
		m.EndBlock()
	}
	return
}

// Words makes a parser scanner
// User implements: ScannerMaker.
type Words struct {
	Words []string `if:"label=_,type=text"`
}

func (*Words) Compose() composer.Spec {
	return composer.Spec{
		Name: Words_Type,
		Uses: composer.Type_Flow,
	}
}

const Words_Type = "words"

const Words_Field_Words = "$WORDS"

func (op *Words) Marshal(m jsn.Marshaler) error {
	return Words_Marshal(m, op)
}

type Words_Slice []Words

func (op *Words_Slice) GetType() string { return Words_Type }

func (op *Words_Slice) Marshal(m jsn.Marshaler) error {
	return Words_Repeats_Marshal(m, (*[]Words)(op))
}

func (op *Words_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Words_Slice) SetSize(cnt int) {
	var els []Words
	if cnt >= 0 {
		els = make(Words_Slice, cnt)
	}
	(*op) = els
}

func (op *Words_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Words_Marshal(m, &(*op)[i])
}

func Words_Repeats_Marshal(m jsn.Marshaler, vals *[]Words) error {
	return jsn.RepeatBlock(m, (*Words_Slice)(vals))
}

func Words_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Words) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Words_Repeats_Marshal(m, pv)
	}
	return
}

type Words_Flow struct{ ptr *Words }

func (n Words_Flow) GetType() string      { return Words_Type }
func (n Words_Flow) GetLede() string      { return Words_Type }
func (n Words_Flow) GetFlow() interface{} { return n.ptr }
func (n Words_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*Words); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func Words_Optional_Marshal(m jsn.Marshaler, pv **Words) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = Words_Marshal(m, *pv)
	} else if !enc {
		var v Words
		if err = Words_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func Words_Marshal(m jsn.Marshaler, val *Words) (err error) {
	if err = m.MarshalBlock(Words_Flow{val}); err == nil {
		e0 := m.MarshalKey("", Words_Field_Words)
		if e0 == nil {
			e0 = literal.Text_Unboxed_Repeats_Marshal(m, &val.Words)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", Words_Field_Words))
		}
		m.EndBlock()
	}
	return
}

var Slots = []interface{}{
	(*GrammarMaker)(nil),
	(*ScannerMaker)(nil),
}

var Slats = []composer.Composer{
	(*Action)(nil),
	(*Alias)(nil),
	(*AllOf)(nil),
	(*AnyOf)(nil),
	(*Directive)(nil),
	(*Grammar)(nil),
	(*Noun)(nil),
	(*Retarget)(nil),
	(*Reverse)(nil),
	(*Self)(nil),
	(*Words)(nil),
}

var Signatures = map[uint64]interface{}{
	18013434347847705365: (*Action)(nil),    /* As: */
	12970384957739247299: (*Alias)(nil),     /* Alias:asNoun: */
	12299258133038749149: (*AllOf)(nil),     /* AllOf: */
	5638555781748976558:  (*AnyOf)(nil),     /* AnyOf: */
	13009220793665599564: (*Directive)(nil), /* Directive:scans: */
	15013060695242199180: (*Grammar)(nil),   /* Grammar: */
	571163134278291657:   (*Noun)(nil),      /* Noun: */
	9105733481983959033:  (*Retarget)(nil),  /* Retarget: */
	11708077258721206605: (*Reverse)(nil),   /* Reverse: */
	7416511141403176695:  (*Self)(nil),      /* Self: */
	1154838578286238320:  (*Words)(nil),     /* Words: */
}
