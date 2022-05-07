package spec

import "strings"

//--------------------------
// TypeSpec:

// fix: this isnt right currently i dont think
func (op *TypeSpec) FriendlyName() string {
	return strings.Join(strings.Split(op.Name, "_"), " ")
}

//--------------------------
// TermSpec:

func (op *TermSpec) IsAnonymous() bool {
	return op.Label == "_"
}

// return the explicit parameter name, or the label if nothing explicit was defined.
func (op *TermSpec) ParameterName() (ret string) {
	if len(op.Name) > 0 {
		ret = op.Name
	} else {
		ret = op.Label
	}
	return
}

// return the name of the spec which describes the content of this term.
// ( either the explicitly specified type, or the parameter name if nothing explicitly declared. )
func (op *TermSpec) TypeName() (ret string) {
	if len(op.Type) > 0 {
		ret = op.Type
	} else {
		ret = op.ParameterName()
	}
	return
}

// tokenized field
func (op *TermSpec) Value() (ret string) {
	return "$" + strings.ToUpper(op.ParameterName())
}

//--------------------------
// ChoiceSpec:

// return english human friendly text.
func (op *ChoiceSpec) FriendlyName() (ret string) {
	if n := op.Label; len(n) > 0 {
		ret = n
	} else {
		ret = op.Name
	}
	return
}

// the value for the swap
func (op *ChoiceSpec) Value() string {
	return "$" + strings.ToUpper(op.Name)
}

// return the name of the spec which describes the content of this Choice.
func (op *ChoiceSpec) TypeName() (ret string) {
	if len(op.Type) > 0 {
		ret = op.Type
	} else {
		ret = op.Name
	}
	return
}

//--------------------------
// OptionSpec:

func (op *OptionSpec) FriendlyName() (ret string) {
	if n := op.Label; len(n) > 0 {
		ret = n
	} else {
		ret = op.Name
	}
	return
}

func (op *OptionSpec) Value() (ret string) {
	return "$" + strings.ToUpper(op.Name)
}
