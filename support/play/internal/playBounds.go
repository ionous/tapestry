package internal

import (
	"errors"
	"log"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/rt"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

func (pt *Playtime) getPawn() (ret string, err error) {
	// probably could cache the player agent object
	if pawn, e := pt.GetField(pt.player, "pawn"); e != nil {
		err = e
	} else {
		ret = pawn.String()
	}
	return
}

// note: play has to attach to some specifics of a particular runtime to work
// ex. childrenOf, transparentOf, etc.
// enc is the enclosureOf the player
func (pt *Playtime) locationBounded(enc string) parser.Bounds {
	return func(cb parser.NounVisitor) (ret bool) {
		if kids, e := pt.Call(
			pt.bounds,
			affine.TextList,
			[]string{"obj"},
			[]g.Value{g.StringOf(enc)},
		); e != nil && !errors.Is(e, rt.NoResult) {
			log.Println(e)
		} else {
			ret = pt.visitStrings(cb, kids)
		}
		return
	}
}

// return bounds which includes only the player agent and nothing else.
// fix: why this is the agent and not the pawn?
func (pt *Playtime) selfBounded() (ret parser.Bounds, err error) {
	ret = func(cb parser.NounVisitor) (ret bool) {
		return cb(&Noun{pt, pt.player})
	}
	return
}

func (pt *Playtime) visitStrings(cb parser.NounVisitor, kids g.Value) (ret bool) {
	for _, k := range kids.Strings() {
		if ok := cb(MakeNoun(pt, k)); ok {
			ret = ok
			break
		}
	}
	return
}
