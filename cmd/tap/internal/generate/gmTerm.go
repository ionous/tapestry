package generate

import "git.sr.ht/~ionous/tapestry/dl/spec"

// Wraps the existing spec Term to add helpers which make generating the go-code easier.
type Term struct {
	ctx  *Context
	Flow *spec.FlowSpec
	spec.TermSpec
}

func (t *Term) IsUnboxed() (okay bool) {
	_, okay = t.ctx.unbox[t.TypeName()]
	return
}

func (t *Term) GoDecl() string {
	var qualifier string
	typeName := t.TypeName()
	termType, ok := t.ctx.GetTypeSpec(typeName) // the referenced type.
	if t.Repeats {
		qualifier = "[]"
	} else if t.Optional && ok && termType.Spec.Choice == spec.UsesSpec_Flow_Opt {
		qualifier = "*" // pointer to a flow
	}
	return qualifier + t.ctx.scopedName(typeName, false)
}
