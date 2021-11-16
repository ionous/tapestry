// Code generated by "makeops"; edit at your own risk.
package eph

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

// EphAt
type EphAt struct {
	At  string   `if:"label=at,type=text"`
	Eph Ephemera `if:"label=eph"`
}

func (*EphAt) Compose() composer.Spec {
	return composer.Spec{
		Name: EphAt_Type,
		Uses: composer.Type_Flow,
		Lede: "eph",
	}
}

const EphAt_Type = "eph_at"

const EphAt_Field_At = "$AT"
const EphAt_Field_Eph = "$EPH"

func (op *EphAt) Marshal(m jsn.Marshaler) error {
	return EphAt_Marshal(m, op)
}

type EphAt_Slice []EphAt

func (op *EphAt_Slice) GetType() string { return EphAt_Type }

func (op *EphAt_Slice) Marshal(m jsn.Marshaler) error {
	return EphAt_Repeats_Marshal(m, (*[]EphAt)(op))
}

func (op *EphAt_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *EphAt_Slice) SetSize(cnt int) {
	var els []EphAt
	if cnt >= 0 {
		els = make(EphAt_Slice, cnt)
	}
	(*op) = els
}

func (op *EphAt_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return EphAt_Marshal(m, &(*op)[i])
}

func EphAt_Repeats_Marshal(m jsn.Marshaler, vals *[]EphAt) error {
	return jsn.RepeatBlock(m, (*EphAt_Slice)(vals))
}

func EphAt_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]EphAt) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = EphAt_Repeats_Marshal(m, pv)
	}
	return
}

type EphAt_Flow struct{ ptr *EphAt }

func (n EphAt_Flow) GetType() string      { return EphAt_Type }
func (n EphAt_Flow) GetLede() string      { return "eph" }
func (n EphAt_Flow) GetFlow() interface{} { return n.ptr }
func (n EphAt_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*EphAt); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func EphAt_Optional_Marshal(m jsn.Marshaler, pv **EphAt) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = EphAt_Marshal(m, *pv)
	} else if !enc {
		var v EphAt
		if err = EphAt_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func EphAt_Marshal(m jsn.Marshaler, val *EphAt) (err error) {
	if err = m.MarshalBlock(EphAt_Flow{val}); err == nil {
		e0 := m.MarshalKey("at", EphAt_Field_At)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.At)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", EphAt_Field_At))
		}
		e1 := m.MarshalKey("eph", EphAt_Field_Eph)
		if e1 == nil {
			e1 = Ephemera_Marshal(m, &val.Eph)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", EphAt_Field_Eph))
		}
		m.EndBlock()
	}
	return
}

// EphBeginDomain
// User implements: Ephemera.
type EphBeginDomain struct {
	Name     string   `if:"label=domain,type=text"`
	Requires []string `if:"label=requires,type=text"`
}

func (*EphBeginDomain) Compose() composer.Spec {
	return composer.Spec{
		Name: EphBeginDomain_Type,
		Uses: composer.Type_Flow,
		Lede: "eph",
	}
}

const EphBeginDomain_Type = "eph_begin_domain"

const EphBeginDomain_Field_Name = "$NAME"
const EphBeginDomain_Field_Requires = "$REQUIRES"

func (op *EphBeginDomain) Marshal(m jsn.Marshaler) error {
	return EphBeginDomain_Marshal(m, op)
}

type EphBeginDomain_Slice []EphBeginDomain

func (op *EphBeginDomain_Slice) GetType() string { return EphBeginDomain_Type }

func (op *EphBeginDomain_Slice) Marshal(m jsn.Marshaler) error {
	return EphBeginDomain_Repeats_Marshal(m, (*[]EphBeginDomain)(op))
}

func (op *EphBeginDomain_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *EphBeginDomain_Slice) SetSize(cnt int) {
	var els []EphBeginDomain
	if cnt >= 0 {
		els = make(EphBeginDomain_Slice, cnt)
	}
	(*op) = els
}

func (op *EphBeginDomain_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return EphBeginDomain_Marshal(m, &(*op)[i])
}

