// Code generated by "makeops"; edit at your own risk.
package render

import (
	"encoding/json"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/export/jsonexp"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

// RenderExp
type RenderExp struct {
	Expression rt.TextEval `if:"label=_"`
}

func (*RenderExp) Compose() composer.Spec {
	return composer.Spec{
		Name: Type_RenderExp,
		Uses: composer.Type_Flow,
	}
}

const Type_RenderExp = "render_exp"
const RenderExp_Expression = "$EXPRESSION"

func (op *RenderExp) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return RenderExp_Detailed_Marshal(n, op)
}
func (op *RenderExp) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return RenderExp_Detailed_Unmarshal(n, b, op)
}

func RenderExp_Detailed_Marshal(n jsonexp.Context, val *RenderExp) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := rt.TextEval_Detailed_Marshal(n, &val.Expression); e != nil {
		err = errutil.Append(err, e)
	} else if len(b) > 0 {
		fields[RenderExp_Expression] = b
	}
	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   Type_RenderExp,
			Fields: fields,
		})
	}
	return
}

func RenderExp_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *RenderExp) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(Type_RenderExp, "-", e)
	} else if e := rt.TextEval_Detailed_Unmarshal(n, msg.Fields[RenderExp_Expression], &out.Expression); e != nil {
		err = errutil.New(Type_RenderExp+"."+RenderExp_Expression, "-", e)
	}
	return
}

// RenderField in template phrases, picks between record variables, object variables, and named global objects.,ex. could be &quot;ringBearer&quot;, &quot;SamWise&quot;, or &quot;frodo&quot;
type RenderField struct {
	Name rt.TextEval `if:"label=_"`
}

func (*RenderField) Compose() composer.Spec {
	return composer.Spec{
		Name: Type_RenderField,
		Uses: composer.Type_Flow,
	}
}

const Type_RenderField = "render_field"
const RenderField_Name = "$NAME"

func (op *RenderField) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return RenderField_Detailed_Marshal(n, op)
}
func (op *RenderField) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return RenderField_Detailed_Unmarshal(n, b, op)
}

func RenderField_Detailed_Marshal(n jsonexp.Context, val *RenderField) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := rt.TextEval_Detailed_Marshal(n, &val.Name); e != nil {
		err = errutil.Append(err, e)
	} else if len(b) > 0 {
		fields[RenderField_Name] = b
	}
	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   Type_RenderField,
			Fields: fields,
		})
	}
	return
}

func RenderField_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *RenderField) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(Type_RenderField, "-", e)
	} else if e := rt.TextEval_Detailed_Unmarshal(n, msg.Fields[RenderField_Name], &out.Name); e != nil {
		err = errutil.New(Type_RenderField+"."+RenderField_Name, "-", e)
	}
	return
}

// RenderFlags requires a user-specified string.
type RenderFlags struct {
	Str string
}

func (op *RenderFlags) String() (ret string) {
	return op.Str
}

const RenderFlags_RenderAsVar = "$RENDER_AS_VAR"
const RenderFlags_RenderAsObj = "$RENDER_AS_OBJ"
const RenderFlags_RenderAsAny = "$RENDER_AS_ANY"

func (*RenderFlags) Compose() composer.Spec {
	return composer.Spec{
		Name: Type_RenderFlags,
		Uses: composer.Type_Str,
		Choices: []string{
			RenderFlags_RenderAsVar, RenderFlags_RenderAsObj, RenderFlags_RenderAsAny,
		},
		Strings: []string{
			"render_as_var", "render_as_obj", "render_as_any",
		},
	}
}

const Type_RenderFlags = "render_flags"

func (op *RenderFlags) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return RenderFlags_Detailed_Marshal(n, op)
}
func (op *RenderFlags) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return RenderFlags_Detailed_Unmarshal(n, b, op)
}
func RenderFlags_Detailed_Marshal(n jsonexp.Context, val *RenderFlags) ([]byte, error) {
	return json.Marshal(jsonexp.Str{
		Type:  Type_RenderFlags,
		Value: val.Str,
	})
}

func RenderFlags_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *RenderFlags) (err error) {
	var msg jsonexp.Str
	if len(b) > 0 { // generated code collapses optional and empty.
		if e := json.Unmarshal(b, &msg); e != nil {
			err = errutil.New(Type_RenderFlags, "-", e)
		}
	}
	if err == nil {
		out.Str = msg.Value
	}
	return
}

// RenderName handles changing a template like {.boombip} into text.,if the name is a variable containing an object name: return the printed object name ( via &quot;print name&quot; ),if the name is a variable with some other text: return that text.,if the name isn&#x27;t a variable but refers to some object: return that object&#x27;s printed object name.,otherwise, its an error.
type RenderName struct {
	Name string `if:"label=_,type=text"`
}

