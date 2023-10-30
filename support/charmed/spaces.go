package charmed

import (
	"git.sr.ht/~ionous/tapestry/support/charm"
	"github.com/ionous/errutil"
)

// spaces eats whitespace
var OptionalSpaces = charm.Self("spaces", func(self charm.State, r rune) (ret charm.State) {
	if IsSpace(r) {
		ret = self
	}
	return
})

var RequiredSpaces = charm.Statement("spaces", func(r rune) (ret charm.State) {
	if IsSpace(r) {
		ret = OptionalSpaces
	} else {
		e := errutil.New("expected whitespace")
		ret = charm.Error(e)
	}
	return
})