func EphBeginDomain_Repeats_Marshal(m jsn.Marshaler, vals *[]EphBeginDomain) error {
	return jsn.RepeatBlock(m, (*EphBeginDomain_Slice)(vals))
}

func EphBeginDomain_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]EphBeginDomain) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = EphBeginDomain_Repeats_Marshal(m, pv)
	}
	return
}

type EphBeginDomain_Flow struct{ ptr *EphBeginDomain }

func (n EphBeginDomain_Flow) GetType() string      { return EphBeginDomain_Type }
func (n EphBeginDomain_Flow) GetLede() string      { return "eph" }
func (n EphBeginDomain_Flow) GetFlow() interface{} { return n.ptr }
func (n EphBeginDomain_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*EphBeginDomain); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func EphBeginDomain_Optional_Marshal(m jsn.Marshaler, pv **EphBeginDomain) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = EphBeginDomain_Marshal(m, *pv)
	} else if !enc {
		var v EphBeginDomain
		if err = EphBeginDomain_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func EphBeginDomain_Marshal(m jsn.Marshaler, val *EphBeginDomain) (err error) {
	if err = m.MarshalBlock(EphBeginDomain_Flow{val}); err == nil {
		e0 := m.MarshalKey("domain", EphBeginDomain_Field_Name)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", EphBeginDomain_Field_Name))
		}
		e1 := m.MarshalKey("requires", EphBeginDomain_Field_Requires)
		if e1 == nil {
			e1 = value.Text_Unboxed_Repeats_Marshal(m, &val.Requires)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", EphBeginDomain_Field_Requires))
		}
		m.EndBlock()
	}
	return
}

// EphCheckPrint
// User implements: Ephemera.
type EphCheckPrint struct {
	Name string `if:"label=check,type=text"`
	Text string `if:"label=prints,type=text"`
}

func (*EphCheckPrint) Compose() composer.Spec {
	return composer.Spec{
		Name: EphCheckPrint_Type,
		Uses: composer.Type_Flow,
		Lede: "eph",
	}
}

const EphCheckPrint_Type = "eph_check_print"

const EphCheckPrint_Field_Name = "$NAME"
const EphCheckPrint_Field_Text = "$TEXT"

func (op *EphCheckPrint) Marshal(m jsn.Marshaler) error {
	return EphCheckPrint_Marshal(m, op)
}

type EphCheckPrint_Slice []EphCheckPrint

func (op *EphCheckPrint_Slice) GetType() string { return EphCheckPrint_Type }

func (op *EphCheckPrint_Slice) Marshal(m jsn.Marshaler) error {
	return EphCheckPrint_Repeats_Marshal(m, (*[]EphCheckPrint)(op))
}

func (op *EphCheckPrint_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *EphCheckPrint_Slice) SetSize(cnt int) {
	var els []EphCheckPrint
	if cnt >= 0 {
		els = make(EphCheckPrint_Slice, cnt)
	}
	(*op) = els
}

func (op *EphCheckPrint_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return EphCheckPrint_Marshal(m, &(*op)[i])
}

func EphCheckPrint_Repeats_Marshal(m jsn.Marshaler, vals *[]EphCheckPrint) error {
	return jsn.RepeatBlock(m, (*EphCheckPrint_Slice)(vals))
}

func EphCheckPrint_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]EphCheckPrint) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = EphCheckPrint_Repeats_Marshal(m, pv)
	}
	return
}

type EphCheckPrint_Flow struct{ ptr *EphCheckPrint }

func (n EphCheckPrint_Flow) GetType() string      { return EphCheckPrint_Type }
func (n EphCheckPrint_Flow) GetLede() string      { return "eph" }
func (n EphCheckPrint_Flow) GetFlow() interface{} { return n.ptr }
func (n EphCheckPrint_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*EphCheckPrint); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func EphCheckPrint_Optional_Marshal(m jsn.Marshaler, pv **EphCheckPrint) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = EphCheckPrint_Marshal(m, *pv)
	} else if !enc {
		var v EphCheckPrint
		if err = EphCheckPrint_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func EphCheckPrint_Marshal(m jsn.Marshaler, val *EphCheckPrint) (err error) {
	if err = m.MarshalBlock(EphCheckPrint_Flow{val}); err == nil {
		e0 := m.MarshalKey("check", EphCheckPrint_Field_Name)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", EphCheckPrint_Field_Name))
		}
		e1 := m.MarshalKey("prints", EphCheckPrint_Field_Text)
		if e1 == nil {
			e1 = value.Text_Unboxed_Marshal(m, &val.Text)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", EphCheckPrint_Field_Text))
		}
		m.EndBlock()
	}
	return
}