func (*RenderName) Compose() composer.Spec {
	return composer.Spec{
		Name: Type_RenderName,
		Uses: composer.Type_Flow,
	}
}

const Type_RenderName = "render_name"
const RenderName_Name = "$NAME"

func (op *RenderName) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return RenderName_Detailed_Marshal(n, op)
}
func (op *RenderName) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return RenderName_Detailed_Unmarshal(n, b, op)
}

func RenderName_Detailed_Marshal(n jsonexp.Context, val *RenderName) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := value.Text_Detailed_Override_Marshal(n, &val.Name); e != nil {
		err = errutil.Append(err, e)
	} else if len(b) > 0 {
		fields[RenderName_Name] = b
	}
	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   Type_RenderName,
			Fields: fields,
		})
	}
	return
}

func RenderName_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *RenderName) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(Type_RenderName, "-", e)
	} else if e := value.Text_Detailed_Override_Unmarshal(n, msg.Fields[RenderName_Name], &out.Name); e != nil {
		err = errutil.New(Type_RenderName+"."+RenderName_Name, "-", e)
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
		Name: Type_RenderPattern,
		Uses: composer.Type_Flow,
		Lede: "render",
	}
}

const Type_RenderPattern = "render_pattern"
const RenderPattern_Pattern = "$PATTERN"
const RenderPattern_Arguments = "$ARGUMENTS"

func (op *RenderPattern) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return RenderPattern_Detailed_Marshal(n, op)
}
func (op *RenderPattern) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return RenderPattern_Detailed_Unmarshal(n, b, op)
}

func RenderPattern_Detailed_Marshal(n jsonexp.Context, val *RenderPattern) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := value.PatternName_Detailed_Marshal(n, &val.Pattern); e != nil {
		err = errutil.Append(err, e)
	} else if len(b) > 0 {
		fields[RenderPattern_Pattern] = b
	}
	if b, e := core.CallArgs_Detailed_Marshal(n, &val.Arguments); e != nil {
		err = errutil.Append(err, e)
	} else if len(b) > 0 {
		fields[RenderPattern_Arguments] = b
	}
	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   Type_RenderPattern,
			Fields: fields,
		})
	}
	return
}

func RenderPattern_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *RenderPattern) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(Type_RenderPattern, "-", e)
	} else if e := value.PatternName_Detailed_Unmarshal(n, msg.Fields[RenderPattern_Pattern], &out.Pattern); e != nil {
		err = errutil.New(Type_RenderPattern+"."+RenderPattern_Pattern, "-", e)
	} else if e := core.CallArgs_Detailed_Unmarshal(n, msg.Fields[RenderPattern_Arguments], &out.Arguments); e != nil {
		err = errutil.New(Type_RenderPattern+"."+RenderPattern_Arguments, "-", e)
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
		Name: Type_RenderRef,
		Uses: composer.Type_Flow,
	}
}

const Type_RenderRef = "render_ref"
const RenderRef_Name = "$NAME"
const RenderRef_Flags = "$FLAGS"

func (op *RenderRef) MarshalDetailed(n jsonexp.Context) (ret []byte, err error) {
	return RenderRef_Detailed_Marshal(n, op)
}
func (op *RenderRef) UnmarshalDetailed(n jsonexp.Context, b []byte) error {
	return RenderRef_Detailed_Unmarshal(n, b, op)
}

func RenderRef_Detailed_Marshal(n jsonexp.Context, val *RenderRef) (ret []byte, err error) {
	fields := make(jsonexp.Fields)
	if b, e := value.VariableName_Detailed_Marshal(n, &val.Name); e != nil {
		err = errutil.Append(err, e)
	} else if len(b) > 0 {
		fields[RenderRef_Name] = b
	}
	if b, e := RenderFlags_Detailed_Marshal(n, &val.Flags); e != nil {
		err = errutil.Append(err, e)
	} else if len(b) > 0 {
		fields[RenderRef_Flags] = b
	}
	if err == nil {
		ret, err = json.Marshal(jsonexp.Flow{
			Type:   Type_RenderRef,
			Fields: fields,
		})
	}
	return
}

func RenderRef_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *RenderRef) (err error) {
	var msg jsonexp.Flow
	if e := json.Unmarshal(b, &msg); e != nil {
		err = errutil.New(Type_RenderRef, "-", e)
	} else if e := value.VariableName_Detailed_Unmarshal(n, msg.Fields[RenderRef_Name], &out.Name); e != nil {
		err = errutil.New(Type_RenderRef+"."+RenderRef_Name, "-", e)
	} else if e := RenderFlags_Detailed_Unmarshal(n, msg.Fields[RenderRef_Flags], &out.Flags); e != nil {
		err = errutil.New(Type_RenderRef+"."+RenderRef_Flags, "-", e)
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
