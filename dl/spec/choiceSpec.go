package spec

import "strings"

// return english human friendly text.
func (op *ChoiceSpec) FriendlyName() (ret string) {
	if n := op.Label; len(n) > 0 {
		ret = n
	} else {
		ret = op.Name
	}
	return strings.Replace(ret, "_", " ", -1)
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
