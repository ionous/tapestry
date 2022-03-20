package gomake

import "git.sr.ht/~ionous/tapestry/dl/spec"

// Wwraps the existing spec Term to add helpers which make generating the go-code easier.
type Term struct {
	ctx  *Context
	Flow *spec.FlowSpec
	spec.TermSpec
	publicIndex int // index of just the public terms, -1 if private
}

// return field used for the passed term's golang identifier
func (t *Term) GoId() (ret string) {
	if n := t.Name; len(n) > 0 {
		ret = n
	} else {
		ret = t.Key
	}
	return
}

func (t *Term) GoType() (ret string) {
	if typeName := t.Type; len(typeName) > 0 {
		ret = typeName
	} else if n := t.Name; len(n) > 0 {
		ret = n
	} else {
		ret = t.Key
	}
	return
}

func (t *Term) GoDecl() string {
	var qualifier string
	typeName := t.GoType()
	termType, ok := t.ctx.GetTypeSpec(typeName) // the referenced type.
	if t.Repeats {
		qualifier = "[]"
	} else if t.Optional && ok && termType.Spec.Choice == spec.UsesSpec_Flow_Opt {
		qualifier = "*" // pointer to a flow
	}
	return qualifier + t.ctx.scopedName(typeName, false)
}

func (t *Term) IsUnboxed() (okay bool) {
	_, okay = unbox[t.GoType()]
	return
}

func (t *Term) IsInline() (okay bool) {
	return t.Flow.Trim && t.publicIndex == 0
}
