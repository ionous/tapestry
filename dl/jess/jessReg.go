package jess

import (
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// setup the default traits for the passed kind
func AddKindTraits(w weaver.Weaves, kind string, it *Traits) (err error) {
	for ; it != nil; it = it.Next() {
		str := it.Trait.String()
		if e := w.AddKindTrait(kind, inflect.Normalize(str)); e != nil {
			err = e
			break
		}
	}
	return
}
