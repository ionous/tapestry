package spec

import "strings"

// fix: this isnt right currently i dont think
func (op *TypeSpec) FriendlyName() string {
	return strings.Join(strings.Split(op.Name, "_"), " ")
}

func (op *TermSpec) IsAnonymous() bool {
	return op.Key == "_"
}

// return english human friendly text.
func (op *TermSpec) FriendlyName() string {
	return op.Key
}

// return a programmatic name, unique in the flow: usually no spaces, etc.
func (op *TermSpec) Field() (ret string) {
	if len(op.Name) > 0 {
		ret = op.Name
	} else {
		ret = op.Key
	}
	return
}

// return the name of the spec which describes the content of this term.
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
