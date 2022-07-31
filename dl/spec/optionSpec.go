package spec

import "strings"

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
