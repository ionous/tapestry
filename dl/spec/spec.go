package spec

import (
	"strings"
	"unicode"
)

//--------------------------
// TypeSpec:

func (op *TypeSpec) InGroup(g string) (okay bool) {
	for _, n := range op.Groups {
		if n == g {
			okay = true
			break
		}
	}
	return
}

// change one_two into "One two"
func (op *TypeSpec) FriendlyName() (ret string) {
	if flow, ok := op.Spec.Value.(*FlowSpec); ok {
		ret = flow.FriendlyLede(op)
	} else {
		ret = FriendlyName(op.Name, false)
	}
	return
}

func FriendlyName(n string, capAllWords bool) (ret string) {
	var b strings.Builder
	cap := true
	for i, s := range strings.Split(n, "_") {
		if i > 0 {
			b.WriteRune(' ')
		}
		if !cap {
			b.WriteString(s)
		} else {
			for j, r := range s {
				if j == 0 {
					b.WriteRune(unicode.ToUpper(r))
				} else {
					b.WriteRune(r)
				}
			}
			cap = capAllWords
		}
	}
	return b.String()
}

//--------------------------
// TermSpec:

func (op *TermSpec) IsAnonymous() bool {
	return op.Label == "_"
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
func (op *FlowSpec) FriendlyLede(blockType *TypeSpec) (ret string) {
	if lede := op.Name; len(lede) > 0 {
		ret = FriendlyName(lede, false)
	} else {
		ret = FriendlyName(blockType.Name, false)
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
