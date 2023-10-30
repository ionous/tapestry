package charm

import (
	"github.com/ionous/errutil"
)

// zero or more of the runes must pass the filter
func Optional(filter func(r rune) bool) State {
	return Self("optional", func(self State, r rune) (ret State) {
		if filter(r) {
			ret = self
		}
		return
	})
}

// one or more of the runes must pass the filter
func AtleastOne(filter func(r rune) bool) State {
	return Statement("require", func(r rune) (ret State) {
		if filter(r) {
			ret = Optional(filter)
		} else {
			e := errutil.New("unexpected rune")
			ret = Error(e)
		}
		return
	})
}