// EphEndDomain
// User implements: Ephemera.
type EphEndDomain struct {
	Name string `if:"label=domain,type=text"`
}

func (*EphEndDomain) Compose() composer.Spec {
	return composer.Spec{
		Name: EphEndDomain_Type,
		Uses: composer.Type_Flow,
		Lede: "eph",
	}
}

const EphEndDomain_Type = "eph_end_domain"

const EphEndDomain_Field_Name = "$NAME"

func (op *EphEndDomain) Marshal(m jsn.Marshaler) error {
	return EphEndDomain_Marshal(m, op)
}

type EphEndDomain_Slice []EphEndDomain

func (op *EphEndDomain_Slice) GetType() string { return EphEndDomain_Type }

func (op *EphEndDomain_Slice) Marshal(m jsn.Marshaler) error {
	return EphEndDomain_Repeats_Marshal(m, (*[]EphEndDomain)(op))
}

func (op *EphEndDomain_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *EphEndDomain_Slice) SetSize(cnt int) {
	var els []EphEndDomain
	if cnt >= 0 {
		els = make(EphEndDomain_Slice, cnt)
	}
	(*op) = els
}

func (op *EphEndDomain_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return EphEndDomain_Marshal(m, &(*op)[i])
}

func EphEndDomain_Repeats_Marshal(m jsn.Marshaler, vals *[]EphEndDomain) error {
	return jsn.RepeatBlock(m, (*EphEndDomain_Slice)(vals))
}

func EphEndDomain_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]EphEndDomain) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = EphEndDomain_Repeats_Marshal(m, pv)
	}
	return
}

type EphEndDomain_Flow struct{ ptr *EphEndDomain }

func (n EphEndDomain_Flow) GetType() string      { return EphEndDomain_Type }
func (n EphEndDomain_Flow) GetLede() string      { return "eph" }
func (n EphEndDomain_Flow) GetFlow() interface{} { return n.ptr }
func (n EphEndDomain_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*EphEndDomain); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func EphEndDomain_Optional_Marshal(m jsn.Marshaler, pv **EphEndDomain) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = EphEndDomain_Marshal(m, *pv)
	} else if !enc {
		var v EphEndDomain
		if err = EphEndDomain_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func EphEndDomain_Marshal(m jsn.Marshaler, val *EphEndDomain) (err error) {
	if err = m.MarshalBlock(EphEndDomain_Flow{val}); err == nil {
		e0 := m.MarshalKey("domain", EphEndDomain_Field_Name)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", EphEndDomain_Field_Name))
		}
		m.EndBlock()
	}
	return
}

// EphList
type EphList struct {
	All []EphAt `if:"label=list"`
}

func (*EphList) Compose() composer.Spec {
	return composer.Spec{
		Name: EphList_Type,
		Uses: composer.Type_Flow,
		Lede: "eph",
	}
}

const EphList_Type = "eph_list"

const EphList_Field_All = "$ALL"

func (op *EphList) Marshal(m jsn.Marshaler) error {
	return EphList_Marshal(m, op)
}

type EphList_Slice []EphList

func (op *EphList_Slice) GetType() string { return EphList_Type }

func (op *EphList_Slice) Marshal(m jsn.Marshaler) error {
	return EphList_Repeats_Marshal(m, (*[]EphList)(op))
}

func (op *EphList_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *EphList_Slice) SetSize(cnt int) {
	var els []EphList
	if cnt >= 0 {
		els = make(EphList_Slice, cnt)
	}
	(*op) = els
}

func (op *EphList_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return EphList_Marshal(m, &(*op)[i])
}

