// Code generated by "makeops"; edit at your own risk.
package render

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

// RenderExp
// User implements: TextEval.
type RenderExp struct {
	Expression rt.TextEval `if:"label=_"`
}

func (*RenderExp) Compose() composer.Spec {
	return composer.Spec{
		Name: RenderExp_Type,
		Uses: composer.Type_Flow,
	}
}

const RenderExp_Type = "render_exp"

const RenderExp_Field_Expression = "$EXPRESSION"

func (op *RenderExp) Marshal(m jsn.Marshaler) error {
	return RenderExp_Marshal(m, op)
}

type RenderExp_Slice []RenderExp

func (op *RenderExp_Slice) GetType() string { return RenderExp_Type }
func (op *RenderExp_Slice) GetSize() int    { return len(*op) }
func (op *RenderExp_Slice) SetSize(cnt int) { (*op) = make(RenderExp_Slice, cnt) }

func RenderExp_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderExp) (err error) {
	if err = m.MarshalBlock((*RenderExp_Slice)(vals)); err == nil {
		for i := range *vals {
			if e := RenderExp_Marshal(m, &(*vals)[i]); e != nil && e != jsn.Missing {
				m.Error(errutil.New(e, "in slice at", i))
			}
		}
		m.EndBlock()
	}
	return
}

func RenderExp_Optional_Marshal(m jsn.Marshaler, pv **RenderExp) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RenderExp_Marshal(m, *pv)
	} else if !enc {
		var v RenderExp
		if err = RenderExp_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RenderExp_Marshal(m jsn.Marshaler, val *RenderExp) (err error) {
	if err = m.MarshalBlock(jsn.MakeFlow(RenderExp_Type, RenderExp_Type, val)); err == nil {
		e0 := m.MarshalKey("", RenderExp_Field_Expression)
		if e0 == nil {
			e0 = rt.TextEval_Marshal(m, &val.Expression)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RenderExp_Field_Expression))
		}
		m.EndBlock()
	}
	return
}

// RenderField in template phrases, picks between record variables, object variables, and named global objects.,ex. could be &quot;ringBearer&quot;, &quot;SamWise&quot;, or &quot;frodo&quot;
// User implements: FromSourceFields.
type RenderField struct {
	Name rt.TextEval `if:"label=_"`
}

func (*RenderField) Compose() composer.Spec {
	return composer.Spec{
		Name: RenderField_Type,
		Uses: composer.Type_Flow,
	}
}

const RenderField_Type = "render_field"

const RenderField_Field_Name = "$NAME"

func (op *RenderField) Marshal(m jsn.Marshaler) error {
	return RenderField_Marshal(m, op)
}

type RenderField_Slice []RenderField

func (op *RenderField_Slice) GetType() string { return RenderField_Type }
func (op *RenderField_Slice) GetSize() int    { return len(*op) }
func (op *RenderField_Slice) SetSize(cnt int) { (*op) = make(RenderField_Slice, cnt) }

func RenderField_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderField) (err error) {
	if err = m.MarshalBlock((*RenderField_Slice)(vals)); err == nil {
		for i := range *vals {
			if e := RenderField_Marshal(m, &(*vals)[i]); e != nil && e != jsn.Missing {
				m.Error(errutil.New(e, "in slice at", i))
			}
		}
		m.EndBlock()
	}
	return
}

func RenderField_Optional_Marshal(m jsn.Marshaler, pv **RenderField) (err error) {
	if enc := m.IsEncoding(); enc && *pv != nil {
		err = RenderField_Marshal(m, *pv)
	} else if !enc {
		var v RenderField
		if err = RenderField_Marshal(m, &v); err == nil {
			*pv = &v
		}
	}
	return
}

func RenderField_Marshal(m jsn.Marshaler, val *RenderField) (err error) {
	if err = m.MarshalBlock(jsn.MakeFlow(RenderField_Type, RenderField_Type, val)); err == nil {
		e0 := m.MarshalKey("", RenderField_Field_Name)
		if e0 == nil {
			e0 = rt.TextEval_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RenderField_Field_Name))
		}
		m.EndBlock()
	}
	return
}

// RenderFlags requires a user-specified string.
type RenderFlags struct {
	Str string
}

