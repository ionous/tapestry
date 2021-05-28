// Code generated by "makeops"; edit at your own risk.
package render

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/rt"
)

// RenderAsAny
type RenderAsAny struct {
}

func (*RenderAsAny) Compose() composer.Spec {
	return composer.Spec{
		Name: "render_as_any",
		Lede: "as_any",
	}
}

// RenderAsObj
type RenderAsObj struct {
}

func (*RenderAsObj) Compose() composer.Spec {
	return composer.Spec{
		Name: "render_as_obj",
		Lede: "as_obj",
	}
}

// RenderAsVar
type RenderAsVar struct {
}

func (*RenderAsVar) Compose() composer.Spec {
	return composer.Spec{
		Name: "render_as_var",
		Lede: "as_var",
	}
}

// RenderField in template phrases, picks between record variables, object variables, and named global objects.,ex. could be &quot;ringBearer&quot;, &quot;SamWise&quot;, or &quot;frodo&quot;
type RenderField struct {
	Name rt.TextEval `if:"label=_"`
}

func (*RenderField) Compose() composer.Spec {
	return composer.Spec{
		Name: "render_field",
	}
}

// RenderFlags swaps between various options
type RenderFlags struct {
	Opt interface{}
}

func (*RenderFlags) Compose() composer.Spec {
	return composer.Spec{
		Name: "render_flags",
	}
}

func (*RenderFlags) Choices() map[string]interface{} {
	return map[string]interface{}{
		"render_as_var": (*RenderAsVar)(nil),
		"render_as_obj": (*RenderAsObj)(nil),
		"render_as_any": (*RenderAsAny)(nil),
	}
}

// RenderName handles changing a template like {.boombip} into text.,if the name is a variable containing an object name: return the printed object name ( via &quot;print name&quot; ),if the name is a variable with some other text: return that text.,if the name isn&#x27;t a variable but refers to some object: return that object&#x27;s printed object name.,otherwise, its an error.
type RenderName struct {
	Name string `if:"label=_"`
}

func (*RenderName) Compose() composer.Spec {
	return composer.Spec{
		Name: "render_name",
	}
}

// RenderPattern printing is generally an activity b/c say is an activity,and we want the ability to say several things in series.
type RenderPattern struct {
	Pattern   string         `if:"label=_"`
	Arguments core.Arguments `if:"label=arguments"`
}

func (*RenderPattern) Compose() composer.Spec {
	return composer.Spec{
		Name: "render_pattern",
		Lede: "render",
	}
}

// RenderRef returns the value of a variable or the id of an object.
type RenderRef struct {
	Name  string      `if:"label=_"`
	Flags RenderFlags `if:"label=flags"`
}

func (*RenderRef) Compose() composer.Spec {
	return composer.Spec{
		Name: "render_ref",
	}
}

// RenderTemplate Parse text using iffy templates. See: https://github.com/ionous/iffy/wiki/Templates
type RenderTemplate struct {
	Expression rt.TextEval `if:"label=_"`
}

func (*RenderTemplate) Compose() composer.Spec {
	return composer.Spec{
		Name: "render_template",
	}
}

var Swaps = []interface{}{
	(*RenderFlags)(nil),
}
var Slats = []interface{}{
	(*RenderAsAny)(nil),
	(*RenderAsObj)(nil),
	(*RenderAsVar)(nil),
	(*RenderField)(nil),
	(*RenderName)(nil),
	(*RenderPattern)(nil),
	(*RenderRef)(nil),
	(*RenderTemplate)(nil),
}
