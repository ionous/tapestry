package internal

import (
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
		if done := cb(MakeNoun(pt, enc)); done {
			// we can handle this if need be
			panic("unexpected failure of ranking")
		} else {
		Exit:
			for bounds := []string{enc}; len(bounds) > 0; {
				enc, bounds = bounds[0], bounds[1:]
				if kids, e := pt.Call("childrenOf", affine.TextList, []rt.Arg{{"obj",
					&core.FromValue{g.StringOf(enc)}},
				}); e != nil {
					// err = e
					log.Println(e)
				} else {
					// visit children
					for _, kid := range kids.Strings() {
						if done := cb(MakeNoun(pt, kid)); done {
							// we can handle this if need be
							panic("unexpected failure of ranking")
							break Exit
						}

						// step down
						if v, e := pt.Call("fully_opaque", affine.Bool, []rt.Arg{{"obj",
							&core.FromValue{g.StringOf(kid)}},
						}); e != nil {
							// err = e
							log.Println(e)
						} else if !v.Bool() {
							bounds = append(bounds, kid)
						}
					}
				}
			}
		}
		return
	}
}