func (op *RenderFlags) String() string {
	return op.Str
}

const RenderFlags_RenderAsVar = "$RENDER_AS_VAR"
const RenderFlags_RenderAsObj = "$RENDER_AS_OBJ"
const RenderFlags_RenderAsAny = "$RENDER_AS_ANY"

func (*RenderFlags) Compose() composer.Spec {
	return composer.Spec{
		Name: RenderFlags_Type,
		Uses: composer.Type_Str,
		Choices: []string{
			RenderFlags_RenderAsVar, RenderFlags_RenderAsObj, RenderFlags_RenderAsAny,
		},
		Strings: []string{
			"render_as_var", "render_as_obj", "render_as_any",
		},
	}
}

const RenderFlags_Type = "render_flags"

func (op *RenderFlags) Marshal(m jsn.Marshaler) error {
	return RenderFlags_Marshal(m, op)
}

func RenderFlags_Optional_Marshal(m jsn.Marshaler, val *RenderFlags) (err error) {
	var zero RenderFlags
	if enc := m.IsEncoding(); !enc || val.Str != zero.Str {
		err = RenderFlags_Marshal(m, val)
	}
	return
}

func RenderFlags_Marshal(m jsn.Marshaler, val *RenderFlags) (err error) {
	return m.MarshalValue(RenderFlags_Type, jsn.MakeEnum(val, &val.Str))
}

type RenderFlags_Slice []RenderFlags

func (op *RenderFlags_Slice) GetType() string { return RenderFlags_Type }
func (op *RenderFlags_Slice) GetSize() int    { return len(*op) }
func (op *RenderFlags_Slice) SetSize(cnt int) { (*op) = make(RenderFlags_Slice, cnt) }

func RenderFlags_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderFlags) (err error) {
	if err = m.MarshalBlock((*RenderFlags_Slice)(vals)); err == nil {
		for i := range *vals {
			if e := RenderFlags_Marshal(m, &(*vals)[i]); e != nil && e != jsn.Missing {
				m.Error(errutil.New(e, "in slice at", i))
			}
		}
		m.EndBlock()
	}
	return
}

// RenderName handles changing a template like {.boombip} into text.,if the name is a variable containing an object name: return the printed object name ( via &quot;print name&quot; ),if the name is a variable with some other text: return that text.,if the name isn&#x27;t a variable but refers to some object: return that object&#x27;s printed object name.,otherwise, its an error.
// User implements: TextEval.
type RenderName struct {
	Name string `if:"label=_,type=text"`
}

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
func (op *RenderName_Slice) GetSize() int    { return len(*op) }
func (op *RenderName_Slice) SetSize(cnt int) { (*op) = make(RenderName_Slice, cnt) }

func RenderName_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderName) (err error) {
	if err = m.MarshalBlock((*RenderName_Slice)(vals)); err == nil {
		for i := range *vals {
			if e := RenderName_Marshal(m, &(*vals)[i]); e != nil && e != jsn.Missing {
				m.Error(errutil.New(e, "in slice at", i))
			}
		}
		m.EndBlock()
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
	if err = m.MarshalBlock(jsn.MakeFlow(RenderName_Type, RenderName_Type, val)); err == nil {
		e0 := m.MarshalKey("", RenderName_Field_Name)
		if e0 == nil {
			e0 = value.Text_Unboxed_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RenderName_Field_Name))
		}
		m.EndBlock()
	}
	return
}

// RenderPattern printing is generally an activity b/c say is an activity,and we want the ability to say several things in series.
// User implements: Assignment, TextEval.
type RenderPattern struct {
	Pattern   value.PatternName `if:"label=_"`
	Arguments core.CallArgs     `if:"label=args"`
}

func (*RenderPattern) Compose() composer.Spec {
	return composer.Spec{
		Name: RenderPattern_Type,
		Uses: composer.Type_Flow,
		Lede: "render",
	}
}

const RenderPattern_Type = "render_pattern"

const RenderPattern_Field_Pattern = "$PATTERN"
const RenderPattern_Field_Arguments = "$ARGUMENTS"

func (op *RenderPattern) Marshal(m jsn.Marshaler) error {
	return RenderPattern_Marshal(m, op)
}

type RenderPattern_Slice []RenderPattern