func EphList_Repeats_Marshal(m jsn.Marshaler, vals *[]EphList) error {
	return jsn.RepeatBlock(m, (*EphList_Slice)(vals))
}

func EphList_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]EphList) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = EphList_Repeats_Marshal(m, pv)
	}
	return
}

type EphList_Flow struct{ ptr *EphList }

func (n EphList_Flow) GetType() string      { return EphList_Type }
func (n EphList_Flow) GetLede() string      { return "eph" }
func (n EphList_Flow) GetFlow() interface{} { return n.ptr }
func (n EphList_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*EphList); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func EphList_Optional_Marshal(m jsn.Marshaler, pv **EphList) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = EphList_Marshal(m, *pv)
	} else if !enc {
		var v EphList
		if err = EphList_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func EphList_Marshal(m jsn.Marshaler, val *EphList) (err error) {
	if err = m.MarshalBlock(EphList_Flow{val}); err == nil {
		e0 := m.MarshalKey("list", EphList_Field_All)
		if e0 == nil {
			e0 = EphAt_Repeats_Marshal(m, &val.All)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", EphList_Field_All))
		}
		m.EndBlock()
	}
	return
}

// EphNameRef
// User implements: Ephemera.
type EphNameRef struct {
	Name string `if:"label=ref,type=text"`
	Type string `if:"label=of,type=text"`
}

func (*EphNameRef) Compose() composer.Spec {
	return composer.Spec{
		Name: EphNameRef_Type,
		Uses: composer.Type_Flow,
		Lede: "eph",
	}
}

const EphNameRef_Type = "eph_name_ref"

const EphNameRef_Field_Name = "$NAME"
const EphNameRef_Field_Type = "$TYPE"

func (op *EphNameRef) Marshal(m jsn.Marshaler) error {
	return EphNameRef_Marshal(m, op)
}

type EphNameRef_Slice []EphNameRef

func (op *EphNameRef_Slice) GetType() string { return EphNameRef_Type }

func (op *EphNameRef_Slice) Marshal(m jsn.Marshaler) error {
	return EphNameRef_Repeats_Marshal(m, (*[]EphNameRef)(op))
}

func (op *EphNameRef_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *EphNameRef_Slice) SetSize(cnt int) {
	var els []EphNameRef
	if cnt >= 0 {
		els = make(EphNameRef_Slice, cnt)
	}
	(*op) = els
}

func (op *EphNameRef_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return EphNameRef_Marshal(m, &(*op)[i])
}

func EphNameRef_Repeats_Marshal(m jsn.Marshaler, vals *[]EphNameRef) error {
	return jsn.RepeatBlock(m, (*EphNameRef_Slice)(vals))
}

func EphNameRef_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]EphNameRef) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = EphNameRef_Repeats_Marshal(m, pv)
	}
	return
}

type EphNameRef_Flow struct{ ptr *EphNameRef }

func (n EphNameRef_Flow) GetType() string      { return EphNameRef_Type }
func (n EphNameRef_Flow) GetLede() string      { return "eph" }
func (n EphNameRef_Flow) GetFlow() interface{} { return n.ptr }
func (n EphNameRef_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*EphNameRef); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func EphNameRef_Optional_Marshal(m jsn.Marshaler, pv **EphNameRef) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = EphNameRef_Marshal(m, *pv)
	} else if !enc {
		var v EphNameRef
		if err = EphNameRef_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func EphNameRef_Marshal(m jsn.Marshaler, val *EphNameRef) (err error) {
	if err = m.MarshalBlock(EphNameRef_Flow{val}); err == nil {
		e0 := m.MarshalKey("ref", EphNameRef_Field_Name)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", EphNameRef_Field_Name))
		}
		e1 := m.MarshalKey("of", EphNameRef_Field_Type)
		if e1 == nil {
			e1 = value.Text_Unboxed_Marshal(m, &val.Type)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", EphNameRef_Field_Type))
		}
		m.EndBlock()
	}
	return
}

// EphPlural plurals are needed at runtime to help parser what the user inputs.
// plurals are also needed at assembly time to understand what the author wrote.
// User implements: Ephemera.
type EphPlural struct {
	Plural   string `if:"label=plural,type=text"`
	Singular string `if:"label=singular,type=text"`
}

