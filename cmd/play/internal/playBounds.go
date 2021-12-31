package internal

import (
	"errors"
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// note: play has to attach to some specifics of a particular runtime to work
// ex. childrenOf, transparentOf, etc.
// enc is the enclosureOf the player
func (pt *Playtime) LocationBounded(enc string) parser.Bounds {
	return func(cb parser.NounVisitor) (ret bool) {
		if kids, e := pt.Call("parser_bounds", affine.TextList, []rt.Arg{{"obj",
			&core.FromValue{g.StringOf(enc)}},
		}); e != nil && !errors.Is(e, rt.NoResult{}) {
			log.Println(e)
		} else {
			for _, k := range kids.Strings() {
				if ok := cb(MakeNoun(pt, k)); ok {
					ret = ok
					break
				}
			}
		}
		return
	}
}

// return bounds which include only the player object and nothing else.
func (pt *Playtime) SelfBounded() parser.Bounds {
	return func(cb parser.NounVisitor) (ret bool) {
		return cb(&Noun{pt, pt.player})
	}
}