func (op *RenderPattern_Slice) GetType() string { return RenderPattern_Type }
func (op *RenderPattern_Slice) GetSize() int    { return len(*op) }
func (op *RenderPattern_Slice) SetSize(cnt int) { (*op) = make(RenderPattern_Slice, cnt) }

func RenderPattern_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderPattern) (err error) {
	if err = m.MarshalBlock((*RenderPattern_Slice)(vals)); err == nil {
		for i := range *vals {
			if e := RenderPattern_Marshal(m, &(*vals)[i]); e != nil && e != jsn.Missing {
				m.Error(errutil.New(e, "in slice at", i))
			}
		}
		m.EndBlock()
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
	if err = m.MarshalBlock(jsn.MakeFlow("render", RenderPattern_Type, val)); err == nil {
		e0 := m.MarshalKey("", RenderPattern_Field_Pattern)
		if e0 == nil {
			e0 = value.PatternName_Marshal(m, &val.Pattern)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RenderPattern_Field_Pattern))
		}
		e1 := m.MarshalKey("args", RenderPattern_Field_Arguments)
		if e1 == nil {
			e1 = core.CallArgs_Marshal(m, &val.Arguments)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", RenderPattern_Field_Arguments))
		}
		m.EndBlock()
	}
	return
}

// RenderRef returns the value of a variable or the id of an object.
// User implements: Assignment, NumberEval, TextEval.
type RenderRef struct {
	Name  value.VariableName `if:"label=_"`
	Flags RenderFlags        `if:"label=flags"`
}

func (*RenderRef) Compose() composer.Spec {
	return composer.Spec{
		Name: RenderRef_Type,
		Uses: composer.Type_Flow,
	}
}

const RenderRef_Type = "render_ref"

const RenderRef_Field_Name = "$NAME"
const RenderRef_Field_Flags = "$FLAGS"

func (op *RenderRef) Marshal(m jsn.Marshaler) error {
	return RenderRef_Marshal(m, op)
}

type RenderRef_Slice []RenderRef

func (op *RenderRef_Slice) GetType() string { return RenderRef_Type }
func (op *RenderRef_Slice) GetSize() int    { return len(*op) }
func (op *RenderRef_Slice) SetSize(cnt int) { (*op) = make(RenderRef_Slice, cnt) }

func RenderRef_Repeats_Marshal(m jsn.Marshaler, vals *[]RenderRef) (err error) {
	if err = m.MarshalBlock((*RenderRef_Slice)(vals)); err == nil {
		for i := range *vals {
			if e := RenderRef_Marshal(m, &(*vals)[i]); e != nil && e != jsn.Missing {
				m.Error(errutil.New(e, "in slice at", i))
			}
		}
		m.EndBlock()
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
	if err = m.MarshalBlock(jsn.MakeFlow(RenderRef_Type, RenderRef_Type, val)); err == nil {
		e0 := m.MarshalKey("", RenderRef_Field_Name)
		if e0 == nil {
			e0 = value.VariableName_Marshal(m, &val.Name)
		}
		if e0 != nil && e0 != jsn.Missing {
			m.Error(errutil.New(e0, "in flow at", RenderRef_Field_Name))
		}
		e1 := m.MarshalKey("flags", RenderRef_Field_Flags)
		if e1 == nil {
			e1 = RenderFlags_Marshal(m, &val.Flags)
		}
		if e1 != nil && e1 != jsn.Missing {
			m.Error(errutil.New(e1, "in flow at", RenderRef_Field_Flags))
		}
		m.EndBlock()
	}
	return
}

var Slats = []composer.Composer{
	(*RenderExp)(nil),
	(*RenderField)(nil),
	(*RenderFlags)(nil),
	(*RenderName)(nil),
	(*RenderPattern)(nil),
	(*RenderRef)(nil),
}

var Signatures = map[uint64]interface{}{
	16799527360025986462: (*RenderExp)(nil),     /* RenderExp: */
	8103562808853847007:  (*RenderField)(nil),   /* RenderField: */
	2017102261165852124:  (*RenderName)(nil),    /* RenderName: */
	9758431868100851810:  (*RenderPattern)(nil), /* Render:args: */
	615784906923963755:   (*RenderRef)(nil),     /* RenderRef:flags: */
}
