package spec

import "strings"

func (op *OptionSpec) FriendlyName() (ret string) {
	return strings.Replace(op.Key(), "_", " ", -1)
}

func (op *OptionSpec) Key() (ret string) {
	if n := op.Label; len(n) > 0 {
		ret = n
	} else {
		ret = op.Name
	}
	return
}

func (op *OptionSpec) Value() string {
	return "$" + strings.ToUpper(op.Name)
}
