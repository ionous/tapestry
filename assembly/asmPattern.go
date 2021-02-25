package assembly

import (
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// hopefully temporary.
// public because needed for qnaPat_test right now
type PatternFrag struct {
	Name   string
	Return string          // name of return field; empty if none ( could be an index but slightly safer this way )
	Labels []string        // one label for every parameter
	Locals []rt.Assignment // usually equal to the number of locals; or nil for testing.
	Fields []g.Field       // flat list of params and locals and an optional return
	Rules  []rt.Rule
}