func (*EphPlural) Compose() composer.Spec {
	return composer.Spec{
		Name: EphPlural_Type,
		Uses: composer.Type_Flow,
		Lede: "eph",
	}
}

const EphPlural_Type = "eph_plural"

const EphPlural_Field_Plural = "$PLURAL"
const EphPlural_Field_Singular = "$SINGULAR"

func (op *EphPlural) Marshal(m jsn.Marshaler) error {
	return EphPlural_Marshal(m, op)
}

type EphPlural_Slice []EphPlural

func (op *EphPlural_Slice) GetType() string { return EphPlural_Type }

func (op *EphPlural_Slice) Marshal(m jsn.Marshaler) error {
	return EphPlural_Repeats_Marshal(m, (*[]EphPlural)(op))
}

func (op *EphPlural_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *EphPlural_Slice) SetSize(cnt int) {
	var els []EphPlural
	if cnt >= 0 {
		els = make(EphPlural_Slice, cnt)
	}
	(*op) = els
}

func (op *EphPlural_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return EphPlural_Marshal(m, &(*op)[i])
}

func EphPlural_Repeats_Marshal(m jsn.Marshaler, vals *[]EphPlural) error {
	return jsn.RepeatBlock(m, (*EphPlural_Slice)(vals))
}

func EphPlural_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]EphPlural) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = EphPlural_Repeats_Marshal(m, pv)
	}
	return
}

type EphPlural_Flow struct{ ptr *EphPlural }

func (n EphPlural_Flow) GetType() string      { return EphPlural_Type }
func (n EphPlural_Flow) GetLede() string      { return "eph" }
func (n EphPlural_Flow) GetFlow() interface{} { return n.ptr }
func (n EphPlural_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*EphPlural); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func EphPlural_Optional_Marshal(m jsn.Marshaler, pv **EphPlural) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = EphPlural_Marshal(m, *pv)
	} else if !enc {
		var v EphPlural
		if err = EphPlural_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func EphPlural_Marshal(m jsn.Marshaler, val *EphPlural) (err error) {
	if err = m.MarshalBlock(EphPlural_Flow{val}); err == nil {
		e0 := m.MarshalKey("plural", EphPlural_Field_Plural)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.Plural)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", EphPlural_Field_Plural))
		}
		e1 := m.MarshalKey("singular", EphPlural_Field_Singular)
		if e1 == nil {
			e1 = value.Text_Unboxed_Marshal(m, &val.Singular)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", EphPlural_Field_Singular))
		}
		m.EndBlock()
	}
	return
}

// EphRun
// User implements: Ephemera.
type EphRun struct {
	Run []rt.Execute `if:"label=run"`
}

func (*EphRun) Compose() composer.Spec {
	return composer.Spec{
		Name: EphRun_Type,
		Uses: composer.Type_Flow,
		Lede: "eph",
	}
}

const EphRun_Type = "eph_run"

const EphRun_Field_Run = "$RUN"

func (op *EphRun) Marshal(m jsn.Marshaler) error {
	return EphRun_Marshal(m, op)
}

type EphRun_Slice []EphRun

func (op *EphRun_Slice) GetType() string { return EphRun_Type }

func (op *EphRun_Slice) Marshal(m jsn.Marshaler) error {
	return EphRun_Repeats_Marshal(m, (*[]EphRun)(op))
}

func (op *EphRun_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *EphRun_Slice) SetSize(cnt int) {
	var els []EphRun
	if cnt >= 0 {
		els = make(EphRun_Slice, cnt)
	}
	(*op) = els
}

func (op *EphRun_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return EphRun_Marshal(m, &(*op)[i])
}

func EphRun_Repeats_Marshal(m jsn.Marshaler, vals *[]EphRun) error {
	return jsn.RepeatBlock(m, (*EphRun_Slice)(vals))
}

func EphRun_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]EphRun) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = EphRun_Repeats_Marshal(m, pv)
	}
	return
}

type EphRun_Flow struct{ ptr *EphRun }

