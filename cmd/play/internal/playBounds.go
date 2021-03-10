package internal

import (
	"errors"
	"log"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/parser"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// note: play has to attach to some specifics of a particular runtime to work
// ex. childrenOf, transparentOf, etc.
// enc is the enclosureOf the player
func (pt *Playtime) GetNamedBounds(enc string) parser.Bounds {
	return func(cb parser.NounVisitor) (ret bool) {
		if kids, e := pt.Call("parser_bounds", affine.TextList, []rt.Arg{{"obj",
			&core.FromValue{g.StringOf(enc)}},
		}); e != nil && !errors.Is(e, rt.NoResult{}) {
			log.Println(e)
		} else {
			for _, k := range kids.Strings() {
				cb(MakeNoun(pt, k))
			}
		}
		return
	}
}
