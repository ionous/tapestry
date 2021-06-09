package compact

import (
	r "reflect"

	"git.sr.ht/~ionous/iffy/dl/composer"
)

type swapper interface {
	composer.Composer
	Choices() (nameToType map[string]interface{})
}
type stringer interface {
	composer.Composer
	Choices() (tokenToValue map[string]string)
}

var strType = r.TypeOf((*stringer)(nil)).Elem()
var swapType = r.TypeOf((*swapper)(nil)).Elem()

// // translate a choice, typically a $TOKEN, to a value.
// // note: go-code doesnt currently have a way to find a string's label.
// func FindChoice(op strType, choice string) (ret string, found bool) {
// 	spec, keys := op.Compose(), op.Choices()
// 	if str, ok := keys[choice]; ok {
// 		ret, found = str, ok
// 	} else if !ok && spec.OpenStrings {
// 		ret = choice
// 	}
// 	return
// }
