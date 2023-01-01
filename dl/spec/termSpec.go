package spec

import "strings"

func (op *TermSpec) IsAnonymous() bool {
	return op.Label == "_"
}

// returns the empty string if anonymous
func (op *TermSpec) QuietLabel() (ret string) {
	if !op.IsAnonymous() {
		ret = op.Label
	}
	return
}

// return the explicit parameter name, or the label if nothing explicit was defined.
func (op *TermSpec) Field() (ret string) {
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
		ret = op.Field()
	}
	return
}

// tokenized field
func (op *TermSpec) Value() (ret string) {
	return "$" + strings.ToUpper(op.Field())
}
