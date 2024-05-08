package generate

import (
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

var Pascal = inflect.Pascal
var Camelize = inflect.Camelize

// does the passed string list include the passed string?
func Includes(strs []string, str string) (ret bool) {
	for _, el := range strs {
		if el == str {
			ret = true
			break
		}
	}
	return
}
