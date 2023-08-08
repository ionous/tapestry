package story

import (
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave"
)

// Execute - called by the macro runtime during weave.
func (op *StoryAlias) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *StoryAlias) Weave(cat *weave.Catalog) (err error) {
	return cat.Schedule(weave.RequireAll, func(w *weave.Weaver) (err error) {
		name := lang.Normalize(op.AsNoun)
		if n, e := w.GetClosestNoun(name); e != nil {
			err = e
		} else {
			pen := w.Pin()
			for _, a := range op.Names {
				if a := lang.Normalize(a); len(a) > 0 {
					if e := pen.AddName(n, a, -1); e != nil {
						err = e
						break
					}
				}
			}
		}
		return
	})
}

// Execute - called by the macro runtime during weave.
func (op *StoryDirective) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

// an ugly way to ensure that grammar ( and therefore the runtime )
// isnt dependent on story / weave
func (op *StoryDirective) Weave(cat *weave.Catalog) (err error) {
	// jump/skip/hop	{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}
	name := strings.Join(op.Lede, "/")
	return cat.Schedule(weave.RequireRules, func(w *weave.Weaver) error {
		return w.Pin().AddGrammar(name, &grammar.Directive{
			Lede:  op.Lede,
			Scans: op.Scans,
		})
	})
}
