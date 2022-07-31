package spec

import (
	"strings"
	"unicode"
)

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

//--------
