package spec

import "strings"

func (op *OptionSpec) FriendlyName() (ret string) {
	if n := op.Label; len(n) > 0 {
		ret = n
	} else {
		ret = op.Name
	}
	return strings.Replace(ret, "_", " ", -1)
}

func (op *OptionSpec) Value() (ret string) {
	return "$" + strings.ToUpper(op.Name)
}
