// Code generated by "makeops"; edit at your own risk.
package render

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/jsn"
	"git.sr.ht/~ionous/iffy/rt"
)

// RenderExp
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

func (op *RenderExp) Marshal(n jsn.Marshaler) {
	RenderExp_Marshal(n, op)
}

type RenderExp_Slice []RenderExp

func (op *RenderExp_Slice) GetSize() int    { return len(*op) }
func (op *RenderExp_Slice) SetSize(cnt int) { (*op) = make(RenderExp_Slice, cnt) }

func RenderExp_Repeats_Marshal(n jsn.Marshaler, vals *[]RenderExp) {
	if n.RepeatValues(RenderExp_Type, (*RenderExp_Slice)(vals)) {
		for i := range *vals {
			RenderExp_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func RenderExp_Optional_Marshal(n jsn.Marshaler, pv **RenderExp) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		RenderExp_Marshal(n, *pv)
	} else if !enc {
		var v RenderExp
		if RenderExp_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func RenderExp_Marshal(n jsn.Marshaler, val *RenderExp) (okay bool) {
	if okay = n.MapValues(RenderExp_Type, RenderExp_Type); okay {
		if n.MapKey("", RenderExp_Field_Expression) {
			rt.TextEval_Marshal(n, &val.Expression)
		}
		n.EndValues()
	}
	return
}

// RenderField in template phrases, picks between record variables, object variables, and named global objects.,ex. could be &quot;ringBearer&quot;, &quot;SamWise&quot;, or &quot;frodo&quot;
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

func (op *RenderField) Marshal(n jsn.Marshaler) {
	RenderField_Marshal(n, op)
}

type RenderField_Slice []RenderField

func (op *RenderField_Slice) GetSize() int    { return len(*op) }
func (op *RenderField_Slice) SetSize(cnt int) { (*op) = make(RenderField_Slice, cnt) }

func RenderField_Repeats_Marshal(n jsn.Marshaler, vals *[]RenderField) {
	if n.RepeatValues(RenderField_Type, (*RenderField_Slice)(vals)) {
		for i := range *vals {
			RenderField_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func RenderField_Optional_Marshal(n jsn.Marshaler, pv **RenderField) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		RenderField_Marshal(n, *pv)
	} else if !enc {
		var v RenderField
		if RenderField_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func RenderField_Marshal(n jsn.Marshaler, val *RenderField) (okay bool) {
	if okay = n.MapValues(RenderField_Type, RenderField_Type); okay {
		if n.MapKey("", RenderField_Field_Name) {
			rt.TextEval_Marshal(n, &val.Name)
		}
		n.EndValues()
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

func (op *RenderFlags) Marshal(n jsn.Marshaler) {
	RenderFlags_Marshal(n, op)
}

func RenderFlags_Optional_Marshal(n jsn.Marshaler, val *RenderFlags) {
	var zero RenderFlags
	if enc := n.IsEncoding(); !enc || val.Str != zero.Str {
		RenderFlags_Marshal(n, val)
	}
}

func RenderFlags_Marshal(n jsn.Marshaler, val *RenderFlags) {
	n.MarshalValue(RenderFlags_Type, jsn.MakeEnum(val, &val.Str))
}

type RenderFlags_Slice []RenderFlags

func (op *RenderFlags_Slice) GetSize() int    { return len(*op) }
func (op *RenderFlags_Slice) SetSize(cnt int) { (*op) = make(RenderFlags_Slice, cnt) }

func RenderFlags_Repeats_Marshal(n jsn.Marshaler, vals *[]RenderFlags) {
	if n.RepeatValues(RenderFlags_Type, (*RenderFlags_Slice)(vals)) {
		for i := range *vals {
			RenderFlags_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

// RenderName handles changing a template like {.boombip} into text.,if the name is a variable containing an object name: return the printed object name ( via &quot;print name&quot; ),if the name is a variable with some other text: return that text.,if the name isn&#x27;t a variable but refers to some object: return that object&#x27;s printed object name.,otherwise, its an error.
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

func (op *RenderName) Marshal(n jsn.Marshaler) {
	RenderName_Marshal(n, op)
}

type RenderName_Slice []RenderName

func (op *RenderName_Slice) GetSize() int    { return len(*op) }
func (op *RenderName_Slice) SetSize(cnt int) { (*op) = make(RenderName_Slice, cnt) }

func RenderName_Repeats_Marshal(n jsn.Marshaler, vals *[]RenderName) {
	if n.RepeatValues(RenderName_Type, (*RenderName_Slice)(vals)) {
		for i := range *vals {
			RenderName_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func RenderName_Optional_Marshal(n jsn.Marshaler, pv **RenderName) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		RenderName_Marshal(n, *pv)
	} else if !enc {
		var v RenderName
		if RenderName_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func RenderName_Marshal(n jsn.Marshaler, val *RenderName) (okay bool) {
	if okay = n.MapValues(RenderName_Type, RenderName_Type); okay {
		if n.MapKey("", RenderName_Field_Name) {
			value.Text_Unboxed_Marshal(n, &val.Name)
		}
		n.EndValues()
	}
	return
}

// RenderPattern printing is generally an activity b/c say is an activity,and we want the ability to say several things in series.
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

func (op *RenderPattern) Marshal(n jsn.Marshaler) {
	RenderPattern_Marshal(n, op)
}

type RenderPattern_Slice []RenderPattern

func (op *RenderPattern_Slice) GetSize() int    { return len(*op) }
func (op *RenderPattern_Slice) SetSize(cnt int) { (*op) = make(RenderPattern_Slice, cnt) }

func RenderPattern_Repeats_Marshal(n jsn.Marshaler, vals *[]RenderPattern) {
	if n.RepeatValues(RenderPattern_Type, (*RenderPattern_Slice)(vals)) {
		for i := range *vals {
			RenderPattern_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func RenderPattern_Optional_Marshal(n jsn.Marshaler, pv **RenderPattern) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		RenderPattern_Marshal(n, *pv)
	} else if !enc {
		var v RenderPattern
		if RenderPattern_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func RenderPattern_Marshal(n jsn.Marshaler, val *RenderPattern) (okay bool) {
	if okay = n.MapValues("render", RenderPattern_Type); okay {
		if n.MapKey("", RenderPattern_Field_Pattern) {
			value.PatternName_Marshal(n, &val.Pattern)
		}
		if n.MapKey("args", RenderPattern_Field_Arguments) {
			core.CallArgs_Marshal(n, &val.Arguments)
		}
		n.EndValues()
	}
	return
}

// RenderRef returns the value of a variable or the id of an object.
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

func (op *RenderRef) Marshal(n jsn.Marshaler) {
	RenderRef_Marshal(n, op)
}

type RenderRef_Slice []RenderRef

func (op *RenderRef_Slice) GetSize() int    { return len(*op) }
func (op *RenderRef_Slice) SetSize(cnt int) { (*op) = make(RenderRef_Slice, cnt) }

func RenderRef_Repeats_Marshal(n jsn.Marshaler, vals *[]RenderRef) {
	if n.RepeatValues(RenderRef_Type, (*RenderRef_Slice)(vals)) {
		for i := range *vals {
			RenderRef_Marshal(n, &(*vals)[i])
		}
		n.EndValues()
	}
	return
}

func RenderRef_Optional_Marshal(n jsn.Marshaler, pv **RenderRef) {
	if enc := n.IsEncoding(); enc && *pv != nil {
		RenderRef_Marshal(n, *pv)
	} else if !enc {
		var v RenderRef
		if RenderRef_Marshal(n, &v) {
			*pv = &v
		}
	}
}

func RenderRef_Marshal(n jsn.Marshaler, val *RenderRef) (okay bool) {
	if okay = n.MapValues(RenderRef_Type, RenderRef_Type); okay {
		if n.MapKey("", RenderRef_Field_Name) {
			value.VariableName_Marshal(n, &val.Name)
		}
		if n.MapKey("flags", RenderRef_Field_Flags) {
			RenderFlags_Marshal(n, &val.Flags)
		}
		n.EndValues()
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