func (n EphRun_Flow) GetType() string      { return EphRun_Type }
func (n EphRun_Flow) GetLede() string      { return "eph" }
func (n EphRun_Flow) GetFlow() interface{} { return n.ptr }
func (n EphRun_Flow) SetFlow(i interface{}) (okay bool) {
	if ptr, ok := i.(*EphRun); ok {
		*n.ptr, okay = *ptr, true
	}
	return
}

func EphRun_Optional_Marshal(m jsn.Marshaler, pv **EphRun) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = EphRun_Marshal(m, *pv)
	} else if !enc {
		var v EphRun
		if err = EphRun_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func EphRun_Marshal(m jsn.Marshaler, val *EphRun) (err error) {
	if err = m.MarshalBlock(EphRun_Flow{val}); err == nil {
		e0 := m.MarshalKey("run", EphRun_Field_Run)
		if e0 == nil {
			e0 = rt.Execute_Repeats_Marshal(m, &val.Run)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", EphRun_Field_Run))
		}
		m.EndBlock()
	}
	return
}

const Ephemera_Type = "ephemera"

var Ephemera_Optional_Marshal = Ephemera_Marshal

type Ephemera_Slot struct{ ptr *Ephemera }

func (at Ephemera_Slot) GetType() string              { return Ephemera_Type }
func (at Ephemera_Slot) GetSlot() (interface{}, bool) { return *at.ptr, *at.ptr != nil }
func (at Ephemera_Slot) SetSlot(v interface{}) (okay bool) {
	(*at.ptr), okay = v.(Ephemera)
	return
}

func Ephemera_Marshal(m jsn.Marshaler, ptr *Ephemera) (err error) {
	slot := Ephemera_Slot{ptr}
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

type Ephemera_Slice []Ephemera

func (op *Ephemera_Slice) GetType() string { return Ephemera_Type }

func (op *Ephemera_Slice) Marshal(m jsn.Marshaler) error {
	return Ephemera_Repeats_Marshal(m, (*[]Ephemera)(op))
}

func (op *Ephemera_Slice) GetSize() (ret int) {
	if els := *op; els != nil {
		ret = len(els)
	} else {
		ret = -1
	}
	return
}

func (op *Ephemera_Slice) SetSize(cnt int) {
	var els []Ephemera
	if cnt >= 0 {
		els = make(Ephemera_Slice, cnt)
	}
	(*op) = els
}

func (op *Ephemera_Slice) MarshalEl(m jsn.Marshaler, i int) error {
	return Ephemera_Marshal(m, &(*op)[i])
}

func Ephemera_Repeats_Marshal(m jsn.Marshaler, vals *[]Ephemera) error {
	return jsn.RepeatBlock(m, (*Ephemera_Slice)(vals))
}

func Ephemera_Optional_Repeats_Marshal(m jsn.Marshaler, pv *[]Ephemera) (err error) {
	if *pv != nil || !m.IsEncoding() {
		err = Ephemera_Repeats_Marshal(m, pv)
	}
	return
}

var Slots = []interface{}{
	(*Ephemera)(nil),
}

var Slats = []composer.Composer{
	(*EphAt)(nil),
	(*EphBeginDomain)(nil),
	(*EphCheckPrint)(nil),
	(*EphEndDomain)(nil),
	(*EphList)(nil),
	(*EphNameRef)(nil),
	(*EphPlural)(nil),
	(*EphRun)(nil),
}

var Signatures = map[uint64]interface{}{
	9182060341586636438:  (*EphAt)(nil),          /* Eph at:eph: */
	12209727080993772760: (*EphBeginDomain)(nil), /* Eph domain:requires: */
	18354563224792793196: (*EphCheckPrint)(nil),  /* Eph check:prints: */
	4379746949646135194:  (*EphEndDomain)(nil),   /* Eph domain: */
	11648725103497180078: (*EphList)(nil),        /* Eph list: */
	9956475014949920846:  (*EphNameRef)(nil),     /* Eph ref:of: */
	890409142408471553:   (*EphPlural)(nil),      /* Eph plural:singular: */
	4420716908411308437:  (*EphRun)(nil),         /* Eph run: */
}
