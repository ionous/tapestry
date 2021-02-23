package internal

import (
	"log"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/parser"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
)

// note: play has to attach to some specifics of a particular pttime to work
// ex. childrenOf, transparentOf, etc.
// enc is the enclosureOf the player
func (pt *Playtime) GetDefaultBounds(enc string) parser.Bounds {
	return func(cb parser.NounVisitor) (ret bool) {
		// childrenOf enc
		// when done, ask transpentOf
	Exit:
		for len(enc) > 0 {
			// might be some filter here on room, etc.
			// maybe we want a more specialized pattern...

			if v, e := pt.Call("childrenOf", affine.TextList, []rt.Arg{{"obj",
				&core.FromValue{g.StringOf(enc)}},
			}); e != nil {
				// err = e
				log.Println(e)
			} else if done := cb(MakeNoun(pt, enc)); done {
				ret = true
				break Exit
			} else {
				// visit children
				for _, n := range v.Strings() {
					if done := cb(MakeNoun(pt, n)); done {
						ret = true
						break Exit
					}
				}
				// step up
				if v, e := pt.Call("transparentOf", affine.Text, []rt.Arg{{"obj",
					&core.FromValue{g.StringOf(enc)}},
				}); e != nil {
					// err = e
					log.Println(e)
				} else {
					enc = v.String() // i dont know that we have to report the enclosure itself...
				}
			}
		}
		return
	}
}
