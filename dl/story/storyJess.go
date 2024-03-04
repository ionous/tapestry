package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
)

// Execute - called by the macro runtime during weave.
func (op *DeclareStatement) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DeclareStatement) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequirePatterns, func(w *weave.Weaver) (err error) {
		if txt, e := safe.GetText(w, op.Text); e != nil {
			err = e
		} else {
			err = w.Generate(txt.String())
		}
		return
	})
}
