package weave

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/assert"
)

type Weaver struct {
	d     *Domain
	at    string
	phase assert.Phase
	rt.Runtime
}

// fix: replace these with direct access to runtime.
func (w *Weaver) PluralOf(word string) (ret string) {
	if a, e := w.d.Pluralize(word); e == nil {
		ret = a
	}
	if len(ret) == 0 {
		ret = lang.Pluralize(word)
	}
	return
}

func (w *Weaver) SingularOf(word string) (ret string) {
	if a, e := w.d.Singularize(word); e == nil {
		ret = a
	}
	if len(ret) == 0 {
		ret = lang.Singularize(word)
	}
	return
}

func (w *Weaver) OppositeOf(word string) (ret string) {
	if a, e := w.d.FindOpposite(word); e == nil {
		ret = a
	}
	if len(ret) == 0 {
		ret = word
	}
	return
}
